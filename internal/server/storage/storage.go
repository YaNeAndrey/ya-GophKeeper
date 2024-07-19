package storage

import "context"

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
	RemoveAuthDate(ctx context.Context, login string, authDateIDs []int) error

	InsertOrUpdateFiles(ctx context.Context, login string, file []FileInfo) ([]int, error)
	InsertOrUpdateCreditCards(ctx context.Context, login string, creditCard []CreditCardInfo) error
	InsertOrUpdateAuthDate(ctx context.Context, login string, authDate []AuthDateInfo) error

	GetAllFilesInfo(ctx context.Context, login string) ([]FileInfo, error)
	GetAllCreditCardIDs(ctx context.Context, login string) ([]int, error)
	GetAllAuthDateIDs(ctx context.Context, login string) ([]int, error)

	GetCreditCards(ctx context.Context, login string, cardIDs []int) ([]CreditCardInfo, error)
	GetAuthDate(ctx context.Context, login string, authDateIDs []int) ([]AuthDateInfo, error)
}
