// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, NASD Inc. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/noble-assets/nova/types"
	"github.com/noble-assets/nova/types/abi"
)

type VoteExtensionNova struct {
	EpochNumber uint64      `json:"epoch_number"`
	StateRoot   common.Hash `json:"state_root"`
	MailboxRoot common.Hash `json:"mailbox_root"`
}
type VoteExtension struct {
	Nova VoteExtensionNova `json:"nova"`
}

// ExtendVoteHandler implements the Cosmos SDK interface for extending CometBFT
// votes. It extends votes with epoch finalization data including the Noble
// AppLayer state root and mailbox root from the current epoch's end height.
func (k *Keeper) ExtendVoteHandler(txConfig client.TxConfig) sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
		epoch, err := k.GetPendingEpoch(ctx)
		if err != nil {
			return nil, err
		}

		// Because vote extensions are made available in the next block via an
		// injected transaction, we have to check if we have an injection for
		// the current epoch. If we do, this implies the epoch is in the
		// process of being finalized, and we should wait before extending the
		// vote again.
		injection := parseInjection(req.Txs, txConfig.TxDecoder())
		if injection != nil && injection.EpochNumber == epoch.Number {
			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, nil
		}

		height := big.NewInt(int64(epoch.EndHeight))
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		block, err := k.client.BlockByNumber(ctxWithTimeout, height)
		if err != nil {
			if !errors.Is(err, ethereum.NotFound) {
				// An example of this case would be that the local AppLayer
				// node is inaccessible. An error returned during this step
				// doesn't hinder the validator, and it can continue producing
				// blocks.
				return nil, err
			}

			// If the block can't be found, this implies that the epoch isn't
			// ready to be finalized, and so we skip the vote extension process
			// to not pollute blocks with empty injections.
			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, nil
		}
		stateRoot := block.Root()

		var mailboxRoot [32]byte
		hookAddress, err := k.GetHookAddress(ctx)
		if err == nil {
			hook, err := abi.NewMerkleTreeHook(hookAddress, k.client)
			if err == nil {
				mailboxRoot, _ = hook.Root(&bind.CallOpts{
					BlockNumber: height,
					Context:     ctxWithTimeout,
				})
			}
		}

		bz, err := json.Marshal(VoteExtension{
			Nova: VoteExtensionNova{
				EpochNumber: epoch.Number,
				StateRoot:   stateRoot,
				MailboxRoot: mailboxRoot,
			},
		})
		if err != nil {
			return nil, err
		}

		k.logger.Info(fmt.Sprintf("extending vote for epoch %d", epoch.Number), "stateRoot", stateRoot, "mailboxRoot", common.Hash(mailboxRoot), "height", req.Height)

		return &abci.ResponseExtendVote{
			VoteExtension: bz,
		}, nil
	}
}

// PrepareProposalHandler implements the Cosmos SDK interface for modifying the
// default proposal preparation logic. It is called by the current block
// proposer, and injects the vote extensions as the first transaction in the
// block.
func (k *Keeper) PrepareProposalHandler(txConfig client.TxConfig) sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
		res := &abci.ResponsePrepareProposal{Txs: req.Txs}

		if voteExtensionsDisabled(ctx) {
			return res, nil
		}

		err := baseapp.ValidateVoteExtensions(ctx, k.valStore, ctx.BlockHeight(), ctx.ChainID(), req.LocalLastCommit)
		if err != nil {
			return nil, err
		}

		extension := k.computeVoteExtension(ctx, req.LocalLastCommit)
		if extension == nil {
			return res, nil
		}

		builder := txConfig.NewTxBuilder()
		err = builder.SetMsgs(&types.Injection{
			EpochNumber: extension.Nova.EpochNumber,
			StateRoot:   extension.Nova.StateRoot.String(),
			MailboxRoot: extension.Nova.MailboxRoot.String(),
			CommitInfo:  req.LocalLastCommit,
		})
		if err != nil {
			return nil, err
		}
		bz, err := txConfig.TxEncoder()(builder.GetTx())
		if err != nil {
			return nil, err
		}

		txs := slices.Insert(req.Txs, 0, bz)
		return &abci.ResponsePrepareProposal{Txs: txs}, nil
	}
}

// ProcessProposalHandler implements the Cosmos SDK interface for modifying the
// default proposal processing logic. It validates that injected epoch
// finalization data matches the computed vote extension consensus from the
// previous block's commit.
func (k *Keeper) ProcessProposalHandler(txConfig client.TxConfig) sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		accept := &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}
		reject := &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}

		if voteExtensionsDisabled(ctx) {
			return accept, nil
		}

		injection := parseInjection(req.Txs, txConfig.TxDecoder())
		if injection == nil {
			return accept, nil
		}

		err := baseapp.ValidateVoteExtensions(ctx, k.valStore, ctx.BlockHeight(), ctx.ChainID(), injection.CommitInfo)
		if err != nil {
			return nil, err
		}

		extension := k.computeVoteExtension(ctx, injection.CommitInfo)
		if extension == nil {
			return reject, nil
		}

		if injection.EpochNumber != extension.Nova.EpochNumber {
			return reject, nil
		}
		if !bytes.Equal(common.HexToHash(injection.StateRoot).Bytes(), extension.Nova.StateRoot.Bytes()) {
			return reject, nil
		}
		if !bytes.Equal(common.HexToHash(injection.MailboxRoot).Bytes(), extension.Nova.MailboxRoot.Bytes()) {
			return reject, nil
		}

		return accept, nil
	}
}

// PreBlockerHandler implements the Cosmos SDK interface for pre-blockers. It
// processes injected epoch finalization data to start a new epoch.
func (k *Keeper) PreBlockerHandler(txConfig client.TxConfig) sdk.PreBlocker {
	return func(ctx sdk.Context, req *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
		res := &sdk.ResponsePreBlock{ConsensusParamsChanged: false}

		if voteExtensionsDisabled(ctx) {
			return res, nil
		}

		injection := parseInjection(req.Txs, txConfig.TxDecoder())
		if injection != nil {
			stateRoot := common.HexToHash(injection.StateRoot)
			mailboxRoot := common.HexToHash(injection.MailboxRoot)

			err := k.startNewEpoch(ctx, stateRoot, mailboxRoot)
			if err != nil {
				// If we fail to start a new epoch, we simply log the error as we want block production to continue.
				k.logger.Error("failed to start new epoch", "err", err)
				return res, nil
			} else {
				k.logger.Info(fmt.Sprintf("finalized epoch %d", injection.EpochNumber), "height", req.Height)

				err = k.eventService.EventManager(ctx).Emit(ctx, &types.EpochFinalized{
					EpochNumber: injection.EpochNumber,
					StateRoot:   injection.StateRoot,
					MailboxRoot: injection.MailboxRoot,
				})
				if err != nil {
					// If we fail to emit the event, we simply log the error as we want block production to continue.
					k.logger.Error("failed to emit finalized epoch event", "err", err)
				}
			}
		}

		return res, nil
	}
}

// ----- Utilities -----

func (k *Keeper) computeVoteExtension(ctx context.Context, info abci.ExtendedCommitInfo) *VoteExtension {
	enrolledValidators, _ := k.GetEnrolledValidators(ctx)

	var totalPower int64
	tallies := make(map[string]int64)

	var winner string
	var winnerPower int64
	for _, vote := range info.Votes {
		if vote.BlockIdFlag != cmtproto.BlockIDFlagCommit {
			continue
		}
		if len(vote.VoteExtension) == 0 {
			// If there are enrolled validators, we check if this vote
			// extension belongs to an enrolled validator, otherwise we skip
			// them. If there are no enrolled validators, we default to all
			// validators being enrolled.
			if len(enrolledValidators) > 0 {
				if has, _ := k.enrolledValidators.Has(ctx, vote.Validator.Address); !has {
					continue
				}
			}
		}

		totalPower += vote.Validator.Power

		key := string(vote.VoteExtension)
		tallies[key] += vote.Validator.Power
		newPower := tallies[key]
		if newPower > winnerPower {
			winner = key
			winnerPower = newPower
		}
	}

	if len(tallies) == 0 {
		return nil
	}

	// NOTE: This is equivalent to doing winnerPower/totalPower > 2/3
	if winnerPower*3 > totalPower*2 {
		var extension VoteExtension
		if err := json.Unmarshal([]byte(winner), &extension); err != nil {
			return nil
		}

		return &extension
	} else {
		return nil
	}
}

func parseInjection(txs [][]byte, txDecoder sdk.TxDecoder) *types.Injection {
	// Because both Nova and Jester optionally inject transactions, we have to
	// handle all three different cases of injections.
	limit := len(txs)
	maxRange := 2
	if limit > maxRange {
		limit = maxRange
	}

	for _, tx := range txs[:limit] {
		if inj := parseInjectionFromTx(tx, txDecoder); inj != nil {
			return inj
		}
	}

	return nil
}

func parseInjectionFromTx(bz []byte, txDecoder sdk.TxDecoder) *types.Injection {
	tx, err := txDecoder(bz)
	if err != nil {
		return nil
	}

	msgs := tx.GetMsgs()
	if len(msgs) != 1 {
		return nil
	}

	injection, ok := msgs[0].(*types.Injection)
	if !ok {
		return nil
	}

	return injection
}

func voteExtensionsDisabled(ctx sdk.Context) bool {
	voteExtensionsEnableHeight := ctx.ConsensusParams().Abci.VoteExtensionsEnableHeight
	return voteExtensionsEnableHeight == 0 || ctx.BlockHeight() <= voteExtensionsEnableHeight
}
