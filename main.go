package main

import (
	"lrumemcache/api"
	"lrumemcache/utils"
)

func main() {
	config := utils.NewConfig()
	service := api.NewService(config.ServerURL, config.Capacity)
	service.StartService()
}
