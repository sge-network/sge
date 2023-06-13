# **Begin-Block**

Minting parameters are recalculated and inflation is paid at the beginning of each block.

## **Detect Current Phase**

Check if the chain has moved from one phase to another, including the final phase:

```go
// CurrentPhase returns current phase of the inflation
func (m Minter) CurrentPhase(params Params, currentBlock int64) (Phase, int) {
 if currentBlock == 1 {
  return params.GetPhaseAtStep(1), 1
 }

 cumulativeBlock := sdk.NewDec(0)
 var currentStep int
 var found bool

 // add each phase blocks until reaching the range which the current block is in
 for i := 0; i < len(params.Phases); i++ {
  // add current phase blocks to cummulative blocks
  cumulativeBlock = cumulativeBlock.Add(params.getPhaseBlocks(i + 1))

  currentStep = i + 1

  // if the current block is less than or equal to cummulative blocks
  // this means that we are in the i+1 step which is set in above line
  if sdk.NewDec(currentBlock).LTE(cumulativeBlock) {
   found = true
   // it is the current phase
   // so there is no need for furthur phase blocks check
   break
  }
 }

 // if there is no detected phase,
 // this means that the rest of inflation is zero as end phase
 if !found {
  return EndPhase(), EndPhaseAlias
 }

 // the phase has found and we need to return the phase specifications
 return params.GetPhaseAtStep(currentStep), currentStep
}
```

If the phase has changed, set the minter with the changed parameters:

```go
// set the new minter properties if the phase has changed or inflation has changed
 if int32(currentPhaseStep) != minter.PhaseStep || !minter.Inflation.Equal(currentPhase.Inflation) {

  // set new inflation rate
  newInflation := currentPhase.Inflation
  minter.Inflation = newInflation

  // set new phase step
  minter.PhaseStep = int32(currentPhaseStep)

  // set phase provisions of new phase step
  totalSupply := k.TokenSupply(ctx, params.MintDenom)
  minter.PhaseProvisions = minter.NextPhaseProvisions(totalSupply, params.ExcludeAmount, currentPhase)

  // store minter
  k.SetMinter(ctx, minter)

 }
```

---

## **Next Phase Provision**

Since the SGE-Network chain is predominantly reliant on phases for its inflation model, the phase_provision keeps track of the total amount of tokens to be distributed in the current phase. This value is calculated every time there is a phase change.

```go
// NextPhaseProvisions returns the phase provisions based on current total
// supply and inflation rate.
func (m Minter) NextPhaseProvisions(totalSupply sdk.Int, excludeAmount sdk.Int, phase Phase) sdk.Dec {
 // calculate annual provisions as normal
 annualProvisions := m.Inflation.MulInt(totalSupply.Sub(excludeAmount))

 // return this phase provisions according to the year coefficient
 // ex.
 //    year coefficient = 0.5
 //    blocks per year = 100
 //    this phase provisions is 100 * 0.5 => 50
 return annualProvisions.Mul(phase.YearCoefficient)
}
```

---

## **Block Provision**

Calculate the provisions generated for each block based on current phase provisions. The provisions are then minted by the `mint` module's ModuleMinterAccount and then transferred to the auth's FeeCollector ModuleAccount.

```go
// BlockProvisions returns the provisions for a block based on the phase
// provisions rate.
func (m Minter) BlockProvisions(params Params, phaseStep int) (sdk.Coin, sdk.Dec) {

 // get total blocks in this phase
 blocksPerPhase := params.getPhaseBlocks(phaseStep).TruncateDec()

 // detect each block provisions then and the truncated value from previous block
 provisionAmt := m.PhaseProvisions.Quo(blocksPerPhase).Add(m.TruncatedTokens)

 // extract the integer and decimal part of provisions
 // the decimal part is the truncated value because of conversion to sdk.Int
 // so the decimal part is truncated and needs to be added in next block
 intPart := provisionAmt.TruncateDec()
 decPart := provisionAmt.Sub(intPart)

 return sdk.NewCoin(params.MintDenom, intPart.TruncateInt()), decPart
}
```
