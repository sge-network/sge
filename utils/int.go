package utils

// RoundToNearest rounds to the nearest multiplication of toNum
// ex.
// toRound = 397, step= 10 => 400
// toRound = 397, step= 5 => 395
func RoundToNearest(toRound, step int) int {

	remainder := toRound % step
	// if the remainder is zero this means it is already round
	if remainder == 0 {
		return toRound
	}

	// if the remainder is less than half of step
	if remainder < step/2 {
		// return lower bound
		return toRound - remainder
	} else {
		// return upper bound
		return (step - remainder) + toRound
	}
}
