package srverror

import "errors"

var ErrLoginAlreadyTaken = errors.New("login already taken")
var ErrLoginNotFound = errors.New("login not found")
var ErrIncorrectDataTpe = errors.New("incorrect data type")
var ErrIncorrectSyncStep = errors.New("incorrect sync step (must be '1' or '2')")
