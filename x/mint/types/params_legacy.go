package types

import paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

// parameter store keys
var (
	// KeyMintDenom is the mint denom param key
	keyMintDenom = []byte("MintDenom")
	// KeyBlocksPerYear is the blocks per year param key
	keyBlocksPerYear = []byte("BlocksPerYear")
	// KeyPhases is the inflation phases param key
	keyPhases = []byte("Phases")
	// KeyExcludeAmount is the excluded amount from inflation calculation param key
	keyExcludeAmount = []byte("ExcludeAmount")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(keyMintDenom, &p.MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(keyBlocksPerYear, &p.BlocksPerYear, validateBlocksPerYear),
		paramtypes.NewParamSetPair(keyPhases, &p.Phases, validatePhases),
		paramtypes.NewParamSetPair(keyExcludeAmount, &p.ExcludeAmount, validateExcludeAmount),
	}
}
