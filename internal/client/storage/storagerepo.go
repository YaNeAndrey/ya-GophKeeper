package storage

import (
	"ya-GophKeeper/internal/content"
)

type StorageRepo interface {
	AddNewCreditCard(creditCard *content.CreditCardInfo) error
	AddNewCredential(credential *content.CredentialInfo) error
	AddNewFile(file *content.BinaryFileInfo) error
	AddNewText(text *content.TextInfo) error

	RemoveCreditCard(index int) error
	RemoveCredential(index int) error
	RemoveText(index int) error
	RemoveFile(index int) error

	UpdateCreditCards(index int, creditCard *content.CreditCardInfo) error
	UpdateCredentials(index int, credential *content.CredentialInfo) error
	UpdateFiles(index int, file *content.BinaryFileInfo) error
	UpdateTexts(index int, text *content.TextInfo) error

	GetCreditCardsData() *CreditCards
	GetCredentialsData() *Credentials
	GetFilesData() *Files
	GetTextsData() *Texts

	Clear()
}
