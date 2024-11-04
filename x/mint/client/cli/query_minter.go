package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sge-network/sge/x/mint/types"
)

// GetCmdQueryInflation implements a command to return the current minting
// inflation value.
func GetCmdQueryInflation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inflation",
		Short: "Query the current minting inflation value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryInflationRequest{}
			res, err := queryClient.Inflation(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(fmt.Sprintf("%d\n", res.Inflation))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryPhaseStep implements a command to return the current minting
// inflation value.
func GetCmdQueryPhaseStep() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "phase-step",
		Short: "Query the current minting phase step",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryPhaseStepRequest{}
			res, err := queryClient.PhaseStep(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(fmt.Sprintf("%d\n", res.PhaseStep))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryPhaseProvisions implements a command to return the current minting
// phase provisions value.
func GetCmdQueryPhaseProvisions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "phase-provisions",
		Short: "Query the current minting phase provisions value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryPhaseProvisionsRequest{}
			res, err := queryClient.PhaseProvisions(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(fmt.Sprintf("%d\n", res.PhaseProvisions))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryEndPhaseStatus implements a command to return the current minting
// end phase status value.
func GetCmdQueryEndPhaseStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "endphase-status",
		Short: "Query the current status of end phase",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryEndPhaseStatusRequest{}
			res, err := queryClient.EndPhaseStatus(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(fmt.Sprintf("%t\n", res.IsInEndPhase))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
