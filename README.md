# PiCalculation

Calculates Pi in parallel in 3 different ways using Go.

### PiCalcArray

In this directory Pi is calculated using an array where each Go routine stores its result. The sum of the array is Pi.

### PiCalcSharedVar

In this directory Pi is calculated using a shared variable and a mutex. Each Go routine adds its result to that shared variable.

### PiCalcChannels

In this directory Pi is calculated using channels. Each Go routine sends its result in a channel. The main Go routine receives and sums
the results.
