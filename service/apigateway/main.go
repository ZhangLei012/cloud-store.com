package main

import "cloud-store.com/service/apigateway/route"

func main() {
	r := route.Router()
	r.Run(":8080")
}
