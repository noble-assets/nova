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
	"context"

	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/noble-assets/nova/types"
)

func (k *Keeper) InitGenesis(ctx context.Context, genesis types.GenesisState) {
	if err := k.setEpochLength(ctx, genesis.Config.EpochLength); err != nil {
		panic(errors.Wrap(err, "failed to set genesis epoch length"))
	}

	hookAddress := common.HexToAddress(genesis.Config.HookAddress)
	if err := k.setHookAddress(ctx, hookAddress); err != nil {
		panic(errors.Wrap(err, "failed to set genesis hook address"))
	}

	for _, address := range genesis.Config.EnrolledValidators {
		err := k.setEnrolledValidator(ctx, address)
		if err != nil {
			panic(err)
		}
	}

	var pendingEpoch types.Epoch
	if genesis.PendingEpoch == nil {
		pendingEpoch = types.Epoch{
			Number:      0,
			StartHeight: 0,
			EndHeight:   genesis.Config.EpochLength,
		}
	} else {
		pendingEpoch = *genesis.PendingEpoch
	}
	if err := k.setPendingEpoch(ctx, pendingEpoch); err != nil {
		panic(errors.Wrap(err, "failed to set genesis pending epoch"))
	}

	for _, finalizedEpoch := range genesis.FinalizedEpochs {
		if err := k.setFinalizedEpoch(ctx, finalizedEpoch); err != nil {
			panic(errors.Wrapf(err, "failed to set genesis finalized epoch %d", finalizedEpoch.Number))
		}
	}

	for epochNumber, rawStateRoot := range genesis.StateRoots {
		stateRoot := common.HexToHash(rawStateRoot)

		if err := k.setStateRoot(ctx, epochNumber, stateRoot); err != nil {
			panic(errors.Wrapf(err, "failed to set genesis state root %d", epochNumber))
		}
	}

	for epochNumber, rawMailboxRoot := range genesis.MailboxRoots {
		mailboxRoot := common.HexToHash(rawMailboxRoot)

		if err := k.setMailboxRoot(ctx, epochNumber, mailboxRoot); err != nil {
			panic(errors.Wrapf(err, "failed to set genesis mailbox root %d", epochNumber))
		}
	}
}

func (k *Keeper) ExportGenesis(ctx context.Context) *types.GenesisState {
	epochLength, err := k.GetEpochLength(ctx)
	if err != nil {
		k.logger.Warn("unable to get epoch length", "err", err)
	}
	hookAddress, err := k.GetHookAddress(ctx)
	if err != nil {
		k.logger.Warn("unable to get hook address", "err", err)
	}
	enrolledValidators, err := k.GetEnrolledValidators(ctx)
	if err != nil {
		k.logger.Warn("unable to get enrolled validators", "err", err)
	}
	config := types.Config{
		EpochLength:        epochLength,
		HookAddress:        hookAddress.String(),
		EnrolledValidators: enrolledValidators,
	}

	pendingEpoch, err := k.GetPendingEpoch(ctx)
	if err != nil {
		k.logger.Warn("unable to get pending epoch", "err", err)
	}
	finalizedEpochs, err := k.getFinalizedEpochs(ctx)
	if err != nil {
		k.logger.Warn("unable to get finalized epochs", "err", err)
	}
	stateRoots, err := k.getStateRoots(ctx)
	if err != nil {
		k.logger.Warn("unable to get state roots", "err", err)
	}
	mailboxRoots, err := k.getMailboxRoots(ctx)
	if err != nil {
		k.logger.Warn("unable to get mailbox roots", "err", err)
	}

	return &types.GenesisState{
		Config:          config,
		PendingEpoch:    &pendingEpoch,
		FinalizedEpochs: finalizedEpochs,
		StateRoots:      stateRoots,
		MailboxRoots:    mailboxRoots,
	}
}
