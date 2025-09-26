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

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (s msgServer) SetEpochLength(ctx context.Context, msg *types.MsgSetEpochLength) (*types.MsgSetEpochLengthResponse, error) {
	if msg.Signer != s.authority {
		return nil, errors.Wrapf(types.ErrInvalidAuthority, "expected %s, got %s", s.authority, msg.Signer)
	}

	if msg.EpochLength <= 0 {
		return nil, errors.Wrap(types.ErrInvalidRequest, "invalid epoch length")
	}

	oldEpochLength, err := s.GetEpochLength(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get current epoch length from state")
	}

	err = s.setEpochLength(ctx, msg.EpochLength)
	if err != nil {
		return nil, errors.Wrap(err, "unable to set epoch length in state")
	}

	return &types.MsgSetEpochLengthResponse{}, s.eventService.EventManager(ctx).Emit(ctx, &types.EpochLengthSet{
		OldEpochLength: oldEpochLength,
		NewEpochLength: msg.EpochLength,
	})
}

func (s msgServer) SetHookAddress(ctx context.Context, msg *types.MsgSetHookAddress) (*types.MsgSetHookAddressResponse, error) {
	if msg.Signer != s.authority {
		return nil, errors.Wrapf(types.ErrInvalidAuthority, "expected %s, got %s", s.authority, msg.Signer)
	}

	if !common.IsHexAddress(msg.HookAddress) {
		return nil, errors.Wrap(types.ErrInvalidRequest, "invalid hook address")
	}
	// Because of the check above, we can safely decode the hex-encoded string.
	hookAddress := common.HexToAddress(msg.HookAddress)

	oldHookAddress, err := s.GetHookAddress(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get current hook address from state")
	}

	err = s.setHookAddress(ctx, hookAddress)
	if err != nil {
		return nil, errors.Wrap(err, "unable to set hook address in state")
	}

	return &types.MsgSetHookAddressResponse{}, s.eventService.EventManager(ctx).Emit(ctx, &types.HookAddressSet{
		OldHookAddress: oldHookAddress.String(),
		NewHookAddress: msg.HookAddress,
	})
}

func (s msgServer) SetEnrolledValidators(ctx context.Context, msg *types.MsgSetEnrolledValidators) (*types.MsgSetEnrolledValidatorsResponse, error) {
	if msg.Signer != s.authority {
		return nil, errors.Wrapf(types.ErrInvalidAuthority, "expected %s, got %s", s.authority, msg.Signer)
	}

	oldEnrolledValidators, err := s.GetEnrolledValidators(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get current enrolled validators from state")
	}

	err = s.enrolledValidators.Clear(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to clear old enrolled validators from state")
	}

	for _, address := range msg.EnrolledValidators {
		err := s.setEnrolledValidator(ctx, address)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgSetEnrolledValidatorsResponse{}, s.eventService.EventManager(ctx).Emit(ctx, &types.EnrolledValidatorsSet{
		OldEnrolledValidators: oldEnrolledValidators,
		NewEnrolledValidators: msg.EnrolledValidators,
	})
}
