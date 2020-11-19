package main

import (
	"github.com/go-ini/ini"
)

func main() {
	cfg := ini.Empty()
	haSec, err := cfg.NewSection("hamon")
	if err != nil {
		println(err)
	}
	haSec.NewKey("status", "back up")
	cfg.SaveTo("conf.ini")

}
