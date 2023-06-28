package types

// MarketBetConstraintsTest is a wrapper object for the bet constraints, It is being used
// to export unexported methods of the bet constraints type
type MarketBetConstraintsTest = MarketBetConstraints

func (bc *MarketBetConstraintsTest) Validate(params *Params) error {
	return bc.validate(params)
}
