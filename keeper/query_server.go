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

func (s queryServer) Config(ctx context.Context, req *types.QueryConfig) (*types.QueryConfigResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	epochLength, err := s.Keeper.GetEpochLength(ctx)
	if err != nil {
		return nil, err
	}
	hookAddress, err := s.Keeper.GetHookAddress(ctx)
	if err != nil {
		return nil, err
	}
	enrolledValidators, err := s.Keeper.GetEnrolledValidators(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryConfigResponse{
		EpochLength:        epochLength,
		HookAddress:        hookAddress.String(),
		EnrolledValidators: enrolledValidators,
	}, nil
}

func (s queryServer) PendingEpoch(ctx context.Context, req *types.QueryPendingEpoch) (*types.QueryEpochResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	epoch, err := s.GetPendingEpoch(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryEpochResponse{Epoch: epoch}, nil
}

func (s queryServer) FinalizedEpochs(ctx context.Context, req *types.QueryFinalizedEpochs) (*types.QueryFinalizedEpochsResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	finalizedEpochs, pagination, err := s.GetFinalizedEpochsPaginated(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.QueryFinalizedEpochsResponse{
		FinalizedEpochs: finalizedEpochs,
		Pagination:      pagination,
	}, nil
}

func (s queryServer) LatestFinalizedEpoch(ctx context.Context, req *types.QueryLatestFinalizedEpoch) (*types.QueryEpochResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	finalizedEpoch, err := s.GetLatestFinalizedEpoch(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryEpochResponse{Epoch: finalizedEpoch}, nil
}

func (s queryServer) FinalizedEpoch(ctx context.Context, req *types.QueryFinalizedEpoch) (*types.QueryEpochResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	epoch, err := s.GetFinalizedEpoch(ctx, req.EpochNumber)
	if err != nil {
		return nil, err
	}

	return &types.QueryEpochResponse{Epoch: epoch}, nil
}

func (s queryServer) StateRoots(ctx context.Context, req *types.QueryStateRoots) (*types.QueryStateRootsResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	stateRoots, pagination, err := s.GetStateRootsPaginated(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.QueryStateRootsResponse{
		StateRoots: stateRoots,
		Pagination: pagination,
	}, nil
}

func (s queryServer) LatestStateRoot(ctx context.Context, req *types.QueryLatestStateRoot) (*types.QueryStateRootResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	stateRoot, err := s.GetLatestStateRoot(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryStateRootResponse{StateRoot: stateRoot.String()}, nil
}

func (s queryServer) StateRoot(ctx context.Context, req *types.QueryStateRoot) (*types.QueryStateRootResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	stateRoot, err := s.GetStateRoot(ctx, req.EpochNumber)
	if err != nil {
		return nil, err
	}

	return &types.QueryStateRootResponse{StateRoot: stateRoot.String()}, nil
}

func (s queryServer) MailboxRoots(ctx context.Context, req *types.QueryMailboxRoots) (*types.QueryMailboxRootsResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	mailboxRoots, pagination, err := s.GetMailboxRootsPaginated(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.QueryMailboxRootsResponse{
		MailboxRoots: mailboxRoots,
		Pagination:   pagination,
	}, err
}

func (s queryServer) LatestMailboxRoot(ctx context.Context, req *types.QueryLatestMailboxRoot) (*types.QueryMailboxRootResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	mailboxRoot, err := s.GetLatestMailboxRoot(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryMailboxRootResponse{MailboxRoot: mailboxRoot.String()}, nil
}

func (s queryServer) MailboxRoot(ctx context.Context, req *types.QueryMailboxRoot) (*types.QueryMailboxRootResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	mailboxRoot, err := s.GetMailboxRoot(ctx, req.EpochNumber)
	if err != nil {
		return nil, err
	}

	return &types.QueryMailboxRootResponse{MailboxRoot: mailboxRoot.String()}, nil
}
