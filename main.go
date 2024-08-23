package main

import (
	"SiteForDsBot/models"
	_ "SiteForDsBot/routers"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	err := models.CreateDB()
	if err == nil {
		beego.Run()
	} else {
		fmt.Println(err)
	}
}

