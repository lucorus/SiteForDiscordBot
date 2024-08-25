package controllers

import (
	"fmt"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"

	"SiteForDsBot/conf"
	models "SiteForDsBot/models"
)

type DiscordUserController struct {
	beego.Controller
}

// выводит топ 100 пользователей
func (this *DiscordUserController) ListUsers() {
	res, err := models.AllDsBotUsers("100")
	if err != nil {
		this.Data["json"] = err
		this.ServeJSON()
		return
	}
	// this.TplName = "main_page.html"
	this.Data["json"] = res
	this.ServeJSON()
}


// выводит все учётные записи, на которых зарегестрирован текущий пользователь
func (this *DiscordUserController) ListUsersForRequstUser() {
		authHeader := this.Ctx.Request.Header["Authorization"]
	if len(authHeader) == 0 {
		this.Redirect("/user/", 302)
		return
	}

	AuthorizationToken := authHeader[0]

	uuid, err := GetUserUuidFromJWT(AuthorizationToken)
	if err != nil {
		this.Redirect("/user/", 302)
		return
	}

	user, err := models.Find(uuid)
	if err != nil {
		this.Redirect("/user/", 302)
		return
	}

	if user.Discord_server_id == "" {
		this.Data["json"] = nil
	} else {
		servers_accounts, err := models.FindDsBotUsers(user.Discord_server_id)
		if err != nil {
			this.Data["json"] = err
			this.ServeJSON()
		}
		this.Data["json"] = servers_accounts
		this.ServeJSON()
	}
}


// получает токен доступа, токен юзера, id юзера в дискорд и добавляет асоциацию с аккаунтом на сайте
func (this *DiscordUserController) AuthorizeUser() {
	Access := this.Ctx.Request.Header["Access"]
	if len(Access) == 0 {
		fmt.Println("нет токена доступа")
		return
	}
	access_token := Access[0]
	if access_token != conf.AcessToken {
		fmt.Println("неверный токен доступа")
		return
	}

	User := this.Ctx.Request.Header["User"]
	if len(User) == 0 {
		fmt.Println("нет id пользователя")
		return
	}
	user_id, err := strconv.Atoi(User[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	Token := this.Ctx.Request.Header["Token"]
	if len(User) == 0 {
		fmt.Println("нет токена пользователя")
		return
	}
	user_token := Token[0]

  ok := models.Authorize(user_token, user_id)
	if !ok {
		this.Ctx.Output.SetStatus(400)
		return
	}

	this.Ctx.Output.SetStatus(200)
	this.Data["json"] = "success!"
	this.ServeJSON()
}


// получает токен доступа, токен юзера, id юзера в дискорд и убирает асоциацию с аккаунтом на сайте
func (this *DiscordUserController) AnAuthorizeUser() {
	Access := this.Ctx.Request.Header["Access"]
	if len(Access) == 0 {
		fmt.Println("нет токена доступа")
		return
	}
	access_token := Access[0]
	if access_token != conf.AcessToken {
		fmt.Println("неверный токен доступа")
		return
	}

	User := this.Ctx.Request.Header["User"]
	if len(User) == 0 {
		fmt.Println("нет id пользователя")
		return
	}
	user_id, err := strconv.Atoi(User[0])
	if err != nil {
		fmt.Println(err)
		return
	}

  ok := models.AnAuthorize(user_id)
	if !ok {
		this.Ctx.Output.SetStatus(400)
		return
	}
	
	this.Ctx.Output.SetStatus(200)
	this.Data["json"] = "success!"
	this.ServeJSON()
}
