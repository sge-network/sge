package types

import (
	cosmerrors "cosmossdk.io/errors"
)

func (a *Allocation) CheckExpiration(blockTime uint64) error {
	if blockTime > a.ExpTs {
		return cosmerrors.Wrapf(ErrAllocationExpired, "expire timestamp %d, block time %d", a.ExpTs, blockTime)
	}
	return nil
}
