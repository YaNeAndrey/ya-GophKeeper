package main

import (
	"net/http"
	"ya-GophKeeper/internal/server/otp"
	"ya-GophKeeper/internal/server/storage"
	"ya-GophKeeper/internal/server/transport/http/router"
)

func main() {
	st := storage.StorageRepo(storage.InitStorageDB("postgresql://postgres:Qwerty123!@localhost:5432/keeper"))
	managerOTP := otp.InitManagerOTP()
	r := router.InitRouter(st, managerOTP)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
