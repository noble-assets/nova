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

import hyperlaneutil "github.com/bcp-innovations/hyperlane-cosmos/util"

const SubmoduleName = "nova/ism"

// ExpectedId defines the expected ISM ID for this submodule. It was derived by
// concatenating a module specifier ("router_ism"), internal type (255), and
// internal id (0). This is aligned with how the Hyperlane x/core module
// derives their default ISM IDs.
var ExpectedId, _ = hyperlaneutil.DecodeHexAddress("0x726f757465725f69736d00000000000000000000000000ff0000000000000000")

// ModuleId defines the expected Module ID for this ISM. We have chosen 255, as
// it's the largest uint8, to not run into duplicate registrations with the
// default Hyperlane ISMs.
const ModuleId = uint8(255)

var PausedKey = []byte("ism/paused")
