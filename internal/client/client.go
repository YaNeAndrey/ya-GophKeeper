package client

import (
	"fmt"
	"ya-GophKeeper/internal/client/config"
	"ya-GophKeeper/internal/client/storage"
)

type Client struct {
	config  *config.Config
	storage storage.StorageRepo
}

func NewClient(c *config.Config, st storage.StorageRepo) *Client {
	return &Client{
		config:  c,
		storage: st,
	}
}

func (c *Client) Start() {
	fmt.Println("*******************************************")
	fmt.Println("Client settings: ")
	fmt.Println(c.config)
	fmt.Println("Build information: ")
	fmt.Println(c.config)
	RunConsoleFunc(c, StartPage)
}
