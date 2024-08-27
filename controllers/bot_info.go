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
		this.Ctx.Output.SetStatus(400)
		return
	}
	// this.TplName = "main_page.html"
	this.Data["json"] = res
	this.Ctx.Output.SetStatus(200)
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
		this.Ctx.Output.SetStatus(400)
		return
	} else {
		servers_accounts, err := models.FindDsBotUsers(user.Discord_server_id)
		if err != nil {
			this.Ctx.Output.SetStatus(400)
			return
		}
		this.Data["json"] = servers_accounts
		this.Ctx.Output.SetStatus(200)
		this.ServeJSON()
	}
}


// возвращает всех пользователей на сервере с переданным id
func (this *DiscordUserController) ListUsersInGuild() {
	GuildId := this.Ctx.Input.Param(":id")
	users, err := models.ListUsersInGuild(GuildId)
	if err != nil {
		this.Ctx.Output.SetStatus(404)
		return
	}
	this.Data["json"] = users
	this.Ctx.Output.SetStatus(200)
	this.ServeJSON()
}


// получает токен доступа, токен юзера, id юзера в дискорд и добавляет асоциацию с аккаунтом на сайте
func (this *DiscordUserController) AuthorizeUser() {
	Access := this.Ctx.Request.Header["Access"]
	if len(Access) == 0 {
		this.Ctx.Output.SetStatus(401)
		return
	}
	access_token := Access[0]
	if access_token != conf.AcessToken {
		this.Ctx.Output.SetStatus(401)
		return
	}

	User := this.Ctx.Request.Header["User"]
	if len(User) == 0 {
		this.Ctx.Output.SetStatus(400)
		return
	}
	user_id, err := strconv.Atoi(User[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	Token := this.Ctx.Request.Header["Token"]
	if len(User) == 0 {
		this.Ctx.Output.SetStatus(400)
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
		this.Ctx.Output.SetStatus(401)
		return
	}
	access_token := Access[0]
	if access_token != conf.AcessToken {
		this.Ctx.Output.SetStatus(401)
		return
	}

	User := this.Ctx.Request.Header["User"]
	if len(User) == 0 {
		this.Ctx.Output.SetStatus(400)
		return
	}
	user_id, err := strconv.Atoi(User[0])
	if err != nil {
		this.Ctx.Output.SetStatus(400)
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
