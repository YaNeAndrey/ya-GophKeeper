package content

import (
	"fmt"
	"time"
)

type CredentialInfo struct {
	ID               int
	Resource         string
	Login            string
	Password         string
	ModificationTime time.Time
}

func (c *CredentialInfo) String() string {
	resStr := fmt.Sprintf("Resource: %s; Login: %s; Password: %s;\r\n", c.Resource, c.Login, c.Password)
	return resStr
}

type BinaryFileInfo struct {
	ID               int
	BaseFileName     string
	FilePath         string
	Description      string
	ModificationTime time.Time
}

func (b *BinaryFileInfo) String() string {
	resStr := fmt.Sprintf("File name: %s; File path: %s; Description: %s;\r\n", b.BaseFileName, b.FilePath, b.Description)
	return resStr
}

type TextInfo struct {
	ID               int
	Content          string
	Description      string
	ModificationTime time.Time
}

func (t *TextInfo) String() string {
	resStr := fmt.Sprintf("Description: %s; Content: %s;\r\n", t.Description, t.Content)
	return resStr
}

type CreditCardInfo struct {
	ID               int
	CardNumber       string
	CVV              string
	ValidThru        time.Time
	Bank             string
	ModificationTime time.Time
}

func (c *CreditCardInfo) String() string {
	resStr := fmt.Sprintf("Bank: %s;CardNumber: %s;CVV: %s;ValidThru: %s;\r\n", c.Bank, c.CardNumber, c.CVV, c.ValidThru)
	return resStr
}
