package domain

import "errors"

var (
	ErrInvalidClient        = errors.New("invalid_client")
	ErrUnauthorizedClient   = errors.New("unauthorized_client")
	ErrInvalidRedirectURI   = errors.New("invalid_redirect_uri")
	ErrInvalidScope         = errors.New("invalid_scope")
	ErrInvalidRequest       = errors.New("invalid_request")
	ErrInvalidGrant         = errors.New("invalid_grant")
	ErrInvalidCodeVerifier  = errors.New("invalid_code_verifier")
	ErrCodeAlreadyUsed      = errors.New("authorization_code_already_used")
	ErrCodeExpired          = errors.New("authorization_code_expired")
	ErrStateMismatch        = errors.New("state_mismatch")
	ErrNonceMismatch        = errors.New("nonce_mismatch")
	ErrTokenExpired         = errors.New("token_expired")
	ErrTokenRevoked         = errors.New("token_revoked")
	ErrSignatureInvalid     = errors.New("signature_invalid")
	ErrProviderVerification = errors.New("provider_verification_failed")
)
