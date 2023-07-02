package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"wannabe/response"
)

type AuthController struct {
	Db *sql.DB
}

type RegisterController struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	ImgUrl   string `json:"img_url" validate:"required"`
}

type LoginController struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Auth struct {
	Id       int
	Email    string
	Password string
}

var (
	queryCreate = `
		INSERT INTO auth (email, password, img_url)
		VALUES ($1, $2, $3)
	`
	queryFindByEmail = `
		SELECT id, email, password
		FROM auth
		WHERE email=$1
	`
)

func (a *AuthController) Register(c *gin.Context) {
	var req = RegisterController{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(500, err.Error())
		return
	}

	//validate required info
	val := validator.New()
	if err := val.Struct(req); err != nil {
		c.JSON(500, gin.H{
			"Messages": "Missed requirement info",
		})
		return
	}

	//Implement encrypt password using bcrpyt
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, err.Error())
	}
	req.Password = string(hash)

	//implement query db
	stmt, err := a.Db.Prepare(queryCreate)
	if err != nil {
		c.JSON(401, err.Error())
		return
	}
	_, err = stmt.Exec(req.Email, req.Password, req.ImgUrl)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	//200 ok resp
	resp := response.ResponseApi{
		StatusCode: 201,
		Messages:   "Successfully Created",
	}
	c.JSON(resp.StatusCode, resp)
}

func (a *AuthController) Login(c *gin.Context) {

	var req = LoginController{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(500, err.Error())
	}

	val := validator.New()
	if err := val.Struct(req); err != nil {
		c.AbortWithStatusJSON(500, err.Error())
		return
	}

	stmt, err := a.Db.Prepare(queryFindByEmail)
	if err != nil {
		c.AbortWithStatusJSON(500, err.Error())
		return
	}
	row := stmt.QueryRow(req.Email)

	var auth = Auth{}

	if err = row.Scan(&auth.Id, &auth.Email, &auth.Password); err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"Messages": "Wrong Email / Unregistered",
		})
		return
	}

	//compare password
	if err = bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(req.Password)); err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"Messages": "Wrong password",
		})
		return
	}

	c.JSON(200, gin.H{
		"Messages": "Success",
	})
}
