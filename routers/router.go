package routers

import (
	"SiteForDsBot/controllers"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/google/uuid"
)

func init() {
		beego.Router("/user/", &controllers.UserController{}, "get:ListUsers;post:NewUser;delete:DeleteUser")
		beego.Router("/user/:uuid", &controllers.UserController{}, "get:GetUser")
		beego.Router("/login/", &controllers.UserController{}, "post:LoginUser")
}
