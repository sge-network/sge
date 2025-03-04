package types

import sdkerrors "cosmossdk.io/errors"

func (cf *PromoterConf) Validate() error {
	catMap := make(map[RewardCategory]struct{})
	for _, v := range cf.CategoryCap {
		_, ok := catMap[v.Category]
		if ok {
			return sdkerrors.Wrapf(ErrDuplicateCategoryInConf, "%s", v.Category)
		}
		if v.CapPerAcc <= 0 {
			return sdkerrors.Wrapf(ErrCategoryCapShouldBePos, "%s", v.Category)
		}

		catMap[v.Category] = struct{}{}
	}

	return nil
}
