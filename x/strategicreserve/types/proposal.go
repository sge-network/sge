package types

import (
	fmt "fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeDataFeeCollectorFeed defines the type for a DataFeeCollectorFeedProposal
	ProposalTypeDataFeeCollectorFeed = "DataFeeCollectorFeed"
)

// Assert DataFeeCollectorFeedProposal implements govtypes.Content at compile-time
var _ govtypes.Content = &DataFeeCollectorFeedProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeDataFeeCollectorFeed)
	govtypes.RegisterProposalTypeCodec(&DataFeeCollectorFeedProposal{}, "strategicreserve/DataFeeCollectorFeedProposal")
}

// NewDataFeeCollectorFeedProposal creates a new data fee collector feed proposal.
//
//nolint:interfacer
func NewDataFeeCollectorFeedProposal(title, description string, houseFeeSpend, betFeeSpend sdk.Int) *DataFeeCollectorFeedProposal {
	return &DataFeeCollectorFeedProposal{title, description, houseFeeSpend, betFeeSpend}
}

// GetTitle returns the title of a data fee collector feed proposal.
func (dff *DataFeeCollectorFeedProposal) GetTitle() string { return dff.Title }

// GetDescription returns the description of a data fee collector feed proposal.
func (dff *DataFeeCollectorFeedProposal) GetDescription() string { return dff.Description }

// GetDescription returns the routing key of a data fee collector feed proposal.
func (dff *DataFeeCollectorFeedProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a data fee collector feed proposal.
func (dff *DataFeeCollectorFeedProposal) ProposalType() string {
	return ProposalTypeDataFeeCollectorFeed
}

// ValidateBasic runs basic stateless validity checks
func (dff *DataFeeCollectorFeedProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(dff)
	if err != nil {
		return err
	}
	if !dff.HouseFeeSpend.LT(sdk.ZeroInt()) {
		return ErrInvalidDataFeeCollectorProposalAmount
	}
	if !dff.BetFeeSpend.LT(sdk.ZeroInt()) {
		return ErrInvalidDataFeeCollectorProposalAmount
	}
	if dff.HouseFeeSpend.IsZero() && dff.BetFeeSpend.IsZero() {
		return ErrInvalidDataFeeCollectorProposalAmount
	}

	return nil
}

// String implements the Stringer interface.
func (dff DataFeeCollectorFeedProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Data Fee Collector Feed Proposal:
  Title:       %s
  Description: %s
  HouseFeeSpend:      %s
  BetFeeSpend:      %s
`, dff.Title, dff.Description, dff.HouseFeeSpend, dff.BetFeeSpend))
	return b.String()
}
