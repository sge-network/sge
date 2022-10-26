package types

// NewReserver returns a new reserver and initializes it with
// the passed values
func NewReserver(srpool *SRPool) Reserver {
	return Reserver{
		SrPool: srpool,
	}
}

// InitialReserver returns an initial Reserver object.
func InitialReserver() Reserver {
	return NewReserver(
		&InitialSrPool,
	)
}

// DefaultInitialReserver returns the default reserver
func DefaultInitialReserver() Reserver {
	return InitialReserver()
}

// ValidateReserver validates the reserver
func ValidateReserver(reserver Reserver) error {
	return nil
}
