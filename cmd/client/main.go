package main

import (
	"ya-GophKeeper/internal/client"
	"ya-GophKeeper/internal/client/storage"
	"ya-GophKeeper/internal/client/transport"
	"ya-GophKeeper/internal/client/transport/http"
)

func main() {
	myClient := client.NewClient(nil, storage.StorageRepo(storage.NewBaseStorage(".\\temp")), transport.Transport(http.InitTransport("http://localhost:8080", 5*1024*1024)))
	myClient.Start()
}
