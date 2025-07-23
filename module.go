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

package nova

import (
	"context"
	"encoding/json"
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	modulev1 "github.com/noble-assets/nova/api/module/v1"
	novav1 "github.com/noble-assets/nova/api/v1"
	"github.com/noble-assets/nova/client/cli"
	"github.com/noble-assets/nova/keeper"
	"github.com/noble-assets/nova/types"
)

// ConsensusVersion defines the current Nova module consensus version.
const ConsensusVersion = 1

var (
	_ module.AppModuleBasic      = AppModule{}
	_ appmodule.AppModule        = AppModule{}
	_ module.HasConsensusVersion = AppModule{}
	_ module.HasGenesis          = AppModule{}
	_ module.HasGenesisBasics    = AppModuleBasic{}
	_ module.HasServices         = AppModule{}
)

//

type AppModuleBasic struct{}

func NewAppModuleBasic() AppModuleBasic {
	return AppModuleBasic{}
}

func (AppModuleBasic) Name() string { return types.ModuleName }

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (AppModuleBasic) RegisterInterfaces(reg codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (b AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var genesis types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genesis); err != nil {
		return fmt.Errorf("failed to unmarshal Nova genesis state: %w", err)
	}

	return genesis.Validate()
}

//

type AppModule struct {
	AppModuleBasic

	keeper *keeper.Keeper
}

func NewAppModule(keeper *keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(),
		keeper:         keeper,
	}
}

func (AppModule) IsOnePerModuleType() {}

func (AppModule) IsAppModule() {}

func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

func (m AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) {
	var genesis types.GenesisState
	cdc.MustUnmarshalJSON(bz, &genesis)

	m.keeper.InitGenesis(ctx, genesis)
}

func (m AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genesis := m.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(genesis)
}

func (m AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(m.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(m.keeper))
}

//

func (AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: novav1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod:      "SetEpochLength",
					Use:            "set-epoch-length [epoch-length]",
					Short:          "Set a new epoch length (authority gated)",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "epoch_length"}},
				},
				{
					RpcMethod:      "SetHookAddress",
					Use:            "set-hook-address [hook-address]",
					Short:          "Set a new hook address (authority gated)",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "hook_address"}},
				},
			},
		},
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: novav1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Config",
					Use:       "config",
					Short:     "Query the module configuration",
				},
				{
					RpcMethod: "PendingEpoch",
					Use:       "pending-epoch",
					Short:     "Query the currently pending epoch",
				},
				{
					RpcMethod: "FinalizedEpochs",
					Use:       "finalized-epochs",
					Short:     "Query all finalized epochs",
				},
				// NOTE: LatestFinalizedEpoch and FinalizedEpoch are combined into a single custom command.
				{
					RpcMethod: "LatestFinalizedEpoch",
					Skip:      true,
				},
				{
					RpcMethod: "FinalizedEpoch",
					Skip:      true,
				},
				{
					RpcMethod: "StateRoots",
					Use:       "state-roots",
					Short:     "Query all finalized state roots",
				},
				// NOTE: LastestStateRoot and StateRoot are combined into a single custom command.
				{
					RpcMethod: "LatestStateRoot",
					Skip:      true,
				},
				{
					RpcMethod: "StateRoot",
					Skip:      true,
				},
				{
					RpcMethod: "MailboxRoots",
					Use:       "mailbox-roots",
					Short:     "Query all finalized mailbox roots",
				},
				// NOTE: LatestMailboxRoot and MailboxRoot are combines into a single custom command.
				{
					RpcMethod: "LatestMailboxRoot",
					Skip:      true,
				},
				{
					RpcMethod: "MailboxRoot",
					Skip:      true,
				},
			},
			EnhanceCustomCommand: true,
		},
	}
}

func (AppModule) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

//

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config *modulev1.Module

	Codec          codec.Codec
	StoreService   store.KVStoreService
	EventService   event.Service
	Logger         log.Logger
	ValidatorStore baseapp.ValidatorStore

	AppOpts servertypes.AppOptions `optional:"true"`
	Viper   *viper.Viper           `optional:"true"`
}

type ModuleOutputs struct {
	depinject.Out

	Keeper *keeper.Keeper
	Module appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	if in.Config.Authority == "" {
		panic("authority for nova module must be set")
	}

	var rpcAddress string
	if in.Viper != nil { // viper takes precedence over app options
		rpcAddress = in.Viper.GetString(FlagRPCAddress)
	} else if in.AppOpts != nil {
		rpcAddress = cast.ToString(in.AppOpts.Get(FlagRPCAddress))
	}

	authority := authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	k := keeper.NewKeeper(authority.String(), in.Codec, in.StoreService, in.EventService, in.Logger, rpcAddress, in.ValidatorStore)
	m := NewAppModule(k)

	return ModuleOutputs{Keeper: k, Module: m}
}
