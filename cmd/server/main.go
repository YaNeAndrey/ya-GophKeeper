package main

import (
	"fmt"
	"net/http"
	"ya-GophKeeper/internal/server/config"
	"ya-GophKeeper/internal/server/otp"
	"ya-GophKeeper/internal/server/storage"
	"ya-GophKeeper/internal/server/transport/http/router"
)

func main() {
	//flags will be added later
	cnfg, err := config.ParseConfigFromJSON(".\\server.config")
	if err != nil {
		panic(err)
	}
	st := storage.InitStorageDB(cnfg.DBconnectionString)
	//st := storage.InitStorageDB("postgresql://postgres:Qwerty123!@localhost:5432/keeper")
	managerOTP := otp.InitManagerOTP()
	fm := storage.InitFileManager(cnfg.FileStorageDir)
	r := router.InitRouter(st, managerOTP, fm)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cnfg.ServerPort), r)
	if err != nil {
		panic(err)
	}
}
