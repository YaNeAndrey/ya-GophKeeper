package client

import (
	"fmt"
	"ya-GophKeeper/internal/client/config"
)

type Client struct {
	config *config.Config
	//storage *storage
}

func (c *Client) Start() {
	fmt.Println("*******************************************")
	fmt.Println("Client settings: ")
	fmt.Println(c.config)
	fmt.Println("Build information: ")
	fmt.Println(c.config)
	RunConsoleFunc(c, RootPage)
}
