package storage

import "time"

type CredentialInfo struct {
	ID               int
	Resource         string
	Login            string
	Password         string
	ModificationTime time.Time
}

type FileInfo struct {
	ID               int
	BaseFileName     string
	ContentBase64    string
	Description      string
	ModificationTime time.Time
}

type CreditCardInfo struct {
	ID               int
	CardNumber       string
	CVV              string
	ValidThru        string
	Bank             string
	ModificationTime time.Time
}
