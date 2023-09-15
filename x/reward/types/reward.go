package types

import (
	"errors"

	cosmerrors "cosmossdk.io/errors"
)

// IRewardFactory defines the methods that should be implemented for all of reward types.
type IRewardFactory interface {
	ValidateBasic(campaign Campaign) error
	VaidateDefinitions(campaign Campaign) error
	CalculateDistributions(definitions []Definition, ticket string) ([]Distribution, error)
}

// SignUpReward is the type for signup rewards calculations
type SignUpReward struct{}

// NewSignUpReward create new object of signup reward calculator type.
func NewSignUpReward() SignUpReward { return SignUpReward{} }

// ValidateBasic validates basic and common campaign constraints.
func (sur SignUpReward) ValidateBasic(campaign Campaign) error {
	return errors.New("not implemented")
}

// VaidateDefinitions validates campaign definitions.
func (sur SignUpReward) VaidateDefinitions(campaign Campaign) error {
	if len(campaign.RewardDefs) != 1 {
		return cosmerrors.Wrapf(ErrWrongDefinitionsCount, "signup rewards can only have single definition")
	}
	if campaign.RewardDefs[0].ReceiverType != ReceiverType_RECEIVER_TYPE_SINGLE {
		return cosmerrors.Wrapf(ErrInvalidReceiverType, "signup rewards can be defined for subaccount only")
	}
	return nil
}

// CalculateDistributions parses ticket payload and returns the distribution list of signup reward.
func (sur SignUpReward) CalculateDistributions(definitions []Definition, ticket string) ([]Distribution, error) {
	return []Distribution{}, errors.New("not implemented")
}

// ReferralReward is the type for referral rewards calculations
type ReferralReward struct{}

// NewReferralReward create new object of referral reward calculator type.
func NewReferralReward() ReferralReward { return ReferralReward{} }

// ValidateBasic validates basic and common campaign constraints.
func (rfr ReferralReward) ValidateBasic(campaign Campaign) error {
	return errors.New("not implemented")
}

// VaidateDefinitions validates campaign definitions.
func (rfr ReferralReward) VaidateDefinitions(campaign Campaign) error {
	hasReferrer := false
	hasReferee := false
	for _, d := range campaign.RewardDefs {
		switch d.ReceiverType {
		case ReceiverType_RECEIVER_TYPE_REFEREE:
			hasReferee = true
		case ReceiverType_RECEIVER_TYPE_REFERRER:
			hasReferrer = true
		default:
			return cosmerrors.Wrapf(ErrInvalidReceiverType, "%s", d.ReceiverType)
		}
	}

	if !hasReferee || !hasReferrer {
		return cosmerrors.Wrapf(ErrMissingDefinition, "referral rewards should have the referrer and the referee")
	}
	return nil
}

// CalculateDistributions parses ticket payload and returns the distribution list of referral reward.
func (rfr ReferralReward) CalculateDistributions(definitions []Definition, ticket string) ([]Distribution, error) {
	return []Distribution{}, errors.New("not implemented")
}

// AffiliationReward is the type for affiliation rewards calculations
type AffiliationReward struct{}

// NewAffiliationReward create new object of affiliation reward calculator type.
func NewAffiliationReward() AffiliationReward { return AffiliationReward{} }

// ValidateBasic validates basic and common campaign constraints.
func (afr AffiliationReward) ValidateBasic(campaign Campaign) error {
	return errors.New("not implemented")
}

// VaidateDefinitions validates campaign definitions.
func (afr AffiliationReward) VaidateDefinitions(campaign Campaign) error {
	if len(campaign.RewardDefs) != 1 {
		return cosmerrors.Wrapf(ErrWrongDefinitionsCount, "affiliation rewards can only have single definition")
	}
	if campaign.RewardDefs[0].ReceiverType != ReceiverType_RECEIVER_TYPE_SINGLE {
		return cosmerrors.Wrapf(ErrInvalidReceiverType, "affiliation rewards can be defined for subaccount only")
	}
	return nil
}

// CalculateDistributions parses ticket payload and returns the distribution list of affiliation reward.
func (afr AffiliationReward) CalculateDistributions(definitions []Definition, ticket string) ([]Distribution, error) {
	return []Distribution{}, errors.New("not implemented")
}

// NoLossBetsReward is the type for no loss bets rewards calculations
type NoLossBetsReward struct{}

// NewNoLossBetsReward create new object of no loss bets reward calculator type.
func NewNoLossBetsReward() NoLossBetsReward { return NoLossBetsReward{} }

// ValidateBasic validates basic and common campaign constraints.
func (afr NoLossBetsReward) ValidateBasic(campaign Campaign) error {
	return errors.New("not implemented")
}

// VaidateDefinitions validates campaign definitions.
func (afr NoLossBetsReward) VaidateDefinitions(campaign Campaign) error {
	if len(campaign.RewardDefs) != 1 {
		return cosmerrors.Wrapf(ErrWrongDefinitionsCount, "noloss bets rewards can only have single definition")
	}
	if campaign.RewardDefs[0].ReceiverType != ReceiverType_RECEIVER_TYPE_SINGLE {
		return cosmerrors.Wrapf(ErrInvalidReceiverType, "noloss bets rewards can be defined for subaccount only")
	}
	return nil
}

// CalculateDistributions parses ticket payload and returns the distribution list of no loss bets reward.
func (afr NoLossBetsReward) CalculateDistributions(definitions []Definition, ticket string) ([]Distribution, error) {
	return []Distribution{}, errors.New("not implemented")
}
