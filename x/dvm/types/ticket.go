package types

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	gtime "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tendermint/tendermint/types/time"
)

// Ticket is the Interface of ticket.
type Ticket interface {
	// Unmarshal unmarshals the information of the ticket to the v. v must be a pointer.
	Unmarshal(v interface{}) error

	// Verify verifies the ticket signature with the given public keys. if the ticket verifies with any
	// of the keys, then return nil else return tyeps.ErrInvalidSignature
	Verify(pubKeys ...string) error

	// Consensus verifies that 2/3 of given public keys signed the ticket.
	Consensus(pubKeys ...string) error

	// IsValid verifies that the thicket is not expired yet.
	IsValid(ctx sdk.Context) error
}

// ticket is the Ticket implementer.
type ticket struct {
	value      string
	header     string
	payload    string
	signatures []string
	exp        time.WeightedTime
	clm        *jwt.RegisteredClaims
}

// NewTicket create a new Ticket from the given ticket.
func NewTicket(ticketStr string) (Ticket, error) {
	var err error
	t := ticket{
		value: ticketStr,
	}

	err = t.initFromValue()
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Unmarshal unmarshals the information of the ticket to the v. v must be a pointer.
func (t *ticket) Unmarshal(v interface{}) error {
	// data = json.Unmarshal(base64.Decode(payload))

	var err error
	bs, err := base64.RawURLEncoding.DecodeString(t.payload)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bs, v)
	if err != nil {
		return err
	}
	return nil
}

// Verify verifies the ticket signature with the given public keys.
// If the ticket verifies with any of the keys, then return nil else return tyeps.ErrInvalidSignature
func (t *ticket) Verify(pubKeys ...string) error {
	for _, v := range pubKeys {
		_, err := t.verifyWithKey(v)
		if err == nil {
			return nil
		}
		if err != ErrInvalidSignature {
			return err
		}
	}
	return ErrInvalidSignature
}

// Consensus verifies that 2/3 of given public keys signed the ticket.
// consensus mechanism will reside in this fuction if requested
func (t *ticket) Consensus(pubKeys ...string) error {
	if len(pubKeys) == 0 {
		return nil
	}
	for _, v := range pubKeys {
		_, err := t.verifyWithKey(v)
		if err == nil {
			return nil
		}
		if err != ErrInvalidSignature {
			return err
		}
	}
	return ErrInvalidSignature
}

func (t *ticket) IsValid(ctx sdk.Context) error {

	//validate the expiration
	if !t.exp.Time.After(ctx.BlockTime()) {
		return ErrTicketExpired
	}

	return nil
}

// initFromValue initializes the ticket from the raw value.few validation happening over this process.
func (t *ticket) initFromValue() error {
	var err error
	ts := strings.Split(t.value, JWTSeparator)
	if len(ts) < 3 {
		return ErrInvalidTicketFormat
	}
	t.header = ts[JWTHeaderIndex]
	t.payload = ts[JWTPayloadIndex]
	t.signatures = ts[JWTPayloadIndex+1:]

	clm := jwt.RegisteredClaims{}
	err = t.Unmarshal(&clm)
	if err != nil {
		return err
	}
	t.clm = &clm

	if t.clm.ExpiresAt == nil {
		return ErrExpirationRequired
	}
	gt := gtime.Unix(t.clm.ExpiresAt.Unix(), 0)
	t.exp = *time.NewWeightedTime(gt, DefaultTimeWeight)

	return nil
}

// verifyWithKey verify a Ticket with the key
// ?- is any Performance improvement possible?
func (t *ticket) verifyWithKey(key string) (bool, error) {
	for _, s := range t.signatures {
		token := t.header + "." + t.payload + "." + s
		parser := jwt.NewParser(
			jwt.WithoutClaimsValidation(),
		)
		T, err := parser.Parse(token, func(t *jwt.Token) (interface{}, error) {
			P, err := jwt.ParseEdPublicKeyFromPEM([]byte(key))
			if err != nil {
				return nil, err
			}
			return P, nil
		})

		if err != nil {
			return false, err
		}
		if T.Valid {
			return true, nil
		}

	}

	return false, ErrInvalidSignature
}
