package urlsuff

const (
	OperationRemove        = "remove"
	OperationInsertNew     = "add"
	OperationSync          = "sync"
	OperationRegistration  = "registration"
	OperationLogin         = "login"
	OperationGenerateOTP   = "otp"
	OperationChangPassword = "changepass"

	DatatypeCredential = "cred"
	DatatypeCreditCard = "card"
	DatatypeText       = "text"
	DatatypeFile       = "file"

	LoginTypeOTP    = "otp"
	LoginTypePasswd = "passwd"
)
