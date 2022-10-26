package app

// x/mint module sentinel errors
const (
	ErrTextVestingAccountStartBeforeEnd = "vesting start-time cannot be before end-time"
	ErrTextFailedToCreateAnteHandler    = "failed to create AnteHandler: %s"
)
