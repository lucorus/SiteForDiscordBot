package routers

import (
	"SiteForDsBot/controllers"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/google/uuid"
)

func init() {
		beego.Router("/user/", &controllers.UserController{}, "get:ListUsers;delete:DeleteUser;put:UpdateUser")

		beego.Router("/register/", &controllers.UserController{}, "post:NewUser")
		beego.Router("/login/", &controllers.UserController{}, "post:LoginUser")
		beego.Router("/profile/", &controllers.UserController{}, "get:Profile")
		beego.Router("/profile/:uuid", &controllers.UserController{}, "get:Profile")

		beego.Router("/main_page/", &controllers.DiscordUserController{}, "get:ListUsers")
		beego.Router("/my_accounts/", &controllers.DiscordUserController{}, "get:ListUsersForRequstUser")
		beego.Router("/guild/:id:int", &controllers.DiscordUserController{}, "get:ListUsersInGuild")

		beego.Router("/change_token/", &controllers.UserController{}, "patch:ChangeToken")

		beego.Router("/authorize/", &controllers.DiscordUserController{}, "post:AuthorizeUser")
		beego.Router("/anauthorize/", &controllers.DiscordUserController{}, "post:AnAuthorizeUser")
}
