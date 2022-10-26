# **Concepts**

Unlike the default Mint Module of Cosmos which can vary the inflation rate dynamically based on the ratio of bonded and unbonded tokens at any point in time, the Mint module of the Sge Network chain follows a strict regime of inflation rates defined as phases

Phases are nothing but specific discreet time frames during which a certain inflation rate holds. Phases have two components:

- duration: The duration is defined as the year_coefficient. It defines the time in years for which a phase will hold. For example a yearcoefficient of 0.75 means that the phase will last for 9 months, that is, 3/4th of a year.

- inflation: This parameter defines the inflation rate of the chain for the phase in question. Inflation is defined as a decimal. That is, inflation of 0.10000 means an inflation rate of 10%.

The duration and inflation rate of phases can be modified via governance.

> If all the specified phases are over, the chain enters a special phase called the final_phase, where the phase duration is infinite and the phase inflation is zero.

> Note that the inflation rate does not depend on the number of bonded and unbonded tokens
