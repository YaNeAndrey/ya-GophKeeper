package main

import (
	"ya-GophKeeper/internal/client"
	"ya-GophKeeper/internal/client/config"
	"ya-GophKeeper/internal/client/storage/memory"
	"ya-GophKeeper/internal/client/transport"
	"ya-GophKeeper/internal/client/transport/http"
)

func main() {
	//flags will be added later
	cnfg, err := config.ParseConfigFromJSON(".\\client.config")
	if err != nil {
		panic(err)
	}

	st := memory.NewBaseStorage(cnfg.TempDir)
	tr := transport.Transport(http.InitTransport(cnfg.SrvAddr, cnfg.ChunkSize))
	myClient := client.NewClient(cnfg, st, tr)
	myClient.Start()
}
