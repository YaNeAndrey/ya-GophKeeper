package storage

type AuthDateInfo struct {
	ID       int
	Resource string
	Login    string
	Password string
}

type FileInfo struct {
	ID           int
	BaseFileName string
	Content      []byte
	Description  string
}

type CreditCardInfo struct {
	ID         int
	CardNumber string
	CVV        string
	ValidThru  string
	Bank       string
}
