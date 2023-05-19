package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRewardUser = "reward_user"

var _ sdk.Msg = &MsgRewardUser{}

func NewMsgRewardUser(creator string) *MsgRewardUser {
	return &MsgRewardUser{
		Creator: creator,
	}
}

func (msg *MsgRewardUser) Route() string {
	return RouterKey
}

func (msg *MsgRewardUser) Type() string {
	return TypeMsgRewardUser
}

func (msg *MsgRewardUser) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRewardUser) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRewardUser) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func (msg *MsgRewardUser) ValidateSanity(ctx sdk.Context, p *Params) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	// TODO Check for total awardees, their sum limit
	return nil
}
