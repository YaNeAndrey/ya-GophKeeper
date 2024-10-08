package srverror

import "errors"

var ErrLoginAlreadyTaken = errors.New("login already taken")
var ErrLoginNotFound = errors.New("login not found")
var ErrFileNotFound = errors.New("file not found")
var ErrIncorrectDataTpe = errors.New("incorrect data type")
var ErrIncorrectSyncStep = errors.New("incorrect sync step (must be '1' or '2')")
var ErrIncorrectOTP = errors.New("otp must be int")
var ErrIncorrectFileHash = errors.New("incorrect file hash. try again")
var ErrDBModificationDenied = errors.New("the user does not have rights to change the specified data")
