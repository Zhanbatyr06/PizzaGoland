package main

import (
	app2 "github.com/Zhanbatyr06/PizzaGoland/app"
	"log"
)

func main() {
	app, err := app2.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	app.Start()
}
