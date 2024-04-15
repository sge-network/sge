# **Concepts**

In contrast to the default **Mint Module** in Cosmos-SDK, which dynamically adjusts the inflation rate based on the ratio of bonded and unbonded tokens, the **Mint Module** of the SGE-Network chain adheres to a strict regime of inflation rates defined as **phases**.

These phases represent discrete time frames during which a specific inflation rate remains constant. Each phase has two key components:

1. **Duration**: The duration is measured by the **year coefficient**, representing the number of years the phase will last. For instance, a year coefficient of 0.75 corresponds to a phase lasting 9 months (3/4th of a year).

2. **Inflation**: This parameter specifies the inflation rate for the given phase. Inflation is expressed as a decimal; for example, an inflation value of 0.10000 corresponds to a 10% inflation rate.

Governance mechanisms allow adjustments to both the duration and inflation rate of these phases.

Furthermore, once all specified phases conclude, the chain enters a special phase known as the **final phase**. In this phase, the duration is infinite, and the inflation rate becomes zero. Notably, the inflation rate remains independent of the number of bonded and unbonded tokens.
