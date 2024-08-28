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
	authHeader := this.Ctx.Request.Header["Authorization"]
	if len(authHeader) == 0 {
		this.Redirect("/user/", 302)
		return
	}
	AuthorizationToken := authHeader[0]

	Uuid, err := GetUserUuidFromJWT(AuthorizationToken)
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

	Uuid, err := GetUserUuidFromJWT(AuthorizationToken)
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


func GenerateJWT(uuid string) (string, error) {
    claims := jwt.MapClaims{
        "uuid": uuid,
        "exp":  time.Now().Add(time.Hour * 100).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(conf.Jwt_secret))
}


func GetUserUuidFromJWT(tokenString string) (string, error) {
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
  if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
    return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
  }
    return []byte(conf.Jwt_secret), nil
  })
  if err != nil {
    return "", err
  }

  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    uuid, ok := claims["uuid"].(string)
    if !ok {
      return "", fmt.Errorf("uuid not found or invalid type")
    }
    return uuid, nil
  }
  return "", fmt.Errorf("invalid token")
}


// выводит данные пользователя с переданным uuid или берёт uuid из jwt токена
func (this *UserController) Profile() {
	uuid := this.Ctx.Input.Param(":uuid")
	if len(uuid) == 0 {
		authHeader := this.Ctx.Request.Header["Authorization"]
		if len(authHeader) == 0 {
			this.Redirect("/user/", 302)
			return
		}
		AuthorizationToken := authHeader[0]

		Uuid, err := GetUserUuidFromJWT(AuthorizationToken)
		if err != nil {
			this.Redirect("/user/", 302)
			return
		}
		uuid = Uuid
	}

	user, err := models.Find(uuid)
	if err != nil {
		this.Redirect("/user/", 302)
		return
	}

	this.Data["json"] = user
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

	uuid, err := GetUserUuidFromJWT(AuthorizationToken)
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
