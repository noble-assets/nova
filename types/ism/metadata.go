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

const (
	IndexSize   = 4
	ProofLeaves = 32
	LeafSize    = 32

	// MetadataSize defines the byte encoded size of Metadata.
	MetadataSize = IndexSize + ProofLeaves*LeafSize
)

type Metadata struct {
	Index uint32
	Proof [ProofLeaves][LeafSize]byte
}

func ParseMetadata(bz []byte) (Metadata, error) {
	if len(bz) != MetadataSize {
		return Metadata{}, errors.Wrapf(ErrInvalidMetadata, "length %d != %d", len(bz), MetadataSize)
	}

	offset := 0

	index := binary.BigEndian.Uint32(bz[offset : offset+IndexSize])
	offset += IndexSize

	var proof [ProofLeaves][LeafSize]byte
	for i := range ProofLeaves {
		copy(proof[i][:], bz[offset:offset+LeafSize])
		offset += LeafSize
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
	offset += IndexSize

	for i := range ProofLeaves {
		copy(bz[offset:offset+LeafSize], m.Proof[i][:])
		offset += LeafSize
	}

	return bz
}
