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

	"github.com/ethereum/go-ethereum/common"

	"github.com/noble-assets/nova/types"
)

// GetHookAddress returns the hook address from state.
func (k *Keeper) GetHookAddress(ctx context.Context) (common.Address, error) {
	hookAddress, err := k.hookAddress.Get(ctx)
	if err != nil {
		return common.Address{}, err
	}

	return common.Address(hookAddress), nil
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

// GetCurrentEpoch returns the current epoch from state.
func (k *Keeper) GetCurrentEpoch(ctx context.Context) (types.Epoch, error) {
	return k.currentEpoch.Get(ctx)
}

// setCurrentEpoch saves the current epoch to state.
func (k *Keeper) setCurrentEpoch(ctx context.Context, epoch types.Epoch) error {
	return k.currentEpoch.Set(ctx, epoch)
}

// GetEpoch returns a finalized epoch from state.
func (k *Keeper) GetEpoch(ctx context.Context, number uint64) (types.Epoch, error) {
	return k.epochs.Get(ctx, number)
}

// GetEpochs returns all finalized epochs from state.
func (k *Keeper) GetEpochs(ctx context.Context) (map[uint64]types.Epoch, error) {
	epochs := make(map[uint64]types.Epoch)

	err := k.epochs.Walk(
		ctx, nil,
		func(_ uint64, epoch types.Epoch) (stop bool, err error) {
			epochs[epoch.Number] = epoch
			return false, nil
		},
	)

	return epochs, err
}

// setEpoch saves a finalized epoch to state.
func (k *Keeper) setEpoch(ctx context.Context, epoch types.Epoch) error {
	return k.epochs.Set(ctx, epoch.Number, epoch)
}

// startNewEpoch ... TODO
func (k *Keeper) startNewEpoch(ctx context.Context, stateRoot common.Hash, mailboxRoot common.Hash) error {
	currentEpoch, err := k.GetCurrentEpoch(ctx)
	if err != nil {
		return err
	}
	epochLength, err := k.GetEpochLength(ctx)
	if err != nil {
		return err
	}

	err = k.setEpoch(ctx, currentEpoch)
	if err != nil {
		return err
	}
	err = k.setCurrentEpoch(ctx, types.Epoch{
		Number:    currentEpoch.Number + 1,
		EndHeight: currentEpoch.EndHeight + epochLength,
	})
	if err != nil {
		return err
	}

	err = k.setStateRoot(ctx, currentEpoch.Number, stateRoot)
	if err != nil {
		return err
	}
	err = k.setMailboxRoot(ctx, mailboxRoot)
	if err != nil {
		return err
	}

	return nil
}

// GetStateRoot returns a state root for an epoch from state.
func (k *Keeper) GetStateRoot(ctx context.Context, epoch uint64) (common.Hash, error) {
	stateRoot, err := k.stateRoots.Get(ctx, epoch)
	if err != nil {
		return common.Hash{}, err
	}

	return common.Hash(stateRoot), nil
}

// GetStateRoots returns all state roots from state.
func (k *Keeper) GetStateRoots(ctx context.Context) (map[uint64]string, error) {
	stateRoots := make(map[uint64]string)

	err := k.stateRoots.Walk(
		ctx, nil,
		func(epoch uint64, stateRoot []byte) (stop bool, err error) {
			stateRoots[epoch] = common.Hash(stateRoot).String()
			return false, nil
		},
	)

	return stateRoots, err
}

// setStateRoot saves a state root for an epoch to state.
func (k *Keeper) setStateRoot(ctx context.Context, epoch uint64, stateRoot common.Hash) error {
	return k.stateRoots.Set(ctx, epoch, stateRoot.Bytes())
}

// GetMailboxRoot returns the mailbox root from state.
func (k *Keeper) GetMailboxRoot(ctx context.Context) (common.Hash, error) {
	mailboxRoot, err := k.mailboxRoot.Get(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	return common.Hash(mailboxRoot), nil
}

// setMailboxRoot saves the mailbox root to state.
func (k *Keeper) setMailboxRoot(ctx context.Context, mailboxRoot common.Hash) error {
	return k.mailboxRoot.Set(ctx, mailboxRoot.Bytes())
}
