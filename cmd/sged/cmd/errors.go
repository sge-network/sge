package cmd

const (
	errTextVestingAccountGreaterThanTotal  = "vesting amount cannot be greater than total amount"
	errTextInvalidVestingParams            = "invalid vesting parameters; must supply start and end time or end time"
	errTextGenesisUnmarshalFailed          = "failed to unmarshal genesis state: %w"
	errTextBlankGenesisUnmarshalFailed     = "failed to marshal bank genesis state: %w"
	errTextApplicationGenesisMarshalFailed = "failed to marshal application genesis state: %w"
	errTextGettingAddressFromKeybaseFailed = "failed to get address from Keybase: %w"
	errTextCoinsParsingFailed              = "failed to parse coins: %w"
	errTextGetAccountFromAnyFailed         = "failed to get accounts from any: %w"
	errTextCannotAddExistingAddress        = "cannot add account at existing address %s"
	errTextConvertAccountToAnyFailed       = "failed to convert accounts into any's: %w"
	errTextGenesisAuthMarshalFailed        = "failed to marshal auth genesis state: %w"
	errTextParsingVestingAmountFailed      = "failed to parse vesting amount: %w"
	errTextGenesisAccValidationFailed      = "failed to validate new genesis account: %w"
)
