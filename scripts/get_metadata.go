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

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	hyperlaneutil "github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/ethereum/go-ethereum/common"

	types "github.com/noble-assets/nova/types/ism"
)

type Response struct {
	Siblings  [][]string `json:"siblings"`
	LeafIndex int        `json:"leafIndex"`
}

func main() {
	raw, _ := http.Get("http://localhost:42069/prove/0xfc5d4c1f93dc68b02f4779ccf4988e1d1655a10b09dbe25470071ddbe2415a23")
	body, _ := io.ReadAll(raw.Body)

	var res Response
	_ = json.Unmarshal(body, &res)

	var proof [32][32]byte
	for i := 0; i < 32; i++ {
		proof[i], _ = hyperlaneutil.DecodeHexAddress(res.Siblings[i][0])
	}

	fmt.Println("0x" + common.Bytes2Hex(types.Metadata{
		Index: uint32(res.LeafIndex),
		Proof: proof,
	}.Bytes()))
}
