package handlers

import (
	"net/http"
	"ya-GophKeeper/internal/server/otp"
	"ya-GophKeeper/internal/server/storage"
)

func RegistrationPOST(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
}

func LoginWithPasswordPOST(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
}

func LoginWithOTP_POST(w http.ResponseWriter, r *http.Request, m *otp.ManagerOTP) {
}

// Клиент отправляет список новых элементов
// и список элементов которые у него есть с датами модификации
// Возвращает 3 списка:
// 1.
func SyncDataFirstStepPOST(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
}

func SyncDataSecondStepPOST(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
}

func GenerateOTP_GET(w http.ResponseWriter, r *http.Request, m *otp.ManagerOTP) {
}
