package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/sge-network/sge/x/strategicreserve/types"
	"github.com/spf13/cobra"
)

// CmdSubmitProposal implements the command to submit a ata-fee-collector-fund proposal
func CmdSubmitProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data-fee-collector-fund [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a data fee collector fund proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a data fee collector fund along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal data-fee-collector-fund <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Data Fee Collector Fund",
  "description": "Pay me some Atoms!",
  "house_fee_spend": "1000",
  "bet_fee_spend": "100",
  "deposit": "1000usge"
}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := parseDataFeeCollectorProposalWithDeposit(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			content := types.NewDataFeeCollectorFeedProposal(proposal.Title, proposal.Description, proposal.HouseFeeSpend, proposal.BetFeeSpend)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
