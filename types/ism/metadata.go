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
	"encoding/binary"

	"cosmossdk.io/errors"
)

// MetadataSize defines the byte encoded size of Metadata.
// index size + (# of proof leaves * leaf size)
const MetadataSize = 4 + 32*32

type Metadata struct {
	Index uint32
	Proof [32][32]byte
}

func ParseMetadata(bz []byte) (Metadata, error) {
	if len(bz) != MetadataSize {
		return Metadata{}, errors.Wrapf(ErrInvalidMetadata, "length %d != %d", len(bz), MetadataSize)
	}

	offset := 0

	index := binary.BigEndian.Uint32(bz[offset : offset+4])
	offset += 4

	var proof [32][32]byte
	for i := 0; i < 32; i++ {
		copy(proof[i][:], bz[offset:offset+32])
		offset += 32
	}

	return Metadata{
		Index: index,
		Proof: proof,
	}, nil
}

func (m Metadata) Bytes() []byte {
	bz := make([]byte, MetadataSize)

	offset := 0

	binary.BigEndian.PutUint32(bz[offset:], m.Index)
	offset += 4

	for i := 0; i < 32; i++ {
		copy(bz[offset:offset+32], m.Proof[i][:])
		offset += 32
	}

	return bz
}
