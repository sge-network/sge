# **State**

## **Minter**

The minter serves as a repository for storing the present phase and inflation data.

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

The **minter** stores minting parameters in the global params store. Let's break down the key parameters:

1. **Mint Denom**: This parameter specifies the denomination of the minted token.

2. **Blocks Per Year**: It indicates how many blocks constitute one year. The phases, which are coefficients of years, depend on this parameter. For instance, if the `year_coefficient` is 0.5 and `blocks_per_year` is 1000, a phase will last for 500 blocks.

3. **Phases**: These represent discrete time frames during which a specific inflation rate applies.

   - **Duration**: The duration corresponds to the `year_coefficient`. It defines the phase's duration in years. For example, a `year_coefficient` of 0.75 means the phase lasts for 9 months (3/4th of a year).

   - **Inflation**: This parameter sets the inflation rate for the given phase. Inflation is expressed as a decimal; for instance, an inflation of 0.10000 corresponds to a 10% inflation rate.

4. **Exclude Amount**: This parameter specifies the number of tokens exempt from inflation.

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
