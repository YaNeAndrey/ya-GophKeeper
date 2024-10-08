package content

import (
	"fmt"
	"time"
)

type CredentialInfo struct {
	ID               int       `json:"id"`
	Resource         string    `json:"resource"`
	Login            string    `json:"login"`
	Password         string    `json:"passwd"`
	ModificationTime time.Time `json:"mod_time"`
}

func (c *CredentialInfo) String() string {
	resStr := fmt.Sprintf("Resource: %s; Login: %s; Password: %s;\r\n", c.Resource, c.Login, c.Password)
	return resStr
}

type BinaryFileInfo struct {
	ID               int       `json:"id"`
	FileName         string    `json:"file_name"`
	FilePath         string    `json:"file_path"`
	Description      string    `json:"description"`
	FileSize         int       `json:"file_size"`
	MD5              string    `json:"md5"`
	ModificationTime time.Time `json:"mod_time"`
}

func (b *BinaryFileInfo) String() string {
	resStr := fmt.Sprintf("File name: %s; File path: %s; Description: %s;\r\n", b.FileName, b.FilePath, b.Description)
	return resStr
}

type TextInfo struct {
	ID               int       `json:"id"`
	Content          string    `json:"content"`
	Description      string    `json:"description"`
	ModificationTime time.Time `json:"mod_time"`
}

func (t *TextInfo) String() string {
	resStr := fmt.Sprintf("Description: %s;\r\n Content: %s;\r\n", t.Description, t.Content)
	return resStr
}

type CreditCardInfo struct {
	ID               int       `json:"id"`
	CardNumber       string    `json:"card_number"`
	CVV              string    `json:"cvv"`
	ValidThru        time.Time `json:"valid_thru"`
	Bank             string    `json:"bank"`
	ModificationTime time.Time `json:"mod_time"`
}

func (c *CreditCardInfo) String() string {
	resStr := fmt.Sprintf("Bank: %s;CardNumber: %s;CVV: %s;ValidThru: %s;\r\n", c.Bank, c.CardNumber, c.CVV, c.ValidThru)
	return resStr
}
