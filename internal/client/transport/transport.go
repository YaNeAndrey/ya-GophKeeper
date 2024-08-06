package transport

import (
	"context"
	"ya-GophKeeper/internal/client/storage"
)

type Transport interface {
	Registration(ctx context.Context, userAutData UserInfo) error
	Login(ctx context.Context, userAutData UserInfo, loginType string) error
	Sync(ctx context.Context, items storage.Collection) error
}
