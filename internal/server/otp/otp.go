package otp

import (
	"crypto/rand"
	"math/big"
	"sync"
)

type ManagerOTP struct {
	usersOTP map[string]int
	mutex    sync.Mutex
}

func (m *ManagerOTP) GenerateOTP(login string) (int, error) {
	r, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		return 0, err
	}
	otp := int(r.Int64()) + 1000
	m.mutex.Lock()
	m.usersOTP[login] = otp
	m.mutex.Unlock()
	return otp, nil
}

func (m *ManagerOTP) CheckUserOTP(login string, otp int) bool {
	m.mutex.Lock()
	isCorrect := m.usersOTP[login] == otp
	defer m.mutex.Unlock()
	if isCorrect {
		delete(m.usersOTP, login)
		return true
	}
	return false
}
