package client

import (
	"fmt"
	"ya-GophKeeper/internal/client/config"
	"ya-GophKeeper/internal/client/storage"
	"ya-GophKeeper/internal/client/transport"
)

type Client struct {
	config    *config.Config
	storage   storage.StorageRepo
	transport transport.Transport
}

func NewClient(c *config.Config, st storage.StorageRepo, t transport.Transport) *Client {
	return &Client{
		config:    c,
		storage:   st,
		transport: t,
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
