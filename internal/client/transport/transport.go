package transport

import (
	"context"
	"ya-GophKeeper/internal/client/transport/http"
)

type Transport interface {
	Registration(ctx context.Context, userAutData http.UserInfo) error
	Login(ctx context.Context, userAuthData http.UserInfo, loginType string) error
	ChangePassword(ctx context.Context, newPasswd string) error
	SyncRemovedItems(ctx context.Context, removedIDs []int, dataType string) error
	SyncNewItems(ctx context.Context, bodyWithNewItems []byte, dataType string) ([]byte, error)
	SyncChangesFirstStep(ctx context.Context, bodyIDsWithModtime []byte, dataType string) ([]byte, error)
	SyncChangesSecondStep(ctx context.Context, bodyDataForSrv []byte, dataType string) error
	GetOTP(ctx context.Context) (int, error)
}
