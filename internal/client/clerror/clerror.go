package clerror

import "errors"

var ErrOutOfRange = errors.New("index is out of range")
var ErrCopyFileProblem = errors.New("failed to copy file")
var ErrFileNotFound = errors.New("file not found")
var ErrCreditCardLuhn = errors.New("incorrect credit card. Luhn problem")
var ErrCreditCardIncorrectChar = errors.New("credit card number contain incorrect symbols")
var ErrIncorrectValueCVV = errors.New("CVV must be in range [100:999]")
var ErrAllFieldsMustBeFulled = errors.New("all fields must be filled")
