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
		beego.Router("/profile/", &controllers.UserController{}, "get:Profile")
		beego.Router("/main_page/", &controllers.DiscordUserController{}, "get:ListUsers")
		beego.Router("/authorize/", &controllers.DiscordUserController{}, "post:AuthorizeUser")
		beego.Router("/anauthorize/", &controllers.DiscordUserController{}, "post:AnAuthorizeUser")
}
