package main

import (
	"ya-GophKeeper/internal/client"
	"ya-GophKeeper/internal/client/storage"
)

func main() {
	myClient := client.NewClient(nil, storage.StorageRepo(storage.NewBaseStorage("temp")))
	myClient.Start()
}
