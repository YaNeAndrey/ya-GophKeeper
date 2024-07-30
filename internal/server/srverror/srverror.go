package srverror

import "errors"

var ErrLoginAlreadyTaken = errors.New("login already taken")
var ErrLoginNotFound = errors.New("login not found")
