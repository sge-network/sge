package types

import (
	"encoding/base64"
	"encoding/json"
	fmt "fmt"
	"strings"
	gtime "time"

	"github.com/cometbft/cometbft/types/time"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang-jwt/jwt/v4"
)

// jwtTicket is the Ticket implementer.
type jwtTicket struct {
	value     string
	header    string
	payload   string
	signature string
	exp       time.WeightedTime
	clm       *jwt.RegisteredClaims
}

// NewJwtTicket create a new jwt ticket from the given ticket.
func NewJwtTicket(ticketStr string) (Ticket, error) {
	var err error
	t := jwtTicket{
		value: ticketStr,
	}

	err = t.initFromValue()
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Unmarshal the information of the ticket to the v. v must be a pointer.
func (t *jwtTicket) Unmarshal(v interface{}) error {
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

// Verify verifies the ticket signature with the given public key. if the ticket is verified by
// the key, then return nil else return invalid signature error
func (t *jwtTicket) Verify(key string) error {
	_, err := t.verifyJwtKey(key)
	if err == nil {
		return nil
	}
	if err != ErrInvalidSignature {
		return err
	}

	return ErrInvalidSignature
}

// VerifyAny verifies the ticket signature with the given public keys. if the ticket is verified by
// the key, then return nil else return invalid signature error
func (t *jwtTicket) VerifyAny(keys []string) error {
	verified := false
	for _, pubKey := range keys {
		if err := t.Verify(pubKey); err == nil {
			// if no error, this means that ticked verified with this public key
			verified = true
			break
		}
	}

	if !verified {
		return ErrInvalidSignature
	}

	return nil
}

func (t *jwtTicket) ValidateExpiry(ctx sdk.Context) error {
	// validate the expiration
	if !t.exp.Time.After(ctx.BlockTime()) {
		return ErrTicketExpired
	}

	return nil
}

// initFromValue initializes the ticket from the raw value.few validation happening over this process.
func (t *jwtTicket) initFromValue() error {
	var err error
	ts := strings.Split(t.value, JWTSeparator)
	if len(ts) < 3 {
		return ErrInvalidTicketFormat
	}
	t.header = ts[JWTHeaderIndex]
	t.payload = ts[JWTPayloadIndex]
	t.signature = ts[JWTPayloadIndex+1]

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

// verifyJwtKey verify a Ticket with the key
func (t *jwtTicket) verifyJwtKey(key string) (bool, error) {
	token := t.header + "." + t.payload + "." + t.signature
	parser := jwt.NewParser(
		jwt.WithoutClaimsValidation(),
	)
	parsedToken, err := parser.Parse(token, func(_ *jwt.Token) (interface{}, error) {
		parsedPubKey, err := jwt.ParseEdPublicKeyFromPEM([]byte(key))
		if err != nil {
			return nil, err
		}
		return parsedPubKey, nil
	})
	if err != nil {
		return false, err
	}
	if parsedToken.Valid {
		return true, nil
	}

	return false, ErrInvalidSignature
}

func IsValidJwtToken(s string) error {
	if _, err := jwt.ParseEdPublicKeyFromPEM([]byte(s)); err != nil {
		return fmt.Errorf("public key %s is not valid jwt token %s", s, err)
	}
	return nil
}
