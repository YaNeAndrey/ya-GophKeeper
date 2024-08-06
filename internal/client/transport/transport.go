package transport

import (
	"context"
	"ya-GophKeeper/internal/client/storage"
)

type Transport interface {
	Registration(ctx context.Context, userAutData UserInfo) error
	Login(ctx context.Context, userAuthData UserInfo, loginType string) error
	ChangePassword(ctx context.Context, newLogin string) error
	Sync(ctx context.Context, items storage.Collection) error
	GetOTP(ctx context.Context) (int, error)
}
