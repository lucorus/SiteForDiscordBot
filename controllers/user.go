package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/dgrijalva/jwt-go"

	"SiteForDsBot/conf"
	models "SiteForDsBot/models"
)

type UserController struct {
	beego.Controller
}


func (this *UserController) ListUsers() {
	res, err := models.All()
	if err != nil {
		this.Data["error"] = err
		this.ServeJSON()
		return
	}
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


func (this *UserController) GetUser() {
	uuid := this.Ctx.Input.Param(":uuid")
	user, err := models.Find(uuid)
	if err != nil {
		this.Ctx.Output.SetStatus(404)
		this.Ctx.Output.Body([]byte("user not found"))
	}
	this.Ctx.Output.SetStatus(200)
	this.Data["json"] = user
	this.ServeJSON()
}


func GenerateJWT(uuid string) (string, error) {
    claims := jwt.MapClaims{
        "uuid": uuid,
        "exp":  time.Now().Add(time.Hour * 1).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(conf.Jwt_secret))
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

	token, err := GenerateJWT(uuid)
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
	body, err := ioutil.ReadAll(this.Ctx.Request.Body)
		if err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("empty fields"))
		return
	}

	req := struct { Uuid string}{}
	if err := json.Unmarshal(body, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("empty fields"))
		fmt.Println(req, err)
		return
	}
  err = models.DeleteUser(req.Uuid)
	if err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	this.Ctx.Output.SetStatus(200)
	this.Data["json"] = "success!"
	this.ServeJSON()
}