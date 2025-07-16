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

	var currentEpoch types.Epoch
	if genesis.CurrentEpoch == nil {
		currentEpoch = types.Epoch{
			Number:    0,
			EndHeight: genesis.Config.EpochLength,
		}
	} else {
		currentEpoch = *genesis.CurrentEpoch
	}
	if err := k.setCurrentEpoch(ctx, currentEpoch); err != nil {
		panic(errors.Wrap(err, "failed to set genesis current epoch"))
	}

	for _, epoch := range genesis.Epochs {
		if err := k.setEpoch(ctx, epoch); err != nil {
			panic(errors.Wrapf(err, "failed to set genesis epoch %d", epoch.Number))
		}
	}

	for epoch, rawStateRoot := range genesis.StateRoots {
		stateRoot := common.HexToHash(rawStateRoot)

		if err := k.setStateRoot(ctx, epoch, stateRoot); err != nil {
			panic(errors.Wrapf(err, "failed to set genesis state root %d", epoch))
		}
	}

	mailboxRoot := common.HexToHash(genesis.MailboxRoot)
	if err := k.setMailboxRoot(ctx, mailboxRoot); err != nil {
		panic(errors.Wrap(err, "failed to set genesis mailbox root"))
	}
}

func (k *Keeper) ExportGenesis(ctx context.Context) *types.GenesisState {
	epochLength, _ := k.GetEpochLength(ctx)
	hookAddress, _ := k.GetHookAddress(ctx)
	config := types.Config{
		EpochLength: epochLength,
		HookAddress: hookAddress.String(),
	}

	currentEpoch, _ := k.GetCurrentEpoch(ctx)
	epochs, _ := k.GetEpochs(ctx)
	stateRoots, _ := k.GetStateRoots(ctx)
	mailboxRoot, _ := k.GetMailboxRoot(ctx)

	return &types.GenesisState{
		Config:       config,
		CurrentEpoch: &currentEpoch,
		Epochs:       epochs,
		StateRoots:   stateRoots,
		MailboxRoot:  mailboxRoot.String(),
	}
}
