package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"
)

const (
	flagVestingStart     = "vesting-start-time"
	flagVestingEnd       = "vesting-end-time"
	flagVestingAmt       = "vesting-amount"
	genesisAccountNumber = 0
	genesisSquenceNumber = 0
	defaultVestingStart  = 0
	defaultVestingEnd    = 0
)

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddGenesisAccountCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address_or_key_name] [coin][,[coin]]",
		Short: "Add a genesis account to genesis.json",
		Long: `Add a genesis account to genesis.json. The provided account must specify
the account address or key name and a list of initial coins. If a key name is given,
the address will be looked up in the local Keybase. The list of initial tokens must
contain valid denominations. Accounts may optionally be supplied with vesting parameters.
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			addr, genAccount, balances, err := getAccountAddressAndBalances(cmd, args[0], args[1], clientCtx.HomeDir)
			if err != nil {
				return err
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf(errTextGenesisUnmarshalFailed, err)
			}

			authGenStateBz, err := getAuthGenesisState(clientCtx.Codec, appState, addr, genAccount)
			if err != nil {
				return err
			}

			appState[authtypes.ModuleName] = authGenStateBz

			bankGenState := banktypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			bankGenState.Balances = append(bankGenState.Balances, balances)
			bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)
			bankGenState.Supply = bankGenState.Supply.Add(balances.Coins...)

			bankGenStateBz, err := clientCtx.Codec.MarshalJSON(bankGenState)
			if err != nil {
				return fmt.Errorf(errTextBlankGenesisUnmarshalFailed, err)
			}

			appState[banktypes.ModuleName] = bankGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf(errTextApplicationGenesisMarshalFailed, err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, argAppHome)
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, argKeyringBackend)
	cmd.Flags().String(flagVestingAmt, "", argVestingaccountCoins)
	cmd.Flags().Int64(flagVestingStart, defaultVestingStart, argVestingScheduleStart)
	cmd.Flags().Int64(flagVestingEnd, defaultVestingEnd, argVestingScheduleEnd)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func getAccountAddress(cmd *cobra.Command, address string, homeDir string) (sdk.AccAddress, error) {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		inBuf := bufio.NewReader(cmd.InOrStdin())
		keyringBackend, err := cmd.Flags().GetString(flags.FlagKeyringBackend)
		if err != nil {
			return sdk.AccAddress{}, err
		}

		// attempt to lookup address from Keybase if no address was provided
		kb, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, homeDir, inBuf)
		if err != nil {
			return sdk.AccAddress{}, err
		}

		info, err := kb.Key(address)
		if err != nil {
			return sdk.AccAddress{}, fmt.Errorf(errTextGettingAddressFromKeybaseFailed, err)
		}

		addr = info.GetAddress()
	}
	return addr, nil
}

func getAccountAddressAndBalances(cmd *cobra.Command,
	addRessOrKeyName,
	conins,
	homeDir string,
) (
	sdk.AccAddress,
	authtypes.GenesisAccount,
	banktypes.Balance,
	error,
) {
	addr, err := getAccountAddress(cmd, addRessOrKeyName, homeDir)
	if err != nil {
		return sdk.AccAddress{}, nil, banktypes.Balance{}, err
	}

	coins, err := sdk.ParseCoinsNormalized(conins)
	if err != nil {
		return sdk.AccAddress{}, nil, banktypes.Balance{}, fmt.Errorf(errTextCoinsParsingFailed, err)
	}

	balances := banktypes.Balance{Address: addr.String(), Coins: coins.Sort()}
	genAccount, err := getAndValidateGeneisAccount(cmd, addr, balances)
	if err != nil {
		return sdk.AccAddress{}, nil, banktypes.Balance{}, err
	}
	return addr, genAccount, balances, err
}

func getAuthGenesisState(cdc codec.Codec,
	appState map[string]json.RawMessage,
	addr sdk.Address,
	genAccount authtypes.GenesisAccount,
) ([]byte, error) {
	authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)

	accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
	if err != nil {
		return nil, fmt.Errorf(errTextGetAccountFromAnyFailed, err)
	}

	if accs.Contains(addr) {
		return nil, fmt.Errorf(errTextCannotAddExistingAddress, addr)
	}

	// Add the new account to the set of genesis accounts and sanitize the
	// accounts afterwards.
	accs = append(accs, genAccount)
	accs = authtypes.SanitizeGenesisAccounts(accs)

	genAccs, err := authtypes.PackAccounts(accs)
	if err != nil {
		return nil, fmt.Errorf(errTextConvertAccountToAnyFailed, err)
	}
	authGenState.Accounts = genAccs

	authGenStateBz, err := cdc.MarshalJSON(&authGenState)
	if err != nil {
		return nil, fmt.Errorf(errTextGenesisAuthMarshalFailed, err)
	}
	return authGenStateBz, nil
}

func getVestingParams(cmd *cobra.Command) (int64, int64, sdk.Coins, error) {
	vestingStart, err := cmd.Flags().GetInt64(flagVestingStart)
	if err != nil {
		return 0, 0, sdk.Coins{}, err
	}
	vestingEnd, err := cmd.Flags().GetInt64(flagVestingEnd)
	if err != nil {
		return 0, 0, sdk.Coins{}, err
	}
	vestingAmtStr, err := cmd.Flags().GetString(flagVestingAmt)
	if err != nil {
		return 0, 0, sdk.Coins{}, err
	}

	vestingAmt, err := sdk.ParseCoinsNormalized(vestingAmtStr)
	if err != nil {
		return 0, 0, sdk.Coins{}, fmt.Errorf(errTextParsingVestingAmountFailed, err)
	}
	return vestingStart, vestingEnd, vestingAmt, nil
}

func getAndValidateGeneisAccount(cmd *cobra.Command, addr sdk.AccAddress, balances banktypes.Balance) (authtypes.GenesisAccount, error) {
	// create concrete account type based on input parameters
	var genAccount authtypes.GenesisAccount
	baseAccount := authtypes.NewBaseAccount(addr, nil, genesisAccountNumber, genesisSquenceNumber)
	vestingStart, vestingEnd, vestingAmt, err := getVestingParams(cmd)
	if err != nil {
		return nil, err
	}

	if !vestingAmt.IsZero() {
		baseVestingAccount := authvesting.NewBaseVestingAccount(baseAccount, vestingAmt.Sort(), vestingEnd)

		if (balances.Coins.IsZero() && !baseVestingAccount.OriginalVesting.IsZero()) ||
			baseVestingAccount.OriginalVesting.IsAnyGT(balances.Coins) {
			return nil, errors.New(errTextVestingAccountGreaterThanTotal)
		}

		switch {
		case vestingStart != defaultVestingStart && vestingEnd != defaultVestingEnd:
			genAccount = authvesting.NewContinuousVestingAccountRaw(baseVestingAccount, vestingStart)

		case vestingEnd != defaultVestingEnd:
			genAccount = authvesting.NewDelayedVestingAccountRaw(baseVestingAccount)

		default:
			return nil, errors.New(errTextInvalidVestingParams)
		}
	} else {
		genAccount = baseAccount
	}

	if err := genAccount.Validate(); err != nil {
		return nil, fmt.Errorf(errTextGenesisAccValidationFailed, err)
	}

	return genAccount, nil
}
