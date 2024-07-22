package content

import "time"

type CredentialInfo struct {
	ID               int
	Resource         string
	Login            string
	Password         string
	ModificationTime time.Time
}

type BinaryFileInfo struct {
	ID               int
	BaseFileName     string
	FilePath         string
	Description      string
	ModificationTime time.Time
}

type TextInfo struct {
	ID               int
	Content          string
	Description      string
	ModificationTime time.Time
}

type CreditCardInfo struct {
	ID               int
	CardNumber       string
	CVV              string
	ValidThru        time.Time
	Bank             string
	ModificationTime time.Time
}
