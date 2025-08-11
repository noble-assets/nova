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

package cli

import (
	"fmt"
	"strconv"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/gogoproto/proto"
	"github.com/spf13/cobra"

	"github.com/noble-assets/nova/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(QueryFinalizedEpoch())
	cmd.AddCommand(QueryStateRoot())
	cmd.AddCommand(QueryMailboxRoot())

	return cmd
}

func QueryFinalizedEpoch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "finalized-epoch (epoch-number)",
		Short: "Query the latest or a specific finalized epoch",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			var res proto.Message
			var err error

			if len(args) == 1 {
				// This implies that we are querying for a specific finalized epoch.
				var epochNumber uint64
				epochNumber, err = strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					return errors.Wrap(err, "invalid epoch number")
				}

				res, err = queryClient.FinalizedEpoch(clientCtx.CmdContext, &types.QueryFinalizedEpoch{EpochNumber: epochNumber})
			} else {
				// This implies that we are querying for the latest finalized epoch.
				res, err = queryClient.LatestFinalizedEpoch(clientCtx.CmdContext, &types.QueryLatestFinalizedEpoch{})
			}

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryStateRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state-root (epoch-number)",
		Short: "Query the latest or a specific finalized state root",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			var res proto.Message
			var err error

			if len(args) == 1 {
				// This implies that we are querying for a specific state root.
				var epochNumber uint64
				epochNumber, err = strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					return errors.Wrap(err, "invalid epoch number")
				}

				res, err = queryClient.StateRoot(clientCtx.CmdContext, &types.QueryStateRoot{EpochNumber: epochNumber})
			} else {
				// This implies that we are querying for the latest state root.
				res, err = queryClient.LatestStateRoot(clientCtx.CmdContext, &types.QueryLatestStateRoot{})
			}

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryMailboxRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mailbox-root (epoch-number)",
		Short: "Query the latest or a specific finalized mailbox root",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			var res proto.Message
			var err error

			if len(args) == 1 {
				// This implies that we are querying for a specific mailbox root.
				var epochNumber uint64
				epochNumber, err = strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					return errors.Wrap(err, "invalid epoch number")
				}

				res, err = queryClient.MailboxRoot(clientCtx.CmdContext, &types.QueryMailboxRoot{EpochNumber: epochNumber})
			} else {
				// This implies that we are querying for the latest mailbox root.
				res, err = queryClient.LatestMailboxRoot(clientCtx.CmdContext, &types.QueryLatestMailboxRoot{})
			}

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
