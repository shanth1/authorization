package domain

import "time"

type Event interface {
	EventName() string
	OccurredAt() time.Time
}

type UserLoggedIn struct {
	Time      time.Time
	UserID    UserID
	ClientID  ClientID
	Provider  Provider
	SessionID SessionID
}

func (e UserLoggedIn) EventName() string     { return "UserLoggedIn" }
func (e UserLoggedIn) OccurredAt() time.Time { return e.Time }

type TokensIssued struct {
	Time           time.Time
	UserID         UserID
	ClientID       ClientID
	AccessTokenID  TokenID
	RefreshTokenID *TokenID
}

func (e TokensIssued) EventName() string     { return "TokensIssued" }
func (e TokensIssued) OccurredAt() time.Time { return e.Time }

type TokensRevoked struct {
	Time     time.Time
	ClientID ClientID
	UserID   *UserID
	Reason   string
}

func (e TokensRevoked) EventName() string     { return "TokensRevoked" }
func (e TokensRevoked) OccurredAt() time.Time { return e.Time }

type BackchannelLogoutSent struct {
	Time     time.Time
	ClientID ClientID
	UserID   UserID
	JTI      JTI
}

func (e BackchannelLogoutSent) EventName() string     { return "BackchannelLogoutSent" }
func (e BackchannelLogoutSent) OccurredAt() time.Time { return e.Time }
