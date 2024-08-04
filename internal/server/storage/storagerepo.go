package storage

import (
	"context"
	"time"
	"ya-GophKeeper/internal/content"
)

type StorageRepo interface {
	AddNewUser(ctx context.Context, login string, password string) error
	CheckUserPassword(ctx context.Context, login string, password string) (bool, error)
	ChangeUserPassword(ctx context.Context, login string, password string) error

	AddNewFiles(ctx context.Context, login string, files []content.BinaryFileInfo) ([]content.BinaryFileInfo, error)
	AddNewTexts(ctx context.Context, login string, texts []content.TextInfo) ([]content.TextInfo, error)
	AddNewCreditCards(ctx context.Context, login string, creditCards []content.CreditCardInfo) ([]content.CreditCardInfo, error)
	AddNewCredentials(ctx context.Context, login string, credentials []content.CredentialInfo) ([]content.CredentialInfo, error)

	RemoveFiles(ctx context.Context, login string, fileIDs []int) ([]string, error)
	RemoveTexts(ctx context.Context, login string, textIDs []int) error
	RemoveCreditCards(ctx context.Context, login string, creditCardIDs []int) error
	RemoveCredentials(ctx context.Context, login string, credentialIDs []int) error

	UpdateFiles(ctx context.Context, login string, files []content.BinaryFileInfo) error
	UpdateTexts(ctx context.Context, login string, texts []content.TextInfo) error
	UpdateCreditCards(ctx context.Context, login string, creditCards []content.CreditCardInfo) error
	UpdateCredentials(ctx context.Context, login string, credentials []content.CredentialInfo) error

	GetCreditCards(ctx context.Context, login string, cardIDs []int) ([]content.CreditCardInfo, error)
	GetCredentials(ctx context.Context, login string, credIDs []int) ([]content.CredentialInfo, error)
	GetFiles(ctx context.Context, login string, fileIDs []int) ([]content.BinaryFileInfo, error)
	GetTexts(ctx context.Context, login string, textIDs []int) ([]content.TextInfo, error)

	GetModtimeWithIDs(ctx context.Context, login string, dataType string) (map[int]time.Time, error)
}
