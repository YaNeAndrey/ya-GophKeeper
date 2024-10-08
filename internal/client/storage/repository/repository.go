package repository

import (
	"ya-GophKeeper/internal/client/storage/collection"
	"ya-GophKeeper/internal/content"
)

type CreditCardsProvider interface {
	AddCreditCard(creditCard *content.CreditCardInfo)
	RemoveCreditCard(index int) error
	UpdateCreditCards(index int, creditCard *content.CreditCardInfo) error
	GetCreditCardsData() *collection.CreditCards
}

type CredentialsProvider interface {
	AddCredential(credential *content.CredentialInfo)
	RemoveCredential(index int) error
	UpdateCredentials(index int, credential *content.CredentialInfo) error
	GetCredentialsData() *collection.Credentials
}

type TextsProvider interface {
	AddText(text *content.TextInfo)
	RemoveText(index int) error
	UpdateTexts(index int, text *content.TextInfo) error
	GetTextsData() *collection.Texts
}

type FilesProvider interface {
	AddFile(file *content.BinaryFileInfo)
	RemoveFile(index int) error
	UpdateFiles(index int, file *content.BinaryFileInfo) error
	GetFilesData() *collection.Files
}
