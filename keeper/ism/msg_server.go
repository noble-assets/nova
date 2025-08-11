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

	"cosmossdk.io/errors"

	types "github.com/noble-assets/nova/types/ism"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (s msgServer) Pause(ctx context.Context, msg *types.MsgPause) (*types.MsgPauseResponse, error) {
	if msg.Signer != s.authority {
		return nil, errors.Wrapf(types.ErrInvalidAuthority, "expected %s, got %s", s.authority, msg.Signer)
	}

	if s.GetPaused(ctx) {
		return nil, errors.Wrap(types.ErrInvalidRequest, "already paused")
	}

	err := s.setPaused(ctx, true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to pause")
	}

	return &types.MsgPauseResponse{}, s.eventService.EventManager(ctx).Emit(ctx, &types.Paused{})
}

func (s msgServer) Unpause(ctx context.Context, msg *types.MsgUnpause) (*types.MsgUnpauseResponse, error) {
	if msg.Signer != s.authority {
		return nil, errors.Wrapf(types.ErrInvalidAuthority, "expected %s, got %s", s.authority, msg.Signer)
	}

	if !s.GetPaused(ctx) {
		return nil, errors.Wrap(types.ErrInvalidRequest, "already unpaused")
	}

	err := s.setPaused(ctx, false)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unpause")
	}

	return &types.MsgUnpauseResponse{}, s.eventService.EventManager(ctx).Emit(ctx, &types.Unpaused{})
}
