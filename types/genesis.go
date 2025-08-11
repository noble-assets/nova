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

package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Config: Config{
			EpochLength: 100, // 5 secs @ 50 ms AppLayer block time.
			HookAddress: common.Address{}.String(),
		},
	}
}

func (genesis *GenesisState) Validate() error {
	if genesis.Config.EpochLength <= 0 {
		return fmt.Errorf("invalid nova epoch length: %d", genesis.Config.EpochLength)
	}

	if valid := common.IsHexAddress(genesis.Config.HookAddress); !valid {
		return fmt.Errorf("invalid nova hook address: %s", genesis.Config.HookAddress)
	}

	// TODO: Should we validate finalizedEpochs?

	// TODO(stateRoots, mailboxRoots): go-ethereum doesn't provide a way of validating a hash

	return nil
}
