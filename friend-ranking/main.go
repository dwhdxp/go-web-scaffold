package main

import (
	"friend-ranking/router"
)

func main() {
	r := router.Router()

	r.Run(":9090")
}
