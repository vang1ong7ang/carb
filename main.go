package main

import (
	"carb/app"
	"encoding/json"
	"io/ioutil"
)

func main() {
	var app app.T
	conf, _ := ioutil.ReadFile(".carb.json")
	json.Unmarshal(conf, &app)
	app.Start()
	select {}
}
