package cli

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

// parseDataFeeCollectorProposalWithDeposit reads and parses a parseDataFeeCollectorProposalWithDeposit from a file.
func parseDataFeeCollectorProposalWithDeposit(cdc codec.JSONCodec, proposalFile string) (types.DataFeeCollectorFeedProposalWithDeposit, error) {
	proposal := types.DataFeeCollectorFeedProposalWithDeposit{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
