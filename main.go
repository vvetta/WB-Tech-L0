package main

import (
	"WB-Tech-L0/config"
)

func main() {
	err := config.Init()
	if err != nil {
		return
	}
}
