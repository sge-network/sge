package mint

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	mintv1beta1 "github.com/sge-network/sge/api/sge/mint/v1beta"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: mintv1beta1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current minting parameters",
				},
				{
					RpcMethod: "Inflation",
					Use:       "inflation",
					Short:     "Query the current minting inflation value",
				},
				{
					RpcMethod: "PhaseProvisions",
					Use:       "phase-provisions",
					Short:     "Query the current minting phase provisions value",
				},
				{
					RpcMethod: "EndPhaseStatus",
					Use:       "endphase-status",
					Short:     "Query the current status of end phase",
				},
				{
					RpcMethod: "PhaseStep",
					Use:       "phase-step",
					Short:     "Query the current minting phase step",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: mintv1beta1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
			},
		},
	}
}
