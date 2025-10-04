package domain

type JWK struct {
	KeyType   string `json:"kty"`           // Key type
	KeyID     string `json:"kid"`           // Unique identifier of this key.
	Use       string `json:"use,omitempty"` // [optional] Key assignment. "sig" | "enc"
	Algorithm string `json:"alg"`           // [optional] The algorithm for which the key is intended. "sha256"
	N         string // (Modulus): RSA public key component
	E         string // (Exponent): RSA public key component
	Crv       string // (Elliptic Curve) Curve type, such as P-256
	X         string // (X Coordinate): The X coordinate of the point on the curve.
	Y         string // (Y Coordinate): The Y coordinate of the point on the curve.
}

type JWKS struct {
	Keys []JWK
}
