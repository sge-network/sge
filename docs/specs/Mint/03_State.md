# **State**

## **Minter**

The minter is a space for holding the current phase and inflation information.

```proto
// Minter represents the minting state.
message Minter {
  // current annual inflation rate
  string inflation = 1 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  // phase step is the index of phases slice + 1
  int32 phase_step = 2;
  // current phase expected provisions
  string phase_provisions = 3 [
    (gogoproto.moretags)   = "yaml:\"phase_provisions\"",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable)   = false
  ];
  // returns current truncated tokens because of Dec to Int conversion in the minting
  string truncated_tokens = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];

}
```

---

## **Params**

Minting params are held in the global params store.

1. Mint Denom: This parameter represents the denomination of the minted token

2. Blocks Per Year: This parameter is used to signify how many blocks will comprise 1 year. Note that the phases, which are defined as a coefficient of years, will also be based on this parameter. For example, if the year_coefficient is defined as 0.5, and the blocks_per_year is defined as 1000, that means the phase will last for 500 blocks.

3. Phases: Phases are defined as specific discreet time frames during which a certain inflation rate holds.

   a. duration: The duration is defined as the year_coefficient. It defines the time in years for which a phase will hold. For example a yearcoefficient of 0.75 means that the phase will last for 9 months, that is, 3/4th of a year.

   b. inflation: This parameter defines the inflation rate of the chain for the phase in question. Inflation is defined as a decimal. That is, inflation of 0.10000 means an inflation rate of 10%.

4. Exclude Amount: This parameter defines the number of tokens that will not incur inflation.

```proto
// Params defines the parameters for the module.
message Params {
  option (gogoproto.goproto_stringer) = false;

  // type of coin to mint
  string mint_denom = 1;
  // expected blocks per year
  int64 blocks_per_year = 2 [(gogoproto.moretags) = "yaml:\"blocks_per_year\""];

  // phases
  repeated Phase phases = 3 [(gogoproto.moretags) = "yaml:\"phases\"", (gogoproto.nullable) = false];

  string exclude_amount = 4 [
    (gogoproto.moretags)   = "yaml:\"exclude_amount\"",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable)   = false
  ];
}
```

```proto
// Params defines the phase parameters for the module.
message Phase {
  option (gogoproto.goproto_stringer) = false;

  // the phase inflation rate
  string inflation = 1  [
    (gogoproto.moretags) = "yaml:\"inflation\"",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable)   = false];
  // total coefficient from the beginning
  string year_coefficient = 2 [
    (gogoproto.moretags) = "yaml:\"year_coefficient\"",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable)   = false];

}
```
