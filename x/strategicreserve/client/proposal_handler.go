package client

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	"github.com/sge-network/sge/x/strategicreserve/client/cli"
)

// DataFeeCollectorFeedProposalHandler is the data fee collecttor feed handler.
var DataFeeCollectorFeedProposalHandler = govclient.NewProposalHandler(
	cli.CmdSubmitProposal,
	func(ctx client.Context) govrest.ProposalRESTHandler {
		return govrest.ProposalRESTHandler{
			SubRoute: "data-fee-collector-fund",
			Handler:  emptyHandler(ctx),
		}
	},
)

func emptyHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
