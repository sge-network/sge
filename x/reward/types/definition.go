package types

import sdkerrors "cosmossdk.io/errors"

// ValidateBasic validates the basic properties of a reward definition.
func (d *Definition) ValidateBasic(blockTime uint64) error {
	if d.DstAccType != ReceiverAccType_RECEIVER_ACC_TYPE_SUB {
		if d.UnlockTS != 0 {
			return sdkerrors.Wrapf(ErrUnlockTSIsSubAccOnly, "%d", d.UnlockTS)
		}
	} else if d.UnlockTS <= blockTime {
		return sdkerrors.Wrapf(ErrUnlockTSDefBeforeBlockTime, "%d", d.UnlockTS)
	}
	return nil
}
