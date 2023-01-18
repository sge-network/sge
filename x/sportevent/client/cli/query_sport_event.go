package cli

import (
	"context"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sge-network/sge/x/sportevent/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

// CmdListSportEvents implements a command to return all sport events
func CmdListSportEvents() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sport-events",
		Short: "list sport events",
		Long:  "Get list of sport events in paginated response.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QuerySportEventsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.SportEvents(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdShowSportEvent implements a command to return a specific sport events based on its UID
func CmdShowSportEvent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sport-event [uid]",
		Short: "get sport event",
		Long:  "Get sport event meta by uid.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argUID := args[0]

			params := &types.QuerySportEventRequest{
				Uid: argUID,
			}

			res, err := queryClient.SportEvent(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdSportEventListByUIDs returns command object for querying sport events by uid list
func CmdSportEventListByUIDs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sport-events-by-uids [uid]",
		Short: "Query sport events list by UIDs",
		Long:  "Get list of sport events by list of uids.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqUID := strings.Split(args[0], listSeparator)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QuerySportEventsByUIDsRequest{

				Uids: reqUID,
			}

			res, err := queryClient.SportEventsByUIDs(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
