package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
	"strings"
)

const TypeMsgRewardUser = "reward_user"

var _ sdk.Msg = &MsgRewardUser{}

func NewMsgRewardUser(creator string, addresses string, amounts string, rType string, meta string, incentiveId string) (
	*MsgRewardUser, error) {
	addressList := strings.Split(addresses, ",")
	amountList, err := mapStringsToInts(strings.Split(amounts, ","))
	if err != nil {
		return nil, err
	}
	if len(addressList) != len(amountList) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "unequal amounts and addresses (%s) (%s)", addressList, amountList)
	}
	var metaData map[string]string
	err = json.Unmarshal([]byte(meta), &metaData)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Invalid Meta (%s)", meta)
	}
	rewards := Reward{
		Awardees:    getAwardeesList(addressList, amountList),
		RewardType:  RewardType(RewardType_value[rType]),
		Meta:        metaData,
		IncentiveId: incentiveId,
	}
	return &MsgRewardUser{
		Creator: creator,
		Reward:  &rewards,
	}, nil
}

func getAwardeesList(addressList []string, amountList []int) []*Awardee {
	var awardees []*Awardee
	for i := 0; i < len(addressList); i++ {
		awardees = append(awardees, &Awardee{
			Address: addressList[i],
			Amount:  uint64(amountList[i]),
		})
	}
	return awardees
}

func mapStringsToInts(strList []string) ([]int, error) {
	var err error
	intList := make([]int, len(strList))
	for i, str := range strList {
		intList[i], err = strconv.Atoi(str)
		if err != nil {
			return intList, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "invalid integer string (%s)", err)
		}
	}
	return intList, nil
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
