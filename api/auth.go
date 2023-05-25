package api

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"git.sr.ht/~rottenfishbone/go-cook/pkg/config"
)

// An JSON web token for persistent authentication from a client
type JWT struct {
	Header  JWTHeader
	Payload JWTPayload

	Signature  string // Base64 URLEncoded checksum of all input data
	HeaderStr  string // Base64 URLEncoded JSON of Header
	PayloadStr string // Base64 URLEncoded JSON of Payload
}
type JWTHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}
type JWTPayload struct {
	Username string `json:"sub"` // Username of user issued to
	Issued   int64  `json:"iat"` // Exact time of issue (seconds since epoch)
	Expires  int64  `json:"exp"` // Time of expiration (seconds since epoch)
}

var (
	ErrJWTMalformed   = errors.New("Malformed JWT Json.")
	ErrJWTFutureDated = errors.New("JWT issue time is in the future.")
	ErrJWTExpired     = errors.New("JWT is expired.")
	ErrJWTMismatch    = errors.New("JWT signature mismatch.")
)

// Decodes a JSON string into a JWT object
//
// # Returns token, nil on success
//
// Invalid JSON or a malformed JWT string returns an error.
func DecodeJWTFromJSON(jsonBytes []byte) (JWT, error) {
	var err error

	// Validate and separate the JWT parts
	jwtRegex := `^\"?([A-Za-z0-9_\-]+)\.([A-Za-z0-9_\-]+)\.([A-Za-z0-9_\-]+)\"?$`
	regex := regexp.MustCompile(jwtRegex)
	matches := regex.FindSubmatch(jsonBytes)
	if len(matches) != 4 {
		return JWT{}, ErrJWTMalformed
	}
	headerEnc := string(matches[1])
	payloadEnc := string(matches[2])
	signatureEnc := string(matches[3])

	// Decode payload and header into JSON
	headerJson, _ := base64.RawURLEncoding.DecodeString(headerEnc)
	payloadJson, _ := base64.RawURLEncoding.DecodeString(payloadEnc)

	// Construct the JWT as we have so far, allocating space for Header/Payload
	jwt := JWT{
		Header:     JWTHeader{},
		Payload:    JWTPayload{},
		Signature:  signatureEnc,
		HeaderStr:  headerEnc,
		PayloadStr: payloadEnc,
	}

	// Unmarshal the header and payload into `jwt`
	if err = json.Unmarshal(headerJson, &jwt.Header); err != nil {
		return JWT{}, err
	}
	if err = json.Unmarshal(payloadJson, &jwt.Payload); err != nil {
		return JWT{}, err
	}

	return jwt, nil
}

// Encodes a `JWT` struct as a standard JWT string
func EncodeJWTToString(token JWT) string {
	return fmt.Sprintf("%v.%v.%v",
		token.HeaderStr, token.PayloadStr, token.Signature)
}

// Builds and signs a new JWT for a user.
//
// # Returns `token, nil` on success
//
// The token should not be used on error
func GenerateJWT(username string, issued int64, expiry int64) (JWT, error) {
	var err error

	jwt := JWT{
		Header: JWTHeader{
			Algorithm: "HS512",
			Type:      "JWT",
		},
		Payload: JWTPayload{
			Username: username,
			Issued:   issued,
			Expires:  expiry,
		},
	}

	// Encode header and payload as base64 strings
	var headerJson, payloadJson []byte
	if headerJson, err = json.Marshal(jwt.Header); err != nil {
		return JWT{}, err
	}
	if payloadJson, err = json.Marshal(jwt.Payload); err != nil {
		return JWT{}, err
	}
	headerString := base64.RawURLEncoding.EncodeToString(headerJson)
	jwt.HeaderStr = headerString
	payloadString := base64.RawURLEncoding.EncodeToString(payloadJson)
	jwt.PayloadStr = payloadString

	// Use HMAC-SHA512 to sign the results
	hashFunc := hmac.New(sha512.New, config.GetHMACKeyBytes())
	jwtToHash := fmt.Sprintf("%v.%v", headerString, payloadString)
	signatureBytes := hashFunc.Sum([]byte(jwtToHash))

	// Encode the signature
	signatureString := base64.RawURLEncoding.EncodeToString(signatureBytes)
	jwt.Signature = signatureString

	return jwt, nil
}

// Builds and signs a new JWT for a user with a specified time to live.
//
// # Returns `jsonBytes, nil` on success
//
// The resultant JSON is a standard JWT string (`.` delimited, base64url encoded)
//
// `ttl` is in seconds
func GenerateJWTJSON(username string, ttl int64) ([]byte, error) {
	var err error

	now := time.Now().Unix()
	var jwt JWT
	if jwt, err = GenerateJWT(username, now, now+ttl); err != nil {
		return nil, err
	}

	jwtStr := EncodeJWTToString(jwt)
	return []byte(jwtStr), nil
}

// Validates a JWT with this server's key
//
// # Returns nil if the token is current and the signature is valid
//
// Returns one of the following errors on failure which should likely be handled
// by case:
//   - ErrJWTFutureDated	- `Issued` is in future
//   - ErrJWTExpired		- `Expiry` has passed
//   - ErrJWTMismatch		- `Signature` does not match generated one
//   - parsing/decoding errors
func ValidateJWT(token JWT) error {
	var err error
	// Time validation
	now := time.Now().Unix()
	data := token.Payload
	if data.Issued > now {
		return ErrJWTFutureDated
	}
	if data.Expires < now {
		return ErrJWTExpired
	}

	// Build a JWT to compare against
	compJwt, err := GenerateJWT(data.Username, data.Issued, data.Expires)
	if err != nil {
		return err
	}

	// Decode Base64 signatures
	compSig, err := base64.RawURLEncoding.DecodeString(compJwt.Signature)
	if err != nil {
		return err
	}
	passedSig, err := base64.RawURLEncoding.DecodeString(token.Signature)
	if err != nil {
		return err
	}

	// Compare signatures (using hmac.Equal to avoid timing attacks)
	if !hmac.Equal(compSig, passedSig) {
		return ErrJWTMismatch
	}

	return nil
}

// Validates a JWT with the servers key.
//
// # Returns nil on success
//
// See `ValidateJWT` for more details
func ValidateJWTJSON(jsonBytes []byte) error {
	var err error
	jwt, err := DecodeJWTFromJSON(jsonBytes)
	if err != nil {
		return err
	}

	return ValidateJWT(jwt)
}
