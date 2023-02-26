package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/sge-network/sge/x/orderbook/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// GetCmdQueryParticipationExposures implements the command to query all the participation exposures to a specific orderbook.
func GetCmdQueryParticipationExposures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "participation-exposures [order-book-id]",
		Short: "Query all participation exposures for a specific order book",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query participation exposures on an individual order book.

Example:
$ %s query orderbook participation-exposures %s
`,
				version.AppName, "5531c60f-2025-48ce-ae79-1dc110f16000",
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			orderBookUID := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryParticipationExposuresRequest{
				BookUid:    orderBookUID,
				Pagination: pageReq,
			}

			res, err := queryClient.ParticipationExposures(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "participation exposures")

	return cmd
}

// GetCmdQueryParticipationExposure implements the participationexposure query command.
func GetCmdQueryParticipationExposure() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "participation-exposure [order-book-id] [participation-number]",
		Short: "Query a participation exposure",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about a participation exposure.

Example:
$ %s query orderbook participation-exposure %s %d
`,
				version.AppName, "5531c60f-2025-48ce-ae79-1dc110f16000", 1,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			orderBookUID := args[0]
			particiapationIndex, err := cast.ToUint64E(args[1])
			if err != nil || particiapationIndex < 1 {
				return fmt.Errorf("participation number argument provided must be a non-negative-integer: %v", err)
			}

			params := &types.QueryParticipationExposureRequest{BookUid: orderBookUID, ParticipationIndex: particiapationIndex}
			res, err := queryClient.ParticipationExposure(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryHistoricalParticipationExposures implements the command to query all the historical participation exposures to a specific orderbook.
func GetCmdQueryHistoricalParticipationExposures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "historical-participation-exposures [order-book-id]",
		Short: "Query all historical participation exposures for a specific order book",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query historical participation exposures on an individual order book.

Example:
$ %s query orderbook historical-participation-exposures %s
`,
				version.AppName, "5531c60f-2025-48ce-ae79-1dc110f16000",
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			orderBookUID := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryHistoricalParticipationExposuresRequest{
				BookUid:    orderBookUID,
				Pagination: pageReq,
			}

			res, err := queryClient.HistoricalParticipationExposures(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "historical participation exposures")

	return cmd
}

// GetCmdQueryParticipationBets implements the command to query all the participation fulfilled bets to a specific orderbook.
func GetCmdQueryParticipationBets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "participation-bets [order-book-id] [participation-number]",
		Short: "Query all participation fulfilled bets for a specific order book",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query participation fulfilled bets on an individual order book.

Example:
$ %s query orderbook participation-bets %s %d
`,
				version.AppName, "5531c60f-2025-48ce-ae79-1dc110f16000", 2,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			orderBookUID := args[0]
			particiapationIndex, err := cast.ToUint64E(args[1])
			if err != nil || particiapationIndex < 1 {
				return fmt.Errorf("particiapnt number argument provided must be a non-negative-integer: %v", err)
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryParticipationFulfilledBetsRequest{
				BookUid:            orderBookUID,
				ParticipationIndex: particiapationIndex,
				Pagination:         pageReq,
			}

			res, err := queryClient.ParticipationFulfilledBets(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "participation fulfilled bets")

	return cmd
}
