package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"reflect"
	"strconv"
	"strings"
)

const TypeMsgRewardUser = "reward_user"

var _ sdk.Msg = &MsgRewardUser{}

func NewMsgRewardUser(creator string, addresses string, amounts string, rType string, meta string, ticket string) (
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
	rewardTypeEnum, found := RewardType_value[rType]
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Invalid Reward type (%s)", rType, " Choose from values - ", reflect.ValueOf(RewardType_value).MapKeys())
	}
	rewards := Reward{
		Awardees:   getAwardeesList(addressList, amountList),
		RewardType: RewardType(rewardTypeEnum),
		Meta:       metaData,
	}
	return &MsgRewardUser{
		Creator: creator,
		Reward:  &rewards,
		Ticket:  ticket,
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

func validateAddress(address string, errorString string) error {
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, errorString, err)
	}
	return nil
}

func (msg *MsgRewardUser) validateAwardeeAddresses() error {
	for _, awardee := range msg.Reward.Awardees {
		err := validateAddress(awardee.GetAddress(), "invalid awardee address (%s)")
		if err != nil {
			return err
		}
	}
	return nil
}

func (msg *MsgRewardUser) validateAwardeeAmounts() error {
	for _, awardee := range msg.Reward.Awardees {
		if awardee.Amount <= 0 {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid amount: ", awardee.Amount)
		}
	}
	return nil
}

func (msg *MsgRewardUser) ValidateBasic() error {
	err := validateAddress(msg.Creator, "invalid creator address (%s)")
	if err != nil {
		return err
	}

	err = msg.validateAwardeeAddresses()
	if err != nil {
		return err
	}

	err = msg.validateAwardeeAmounts()
	if err != nil {
		return err
	}

	err = msg.validateAwardeeNumber()
	if err != nil {
		return err
	}
	return nil
}

func (msg *MsgRewardUser) validateAwardeeNumber() error {
	switch msg.Reward.RewardType {
	case RewardType_REFERRAL:
		if len(msg.Reward.Awardees) != 2 {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "For referral there should be 2 awardees")
		}
	default:
		return nil
	}
	return nil
}

func (msg *MsgRewardUser) ValidateSanity(ctx sdk.Context, p *Params) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	return nil
}
