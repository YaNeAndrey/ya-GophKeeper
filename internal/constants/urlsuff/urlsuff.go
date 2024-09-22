package urlsuff

const (
	SyncFirstStep  = "1"
	SyncSecondStep = "2"

	OperationRemove        = "remove"
	OperationInsertNew     = "add"
	OperationSync          = "sync"
	OperationRegistration  = "registration"
	OperationLogin         = "login"
	OperationGenerateOTP   = "otp"
	OperationChangPassword = "changepass"

	FileOperationUpload   = "upload"
	FileOperationDownload = "download"

	DatatypeCredential = "cred"
	DatatypeCreditCard = "card"
	DatatypeText       = "text"
	DatatypeFile       = "file"

	LoginTypeOTP    = "otp"
	LoginTypePasswd = "passwd"
)
