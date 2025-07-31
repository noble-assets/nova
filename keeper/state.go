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
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ethereum/go-ethereum/common"

	"github.com/noble-assets/nova/types"
)

// GetHookAddress returns the hook address from state.
func (k *Keeper) GetHookAddress(ctx context.Context) (common.Address, error) {
	hookAddress, err := k.hookAddress.Get(ctx)
	if err != nil {
		return common.Address{}, err
	}

	return common.BytesToAddress(hookAddress), nil
}

// setHookAddress saves the hook address to state.
func (k *Keeper) setHookAddress(ctx context.Context, hookAddress common.Address) error {
	return k.hookAddress.Set(ctx, hookAddress.Bytes())
}

// GetEpochLength returns the epoch length from state.
func (k *Keeper) GetEpochLength(ctx context.Context) (uint64, error) {
	return k.epochLength.Get(ctx)
}

// setEpochLength saves the epoch length to state.
func (k *Keeper) setEpochLength(ctx context.Context, length uint64) error {
	return k.epochLength.Set(ctx, length)
}

// GetPendingEpoch returns the currently pending epoch from state.
func (k *Keeper) GetPendingEpoch(ctx context.Context) (types.Epoch, error) {
	return k.pendingEpoch.Get(ctx)
}

// setPendingEpoch saves the currently pending epoch to state.
func (k *Keeper) setPendingEpoch(ctx context.Context, epoch types.Epoch) error {
	return k.pendingEpoch.Set(ctx, epoch)
}

// GetLatestFinalizedEpoch returns the latest finalized epoch from state.
func (k *Keeper) GetLatestFinalizedEpoch(ctx context.Context) (types.Epoch, error) {
	pendingEpoch, err := k.GetPendingEpoch(ctx)
	if err != nil {
		return types.Epoch{}, errors.New("no pending epoch")
	}

	if pendingEpoch.Number == 0 {
		return types.Epoch{}, errors.New("no finalized epoch")
	}

	latestEpochNumber := pendingEpoch.Number - 1
	return k.GetFinalizedEpoch(ctx, latestEpochNumber)
}

// GetFinalizedEpoch returns a finalized epoch from state.
func (k *Keeper) GetFinalizedEpoch(ctx context.Context, epochNumber uint64) (types.Epoch, error) {
	if has, _ := k.finalizedEpochs.Has(ctx, epochNumber); !has {
		return types.Epoch{}, fmt.Errorf("finalized epoch %d not found", epochNumber)
	}

	finalizedEpoch, err := k.finalizedEpochs.Get(ctx, epochNumber)
	if err != nil {
		return types.Epoch{}, fmt.Errorf("unable to get finalized epoch %d from state", epochNumber)
	}

	return finalizedEpoch, nil
}

// getFinalizedEpochs returns all finalized epochs from state.
func (k *Keeper) getFinalizedEpochs(ctx context.Context) (map[uint64]types.Epoch, error) {
	finalizedEpochs := make(map[uint64]types.Epoch)

	err := k.finalizedEpochs.Walk(
		ctx, nil,
		func(epochNumber uint64, finalizedEpoch types.Epoch) (stop bool, err error) {
			finalizedEpochs[epochNumber] = finalizedEpoch
			return false, nil
		},
	)

	return finalizedEpochs, err
}

// GetFinalizedEpochsPaginated returns all finalized epochs from state, paginated.
func (k *Keeper) GetFinalizedEpochsPaginated(ctx context.Context, req *query.PageRequest) ([]types.Epoch, *query.PageResponse, error) {
	return query.CollectionPaginate(
		ctx, k.finalizedEpochs, req,
		func(_ uint64, finalizedEpoch types.Epoch) (types.Epoch, error) {
			return finalizedEpoch, nil
		},
	)
}

// setFinalizedEpoch saves a finalized epoch to state.
func (k *Keeper) setFinalizedEpoch(ctx context.Context, epoch types.Epoch) error {
	return k.finalizedEpochs.Set(ctx, epoch.Number, epoch)
}

// startNewEpoch is a utility that starts a new epoch, marking the currently
// pending epoch as finalized given a state root and mailbox root.
func (k *Keeper) startNewEpoch(ctx context.Context, stateRoot common.Hash, mailboxRoot common.Hash) error {
	pendingEpoch, err := k.GetPendingEpoch(ctx)
	if err != nil {
		return err
	}
	epochLength, err := k.GetEpochLength(ctx)
	if err != nil {
		return err
	}

	err = k.setFinalizedEpoch(ctx, pendingEpoch)
	if err != nil {
		return err
	}
	err = k.setPendingEpoch(ctx, types.Epoch{
		Number:      pendingEpoch.Number + 1,
		StartHeight: pendingEpoch.EndHeight,
		EndHeight:   pendingEpoch.EndHeight + epochLength,
	})
	if err != nil {
		return err
	}

	err = k.setStateRoot(ctx, pendingEpoch.Number, stateRoot)
	if err != nil {
		return err
	}
	err = k.setMailboxRoot(ctx, pendingEpoch.Number, mailboxRoot)
	if err != nil {
		return err
	}

	return nil
}

// GetLatestStateRoot returns the latest finalized state root from state.
func (k *Keeper) GetLatestStateRoot(ctx context.Context) (common.Hash, error) {
	pendingEpoch, err := k.GetPendingEpoch(ctx)
	if err != nil {
		return common.Hash{}, errors.New("no pending epoch")
	}

	if pendingEpoch.Number == 0 {
		return common.Hash{}, errors.New("no finalized epoch")
	}

	latestEpochNumber := pendingEpoch.Number - 1
	return k.GetStateRoot(ctx, latestEpochNumber)
}

// GetStateRoot returns a state root for an epoch from state.
func (k *Keeper) GetStateRoot(ctx context.Context, epochNumber uint64) (common.Hash, error) {
	if has, _ := k.stateRoots.Has(ctx, epochNumber); !has {
		return common.Hash{}, fmt.Errorf("state root for epoch %d not found", epochNumber)
	}

	stateRoot, err := k.stateRoots.Get(ctx, epochNumber)
	if err != nil {
		return common.Hash{}, fmt.Errorf("unable to get state root for epoch %d from state", epochNumber)
	}

	return common.BytesToHash(stateRoot), nil
}

// getStateRoots returns all state roots from state.
func (k *Keeper) getStateRoots(ctx context.Context) (map[uint64]string, error) {
	stateRoots := make(map[uint64]string)

	err := k.stateRoots.Walk(
		ctx, nil,
		func(epochNumber uint64, stateRoot []byte) (stop bool, err error) {
			stateRoots[epochNumber] = common.BytesToHash(stateRoot).String()
			return false, nil
		},
	)

	return stateRoots, err
}

// GetStateRootsPaginated paginated returns all state roots from state, paginated.
func (k *Keeper) GetStateRootsPaginated(ctx context.Context, req *query.PageRequest) ([]types.QueryStateRootsResponse_Value, *query.PageResponse, error) {
	return query.CollectionPaginate(
		ctx, k.stateRoots, req,
		func(epochNumber uint64, stateRoot []byte) (types.QueryStateRootsResponse_Value, error) {
			return types.QueryStateRootsResponse_Value{
				EpochNumber: epochNumber,
				StateRoot:   common.BytesToHash(stateRoot).String(),
			}, nil
		},
	)
}

// setStateRoot saves a state root for an epoch to state.
func (k *Keeper) setStateRoot(ctx context.Context, epochNumber uint64, stateRoot common.Hash) error {
	return k.stateRoots.Set(ctx, epochNumber, stateRoot.Bytes())
}

// GetLatestMailboxRoot returns the latest finalized mailbox root from state.
func (k *Keeper) GetLatestMailboxRoot(ctx context.Context) (common.Hash, error) {
	pendingEpoch, err := k.GetPendingEpoch(ctx)
	if err != nil {
		return common.Hash{}, errors.New("no pending epoch")
	}

	if pendingEpoch.Number == 0 {
		return common.Hash{}, errors.New("no finalized epoch")
	}

	latestEpochNumber := pendingEpoch.Number - 1
	return k.GetMailboxRoot(ctx, latestEpochNumber)
}

// GetMailboxRoot returns a mailbox root for an epoch from state.
func (k *Keeper) GetMailboxRoot(ctx context.Context, epochNumber uint64) (common.Hash, error) {
	if has, _ := k.mailboxRoots.Has(ctx, epochNumber); !has {
		return common.Hash{}, fmt.Errorf("mailbox root for epoch %d not found", epochNumber)
	}

	mailboxRoot, err := k.mailboxRoots.Get(ctx, epochNumber)
	if err != nil {
		return common.Hash{}, fmt.Errorf("unable to get mailbox root for epoch %d from state", epochNumber)
	}

	return common.BytesToHash(mailboxRoot), nil
}

// getMailboxRoots returns all mailbox roots from state.
func (k *Keeper) getMailboxRoots(ctx context.Context) (map[uint64]string, error) {
	mailboxRoots := make(map[uint64]string)

	err := k.mailboxRoots.Walk(
		ctx, nil,
		func(epochNumber uint64, stateRoot []byte) (stop bool, err error) {
			mailboxRoots[epochNumber] = common.BytesToHash(stateRoot).String()
			return false, nil
		},
	)

	return mailboxRoots, err
}

// GetMailboxRootsPaginated returns all mailbox roots from state, paginated.
func (k *Keeper) GetMailboxRootsPaginated(ctx context.Context, req *query.PageRequest) ([]types.QueryMailboxRootsResponse_Value, *query.PageResponse, error) {
	return query.CollectionPaginate(
		ctx, k.mailboxRoots, req,
		func(epochNumber uint64, mailboxRoot []byte) (types.QueryMailboxRootsResponse_Value, error) {
			return types.QueryMailboxRootsResponse_Value{
				EpochNumber: epochNumber,
				MailboxRoot: common.BytesToHash(mailboxRoot).String(),
			}, nil
		},
	)
}

// setMailboxRoot saves a mailbox root for an epoch to state.
func (k *Keeper) setMailboxRoot(ctx context.Context, epochNumber uint64, mailboxRoot common.Hash) error {
	return k.mailboxRoots.Set(ctx, epochNumber, mailboxRoot.Bytes())
}
