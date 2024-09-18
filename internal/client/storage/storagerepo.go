package storage

import (
	"ya-GophKeeper/internal/client/storage/collection"
	"ya-GophKeeper/internal/content"
)

type StorageRepo interface {
	AddCreditCard(creditCard *content.CreditCardInfo)
	AddCredential(credential *content.CredentialInfo)
	AddFile(file *content.BinaryFileInfo)
	AddText(text *content.TextInfo)

	RemoveCreditCard(index int) error
	RemoveCredential(index int) error
	RemoveText(index int) error
	RemoveFile(index int) error

	UpdateCreditCards(index int, creditCard *content.CreditCardInfo) error
	UpdateCredentials(index int, credential *content.CredentialInfo) error
	UpdateFiles(index int, file *content.BinaryFileInfo) error
	UpdateTexts(index int, text *content.TextInfo) error

	GetCreditCardsData() *collection.CreditCards
	GetCredentialsData() *collection.Credentials
	GetFilesData() *collection.Files
	GetTextsData() *collection.Texts

	Clear()
}
