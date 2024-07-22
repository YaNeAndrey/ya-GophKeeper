package storage

import (
	"context"
	"ya-GophKeeper/internal/content"
)

type StorageRepo interface {
	AddNewUser(ctx context.Context, login string, password string) error
	CheckUserPassword(ctx context.Context, login string, password string) (bool, error)
	ChangeUserPassword(ctx context.Context, login string, password string) (bool, error)

	/*
		AddNewFile(ctx context.Context, login string, file FileInfo) error
		AddNewCreditCard(ctx context.Context, login string, creditCard CreditCardInfo) error
		AddNewAuthDate(ctx context.Context, login string, authDate AuthDateInfo) error
	*/

	RemoveFiles(ctx context.Context, login string, fileIDs []int) error
	RemoveCreditCards(ctx context.Context, login string, creditCardIDs []int) error
	RemoveCredentials(ctx context.Context, login string, authDateIDs []int) error

	InsertOrUpdateFiles(ctx context.Context, login string, file []content.BinaryFileInfo) ([]int, error)
	InsertOrUpdateTexts(ctx context.Context, login string, file []content.TextInfo) ([]int, error)
	InsertOrUpdateCreditCards(ctx context.Context, login string, creditCard []content.CreditCardInfo) error
	InsertOrUpdateCredentials(ctx context.Context, login string, authDate []content.CredentialInfo) error

	GetAllFilesData(ctx context.Context, login string) ([]content.BinaryFileInfo, error)
	GetAllTextData(ctx context.Context, login string) ([]content.TextInfo, error)
	GetAllCreditCardIDs(ctx context.Context, login string) ([]int, error)
	GetAllCredentialIDs(ctx context.Context, login string) ([]int, error)

	GetCreditCards(ctx context.Context, login string, cardIDs []int) ([]content.CreditCardInfo, error)
	GetCredential(ctx context.Context, login string, authDateIDs []int) ([]content.CredentialInfo, error)
}
