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

package ism

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	hyperlaneutil "github.com/bcp-innovations/hyperlane-cosmos/util"

	types "github.com/noble-assets/nova/types/ism"
)

var _ hyperlaneutil.InterchainSecurityModule = &Keeper{}

type Keeper struct {
	authority string

	eventService event.Service
	logger       log.Logger

	paused collections.Item[bool]
}

func NewKeeper(authority string, storeService store.KVStoreService, eventService event.Service, logger log.Logger, hyperlaneKeeper types.HyperlaneKeeper) *Keeper {
	builder := collections.NewSchemaBuilder(storeService)

	keeper := &Keeper{
		authority: authority,

		eventService: eventService,
		logger:       logger.With("module", types.SubmoduleName),

		paused: collections.NewItem(builder, types.PausedKey, "ism_paused", collections.BoolValue),
	}

	_, err := builder.Build()
	if err != nil {
		panic(err)
	}

	// We must register our ISM with the Hyperlane module so that messages from
	// the Noble AppLayer are correctly routed. We register our ISM as 255, as
	// it's the largest uint8, to not run into duplicate registrations with the
	// default Hyperlane ISMs.
	hyperlaneKeeper.IsmRouter().RegisterModule(255, keeper)

	return keeper
}

// Exists implements the expected Hyperlane InterchainSecurityModule interface.
func (k *Keeper) Exists(ctx context.Context, ismId hyperlaneutil.HexAddress) (bool, error) {
	// TODO implement me
	panic("implement me")
}

// Verify implements the expected Hyperlane InterchainSecurityModule interface.
func (k *Keeper) Verify(ctx context.Context, ismId hyperlaneutil.HexAddress, metadata []byte, message hyperlaneutil.HyperlaneMessage) (bool, error) {
	// TODO implement me
	panic("implement me")
}
