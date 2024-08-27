package transport

import (
	"context"
	"ya-GophKeeper/internal/client/storage"
	"ya-GophKeeper/internal/client/transport/http"
)

type Transport interface {
	Registration(ctx context.Context, userAutData http.UserInfo) error
	Login(ctx context.Context, userAuthData http.UserInfo, loginType string) error
	ChangePassword(ctx context.Context, newPasswd string) error
	Sync(ctx context.Context, items storage.Collection) error
	GetOTP(ctx context.Context) (int, error)
}
