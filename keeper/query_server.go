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

	"github.com/noble-assets/nova/types"
)

var _ types.QueryServer = &queryServer{}

type queryServer struct {
	*Keeper
}

func NewQueryServer(keeper *Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

func (s queryServer) Config(ctx context.Context, _ *types.QueryConfig) (*types.QueryConfigResponse, error) {
	epochLength, err := s.Keeper.GetEpochLength(ctx)
	if err != nil {
		return nil, err
	}
	hookAddress, err := s.Keeper.GetHookAddress(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryConfigResponse{
		EpochLength: epochLength,
		HookAddress: hookAddress.String(),
	}, nil
}

func (s queryServer) PendingEpoch(ctx context.Context, _ *types.QueryPendingEpoch) (*types.QueryEpochResponse, error) {
	epoch, err := s.GetPendingEpoch(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryEpochResponse{Epoch: epoch}, nil
}

func (s queryServer) FinalizedEpochs(ctx context.Context, _ *types.QueryFinalizedEpochs) (*types.QueryFinalizedEpochsResponse, error) {
	finalizedEpochs, err := s.GetFinalizedEpochs(ctx)

	return &types.QueryFinalizedEpochsResponse{FinalizedEpochs: finalizedEpochs}, err
}

func (s queryServer) LatestFinalizedEpoch(ctx context.Context, _ *types.QueryLatestFinalizedEpoch) (*types.QueryEpochResponse, error) {
	finalizedEpoch, err := s.GetLatestFinalizedEpoch(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryEpochResponse{Epoch: finalizedEpoch}, nil
}

func (s queryServer) FinalizedEpoch(ctx context.Context, req *types.QueryFinalizedEpoch) (*types.QueryEpochResponse, error) {
	epoch, err := s.GetFinalizedEpoch(ctx, req.EpochNumber)
	if err != nil {
		return nil, err
	}

	return &types.QueryEpochResponse{Epoch: epoch}, nil
}

func (s queryServer) StateRoots(ctx context.Context, _ *types.QueryStateRoots) (*types.QueryStateRootsResponse, error) {
	stateRoots, err := s.GetStateRoots(ctx)

	return &types.QueryStateRootsResponse{StateRoots: stateRoots}, err
}

func (s queryServer) LatestStateRoot(ctx context.Context, _ *types.QueryLatestStateRoot) (*types.QueryStateRootResponse, error) {
	stateRoot, err := s.GetLatestStateRoot(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryStateRootResponse{StateRoot: stateRoot.String()}, nil
}

func (s queryServer) StateRoot(ctx context.Context, req *types.QueryStateRoot) (*types.QueryStateRootResponse, error) {
	stateRoot, err := s.GetStateRoot(ctx, req.EpochNumber)
	if err != nil {
		return nil, err
	}

	return &types.QueryStateRootResponse{StateRoot: stateRoot.String()}, nil
}

func (s queryServer) MailboxRoots(ctx context.Context, _ *types.QueryMailboxRoots) (*types.QueryMailboxRootsResponse, error) {
	mailboxRoots, err := s.GetMailboxRoots(ctx)

	return &types.QueryMailboxRootsResponse{MailboxRoots: mailboxRoots}, err
}

func (s queryServer) LatestMailboxRoot(ctx context.Context, _ *types.QueryLatestMailboxRoot) (*types.QueryMailboxRootResponse, error) {
	mailboxRoot, err := s.GetLatestMailboxRoot(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryMailboxRootResponse{MailboxRoot: mailboxRoot.String()}, nil
}

func (s queryServer) MailboxRoot(ctx context.Context, req *types.QueryMailboxRoot) (*types.QueryMailboxRootResponse, error) {
	mailboxRoot, err := s.GetMailboxRoot(ctx, req.EpochNumber)
	if err != nil {
		return nil, err
	}

	return &types.QueryMailboxRootResponse{MailboxRoot: mailboxRoot.String()}, nil
}
