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
	"cosmossdk.io/collections"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/noble-assets/nova/types"
)

type Keeper struct {
	authority string

	client        *ethclient.Client
	codec         codec.BinaryCodec
	eventService  event.Service
	logger        log.Logger
	stakingKeeper types.StakingKeeper

	epochLength        collections.Item[uint64]
	hookAddress        collections.Item[[]byte]
	enrolledValidators collections.Map[[]byte, string]
	pendingEpoch       collections.Item[types.Epoch]
	finalizedEpochs    collections.Map[uint64, types.Epoch]
	stateRoots         collections.Map[uint64, []byte]
	mailboxRoots       collections.Map[uint64, []byte]
}

func NewKeeper(authority string, cdc codec.BinaryCodec, storeService store.KVStoreService, eventService event.Service, logger log.Logger, rpcAddress string, stakingKeeper types.StakingKeeper) *Keeper {
	builder := collections.NewSchemaBuilder(storeService)

	client, err := ethclient.Dial(rpcAddress)
	if err != nil {
		panic(err)
	}

	keeper := &Keeper{
		authority: authority,

		client:        client,
		codec:         cdc,
		eventService:  eventService,
		logger:        logger.With("module", types.ModuleName),
		stakingKeeper: stakingKeeper,

		epochLength:        collections.NewItem(builder, types.EpochLengthKey, "epoch_length", collections.Uint64Value),
		hookAddress:        collections.NewItem(builder, types.HookAddressKey, "hook_address", collections.BytesValue),
		enrolledValidators: collections.NewMap(builder, types.EnrolledValidatorPrefix, "enrolled_validators", collections.BytesKey, collections.StringValue),
		pendingEpoch:       collections.NewItem(builder, types.PendingEpochKey, "pending_epoch", codec.CollValue[types.Epoch](cdc)),
		finalizedEpochs:    collections.NewMap(builder, types.FinalizedEpochPrefix, "finalized_epochs", collections.Uint64Key, codec.CollValue[types.Epoch](cdc)),
		stateRoots:         collections.NewMap(builder, types.StateRootPrefix, "state_roots", collections.Uint64Key, collections.BytesValue),
		mailboxRoots:       collections.NewMap(builder, types.MailboxRootPrefix, "mailbox_roots", collections.Uint64Key, collections.BytesValue),
	}

	_, err = builder.Build()
	if err != nil {
		panic(err)
	}

	return keeper
}
