package domain

type JWK struct {
	Kty string
	Use string
	Alg string
	Kid string
	N   string // RSA
	E   string // RSA
	Crv string // EC
	X   string // EC/OKP
	Y   string // EC
}

type JWKS struct {
	Keys []JWK
}
