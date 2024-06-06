package web

import (
	"cronsun/db/entries"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"

	"strings"

	"cronsun/conf"
	"cronsun/log"
	"cronsun/utils"
)

func checkAuthBasicData() error {
	if conf.Config.Web.Auth.Enabled {
		log.Infof("Authentication enabled.")

		list, err := entries.GetAccounts(bson.M{"role": entries.Administrator, "status": entries.UserActived})
		if err != nil {
			return fmt.Errorf("Failed to check available Administrators: %s.", err.Error())
		}

		if len(list) == 0 {
			// create a default administrator with admin@admin.com/admin
			// the email and password can be change from the user profile page.
			salt := genSalt()
			err = entries.CreateAccount(&entries.Account{
				Role:         entries.Administrator,
				Email:        "admin@admin.com",
				Salt:         salt,
				Password:     encryptPassword("admin", salt),
				Status:       entries.UserActived,
				Unchangeable: true,
			})
			if err != nil {
				return fmt.Errorf("Failed to create default Administrators: %s.", err.Error())
			}
		}

		if err = entries.EnsureAccountIndex(); err != nil {
			log.Warnf("Failed to make db index on the `user` collection: %s.", err.Error())
		}
	}

	return nil
}

func encryptPassword(pwd, salt string) string {
	m := md5.Sum([]byte(pwd + salt))
	m = md5.Sum(m[:])
	return hex.EncodeToString(m[:])
}

func genSalt() string {
	return utils.RandString(8)
}

type Authentication struct{}

func (this *Authentication) GetAuthSession(ctx *Context) {
	var authInfo = &struct {
		Role        entries.Role `json:"role,omitempty"`
		Email       string       `json:"email,omitempty"`
		EnabledAuth bool         `json:"enabledAuth"`
	}{}

	if !conf.Config.Web.Auth.Enabled {
		outJSONWithCode(ctx.W, http.StatusOK, authInfo)
		return
	}

	authInfo.EnabledAuth = true

	if ctx.Session.Email != "" {
		authInfo.Email = ctx.Session.Email
		authInfo.Role = ctx.Session.Data["role"].(entries.Role)
		outJSONWithCode(ctx.W, http.StatusOK, authInfo)
		return
	}

	if len(getStringVal("check", ctx.R)) > 0 {
		outJSONWithCode(ctx.W, http.StatusUnauthorized, nil)
		return
	}

	email := getStringVal("email", ctx.R)
	password := getStringVal("password", ctx.R)
	remember := getStringVal("remember", ctx.R) == "on"

	u, err := entries.GetAccountByEmail(email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			outJSONWithCode(ctx.W, http.StatusNotFound, "User ["+email+"] not found.")
		} else {
			outJSONWithCode(ctx.W, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if u.Password != encryptPassword(password, u.Salt) {
		outJSONWithCode(ctx.W, http.StatusBadRequest, "Incorrect password.")
		return
	}

	if u.Status != entries.UserActived {
		outJSONWithCode(ctx.W, http.StatusForbidden, "Access deny.")
		return
	}

	if !remember {
		if c, err := ctx.R.Cookie(conf.Config.Web.Session.CookieName); err == nil {
			c.MaxAge = 0
			c.Path = "/"
			http.SetCookie(ctx.W, c)
		}
	}

	ctx.Session.Email = u.Email
	ctx.Session.Data["role"] = u.Role
	ctx.Session.Store()

	authInfo.Role = u.Role
	authInfo.Email = u.Email

	err = entries.UpdateAccount(bson.M{"email": email}, bson.M{"session": ctx.Session.ID()})
	outJSONWithCode(ctx.W, http.StatusOK, authInfo)
}

func (this *Authentication) DeleteAuthSession(ctx *Context) {
	ctx.Session.Email = ""
	delete(ctx.Session.Data, "role")
	ctx.Session.Store()

	outJSONWithCode(ctx.W, http.StatusOK, nil)
}

func (this *Authentication) SetPassword(ctx *Context) {
	var sp = &struct {
		Password    string `json:"password"`
		NewPassword string `json:"newPassword"`
	}{}

	decoder := json.NewDecoder(ctx.R.Body)
	err := decoder.Decode(&sp)
	if err != nil {
		outJSONWithCode(ctx.W, http.StatusBadRequest, err.Error())
		return
	}
	ctx.R.Body.Close()

	sp.Password = strings.TrimSpace(sp.Password)
	sp.NewPassword = strings.TrimSpace(sp.NewPassword)
	if sp.Password == "" {
		outJSONWithCode(ctx.W, http.StatusBadRequest, "Passowrd is required.")
		return
	}
	if sp.NewPassword == "" {
		outJSONWithCode(ctx.W, http.StatusBadRequest, "New passowrd is required.")
		return
	}

	var email = ctx.Session.Email
	u, err := entries.GetAccountByEmail(email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			outJSONWithCode(ctx.W, http.StatusNotFound, "User ["+email+"] not found.")
		} else {
			outJSONWithCode(ctx.W, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if u.Password != encryptPassword(sp.Password, u.Salt) {
		outJSONWithCode(ctx.W, http.StatusBadRequest, "Incorrect password.")
		return
	}

	salt := genSalt()
	update := bson.M{
		"salt":     salt,
		"password": encryptPassword(sp.NewPassword, salt),
	}

	if err = entries.UpdateAccount(bson.M{"email": email}, update); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			outJSONWithCode(ctx.W, http.StatusBadRequest, "User ["+email+"] not found.")
		} else {
			outJSONWithCode(ctx.W, http.StatusInternalServerError, err.Error())
		}
		return
	}

	outJSONWithCode(ctx.W, http.StatusOK, nil)
}
