package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/sge-network/sge/x/legacy/reward/types"
)

var DefaultRelativePacketTimeoutTimestamp = cast.ToUint64((time.Duration(10) * time.Minute).Nanoseconds())

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdCreatePromoter())
	cmd.AddCommand(CmdSetPromoterConf())
	cmd.AddCommand(CmdCreateCampaign())
	cmd.AddCommand(CmdUpdateCampaign())
	cmd.AddCommand(CmdWithdrawFunds())
	cmd.AddCommand(CmdGrantReward())

	return cmd
}
