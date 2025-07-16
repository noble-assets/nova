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
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/noble-assets/nova/types"
)

type Keeper struct {
	client   *ethclient.Client
	logger   log.Logger
	valStore baseapp.ValidatorStore

	hookAddress  collections.Item[[]byte]
	epochLength  collections.Item[uint64]
	currentEpoch collections.Item[types.Epoch]
	epochs       collections.Map[uint64, types.Epoch]
	stateRoots   collections.Map[uint64, []byte]
	mailboxRoot  collections.Item[[]byte]
}

func NewKeeper(cdc codec.BinaryCodec, store store.KVStoreService, logger log.Logger, rpcAddress string, valStore baseapp.ValidatorStore) *Keeper {
	var err error
	builder := collections.NewSchemaBuilder(store)

	var client *ethclient.Client
	if rpcAddress != "" {
		client, err = ethclient.Dial(rpcAddress)
		if err != nil {
			panic(err)
		}
	}

	keeper := &Keeper{
		client:   client,
		logger:   logger.With("module", types.ModuleName),
		valStore: valStore,

		hookAddress:  collections.NewItem(builder, types.HookAddressKey, "hook_address", collections.BytesValue),
		epochLength:  collections.NewItem(builder, types.EpochLengthKey, "epoch_length", collections.Uint64Value),
		currentEpoch: collections.NewItem(builder, types.CurrentEpochKey, "current_epoch", codec.CollValue[types.Epoch](cdc)),
		epochs:       collections.NewMap(builder, types.EpochPrefix, "epochs", collections.Uint64Key, codec.CollValue[types.Epoch](cdc)),
		stateRoots:   collections.NewMap(builder, types.StateRootPrefix, "state_roots", collections.Uint64Key, collections.BytesValue),
		mailboxRoot:  collections.NewItem(builder, types.MailboxRootKey, "mailbox_root", collections.BytesValue),
	}

	_, err = builder.Build()
	if err != nil {
		panic(err)
	}

	return keeper
}
