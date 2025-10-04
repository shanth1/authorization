package domain

import "time"

type Email struct {
	Address    string     `json:"address"`
	Verified   bool       `json:"verified"`
	VerifiedAt *time.Time `json:"verified_at"`
	Primary    bool       `json:"primary"`
}
