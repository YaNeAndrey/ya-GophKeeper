package main

import (
	"ya-GophKeeper/internal/client"
	"ya-GophKeeper/internal/client/storage"
	"ya-GophKeeper/internal/client/transport"
)

func main() {
	myClient := client.NewClient(nil, storage.StorageRepo(storage.NewBaseStorage("temp")), transport.Transport(transport.InitTransport("http://localhost:8080")))
	myClient.Start()
}
