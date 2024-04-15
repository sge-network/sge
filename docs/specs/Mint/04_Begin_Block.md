# **Begin-Block**

Minting parameters are recalculated, and inflation is disbursed at the start of every block.

## **Detect Current Phase**

Verify whether the chain has transitioned between phases, including the ultimate phase:

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

If there has been a phase change, update the minter with the modified parameters:

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

Given that the SGE-Network chain heavily relies on phases for its inflation model, the phase_provision mechanism monitors the overall token distribution in the current phase. This value is recalculated whenever a phase transition occurs.

```go
// NextPhaseProvisions returns the phase provisions based on current total
// supply and inflation rate.
func (m Minter) NextPhaseProvisions(totalSupply sdkmath.Int, excludeAmount sdkmath.Int, phase Phase) sdk.Dec {
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

1. **Inflation Calculation**:
   - At the beginning of each block, the inflation parameters are recalculated.
   - The target annual inflation rate is adjusted based on the current bonded ratio (the ratio of bonded tokens to the total supply).
   - The inflation rate can change positively or negatively depending on how close it is to the desired ratio (usually 85%).
   - The maximum rate change per year is capped at 10%.
   - The annual inflation is then calculated, ensuring it falls between 5% and 10%.

2. **Annual Provisions**:
   - Next, we calculate the annual provisions based on the current total supply and the inflation rate.
   - This parameter is computed once per block.

3. **Block Provisions**:
   - The provisions generated for each block are derived from the annual provisions.
   - These provisions are minted by the **Mint module's ModuleMinterAccount**.
   - Finally, they are transferred to the **FeeCollector ModuleAccount** within the authentication system.

```go
// BlockProvisions returns the provisions for a block based on the phase
// provisions rate.
func (m Minter) BlockProvisions(params Params, phaseStep int) (sdk.Coin, sdk.Dec) {

 // get total blocks in this phase
 blocksPerPhase := params.getPhaseBlocks(phaseStep).TruncateDec()

 // detect each block provisions then and the truncated value from previous block
 provisionAmt := m.PhaseProvisions.Quo(blocksPerPhase).Add(m.TruncatedTokens)

 // extract the integer and decimal part of provisions
 // the decimal part is the truncated value because of conversion to sdkmath.Int
 // so the decimal part is truncated and needs to be added in next block
 intPart := provisionAmt.TruncateDec()
 decPart := provisionAmt.Sub(intPart)

 return sdk.NewCoin(params.MintDenom, intPart.TruncateInt()), decPart
}
```
