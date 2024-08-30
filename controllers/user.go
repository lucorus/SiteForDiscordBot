package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	beego "github.com/beego/beego/v2/server/web"

	models "SiteForDsBot/models"
	"SiteForDsBot/responses"
	"SiteForDsBot/utils"
)

type UserController struct {
	beego.Controller
}


func (this *UserController) ListUsers() {
	res, err := models.All()
	if err != nil {
		this.Ctx.Output.SetStatus(400)
		return
	}
	//this.TplName = "main_page.html"
	this.Ctx.Output.SetStatus(200)
	this.Data["json"] = res
	this.ServeJSON()
}


func (this *UserController) NewUser() {
	body, err := ioutil.ReadAll(this.Ctx.Request.Body)
		if err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("empty fields"))
		return
	}

	req := struct { Username, Password string}{}
	if err := json.Unmarshal(body, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("empty fields"))
		fmt.Println(req, err)
		return
	}
  err = models.NewUser(req.Username, req.Password)
	if err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	this.Ctx.Output.SetStatus(200)
	this.Data["json"] = "success!"
	this.ServeJSON()
}


func (this *UserController) LoginUser() {
	body, err := ioutil.ReadAll(this.Ctx.Request.Body)
		if err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("empty fields"))
		return
	}
	req := struct { Username, Password string}{}
	if err := json.Unmarshal(body, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("empty fields"))
		return
	}

  uuid, err := models.LoginUser(req.Username, req.Password)
  if err != nil {
    this.Ctx.Output.SetStatus(500)
    this.Ctx.Output.Body([]byte("Error in models"))
    return
  }

	token, err := utils.GenerateJWT(uuid)
  if err != nil {
    this.Ctx.Output.SetStatus(500)
    this.Ctx.Output.Body([]byte("Error in generation jwt"))
    return
  }

  this.Ctx.Output.Header("Authorization", "Bearer "+token)
  this.Ctx.Output.SetStatus(200)
  this.Ctx.Output.Body([]byte(token))
}


func (this *UserController) DeleteUser() {
	authHeader := this.Ctx.Request.Header["Authorization"]
	if len(authHeader) == 0 {
		this.Redirect("/user/", 302)
		return
	}
	AuthorizationToken := authHeader[0]

	Uuid, err := utils.GetUserUuidFromJWT(AuthorizationToken)
	if err != nil {
		this.Redirect("/user/", 302)
		return
	}

  err = models.DeleteUser(Uuid)
	if err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	this.Ctx.Output.SetStatus(200)
	this.Data["json"] = "success!"
	this.ServeJSON()
}


func (this *UserController) UpdateUser() {
	authHeader := this.Ctx.Request.Header["Authorization"]
	if len(authHeader) == 0 {
		this.Redirect("/user/", 302)
		return
	}
	AuthorizationToken := authHeader[0]

	Uuid, err := utils.GetUserUuidFromJWT(AuthorizationToken)
	if err != nil {
		this.Redirect("/user/", 302)
		return
	}

	body, err := ioutil.ReadAll(this.Ctx.Request.Body)
		if err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("empty fields"))
		return
	}

	req := struct { Username, Password string}{}
	if err := json.Unmarshal(body, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("empty fields"))
		fmt.Println(req, err)
		return
	}

	ok := models.UpdateUser(req.Username, req.Password, Uuid)
	if !ok {
		this.Ctx.Output.SetStatus(400)
		return 
	}
	this.Ctx.Output.SetStatus(200)
	this.Data["json"] = "success"
	this.ServeJSON()
}


// возвращает данные пользователя с переданным uuid или берёт uuid из jwt токена
func UserInfo(this *UserController) (*models.User, error) {
	uuid := this.Ctx.Input.Param(":uuid")
	if len(uuid) == 0 {
		authHeader := this.Ctx.Request.Header["Authorization"]
		if len(authHeader) == 0 {
			this.Redirect("/user/", 302)
			return nil, fmt.Errorf("учётные данные не были предоставлены")
		}
		AuthorizationToken := authHeader[0]

		Uuid, err := utils.GetUserUuidFromJWT(AuthorizationToken)
		if err != nil {
			this.Redirect("/user/", 302)
			return nil, fmt.Errorf("uuid не предоставлен")
		}
		uuid = Uuid
	}

	user, err := models.Find(uuid)
	if err != nil {
		this.Redirect("/user/", 302)
		return nil, fmt.Errorf("Пользователь с даннымй uuid не найден")
	}

	return user, nil
}


func (this *UserController) Profile() {
  user, err := UserInfo(this)
  if err != nil {
      this.Redirect("/user/", 302)
      return
  }

  var userDsAccounts []models.DsBotUser
  userDsAccounts, err = ListAccountsUser(this)
  if err != nil {
      fmt.Println(err)
  }

  response := responses.ProfileResponse{
      User:           user,
      UserDsAccounts: userDsAccounts,
  }

  this.Data["json"] = response
  this.Ctx.Output.SetStatus(200)
  this.ServeJSON()
}


// меняет токен текущему пользователю
func (this *UserController) ChangeToken() {
	authHeader := this.Ctx.Request.Header["Authorization"]
	if len(authHeader) == 0 {
		this.Ctx.Output.SetStatus(400)
		return
	}
	AuthorizationToken := authHeader[0]

	uuid, err := utils.GetUserUuidFromJWT(AuthorizationToken)
	if err != nil {
		this.Ctx.Output.SetStatus(400)
		return
	}

	ok := models.ChangeToken(uuid)
	if !ok {
		this.Data["json"] = "Not success"
		this.Ctx.Output.SetStatus(400)
	} else {
		this.Data["json"] = "Success"
		this.Ctx.Output.SetStatus(202)
	}
	this.ServeJSON()
}
