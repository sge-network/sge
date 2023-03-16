# **Parameters**

The minting module contains the following parameters:

|      Key       |       Type       |  Example   |
|:-------------: |:---------------: |:---------: |
| MintDenom      | string           | "usge"     |
| BlocksPerYear  | string (uint64)  | "6311520"  |
| Phases         | Phases           | [{ "inflation": "0.101451460885956644", "year_coefficient": "0.500000000000000000" }]   |
| ExcludeAmount  | sdk.Int          | 100000     |
