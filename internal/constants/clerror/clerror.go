package clerror

import "errors"

var ErrOutOfRange = errors.New("index is out of range")
var ErrCopyFileProblem = errors.New("failed to copy file")
var ErrFileNotFound = errors.New("file not found")
var ErrCreditCardLuhn = errors.New("incorrect credit card. Luhn problem")
var ErrCreditCardIncorrectChar = errors.New("credit card number contain incorrect symbols")
var ErrIncorrectValueCVV = errors.New("CVV must be in range [100:999]")
var ErrAllRequiredFieldsMustBeFulled = errors.New("all required fields must be filled")
var ErrIncorrectType = errors.New("use another type")
var ErrAuthTokenIsEmpty = errors.New("server did not return an authorization token")
var ErrLoginAlreadyTaken = errors.New("login already taken")
var ErrIncorrectPassword = errors.New("incorrect password or OTP")
var ErrMaxTextSizeExceeded = errors.New("text size limit exceeded")
var ErrPathIsDir = errors.New("path points to a directory")
var ErrMaxFileSizeExceeded = errors.New("file size limit exceeded")

var ErrInternalServerError = errors.New("internal server error")
var ErrNotAuthorized = errors.New("not authorized request. try re-login")
var ErrBadRequest = errors.New("bad request")
var ErrStatusNotFound = errors.New("incorrect request url")
