package web

import (
	"cronsun/db/entries"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
	"time"

	"cronsun/log"
	"github.com/gorilla/mux"
)

type Administrator struct{}

func adminAuthHandler(ctx *Context) (abort bool) {
	return authHandler(true, entries.Administrator)(ctx)
}

func NewAdminAuthHandler(f func(ctx *Context)) BaseHandler {
	return BaseHandler{
		BeforeHandle: adminAuthHandler,
		Handle:       f,
	}
}

type Account struct {
	Role       entries.Role       `json:"role"`
	Email      string             `json:"email"`
	Status     entries.UserStatus `json:"status"`
	Session    bool               `json:"session"`
	CreateTime time.Time          `json:"createTime"`
}

func (this *Administrator) GetAccountList(ctx *Context) {

	list, err := entries.GetAccounts(nil)
	if err != nil {
		outJSONWithCode(ctx.W, http.StatusInternalServerError, err.Error())
		return
	}

	var alist = make([]Account, len(list))
	for i := range list {
		alist[i] = Account{
			Role:       list[i].Role,
			Email:      list[i].Email,
			Status:     list[i].Status,
			Session:    list[i].Session != "",
			CreateTime: list[i].CreateTime,
		}
	}

	outJSONWithCode(ctx.W, http.StatusOK, alist)
}

func (this *Administrator) GetAccount(ctx *Context) {
	vars := mux.Vars(ctx.R)
	email := strings.TrimSpace(vars["email"])
	if email == "" {
		outJSONWithCode(ctx.W, http.StatusBadRequest, "Email required.")
		return
	}

	u, err := entries.GetAccountByEmail(email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			outJSONWithCode(ctx.W, http.StatusNotFound, fmt.Sprintf("Email [%s] not found.", email))
		} else {
			outJSONWithCode(ctx.W, http.StatusInternalServerError, err.Error())
		}
		return
	}

	outJSONWithCode(ctx.W, http.StatusOK, &Account{
		Role:       u.Role,
		Email:      u.Email,
		Status:     u.Status,
		Session:    u.Session != "",
		CreateTime: u.CreateTime,
	})
}

func (this *Administrator) AddAccount(ctx *Context) {
	account := struct {
		Role     entries.Role `json:"role"`
		Email    string       `json:"email"`
		Password string       `json:"password"`
	}{}

	decoder := json.NewDecoder(ctx.R.Body)
	err := decoder.Decode(&account)
	if err != nil {
		outJSONWithCode(ctx.W, http.StatusBadRequest, err.Error())
		return
	}
	ctx.R.Body.Close()

	if !account.Role.Defined() {
		outJSONWithCode(ctx.W, http.StatusBadRequest, "Account role undefined.")
		return
	}

	account.Email = strings.TrimSpace(account.Email)
	if len(account.Email) == 0 {
		outJSONWithCode(ctx.W, http.StatusBadRequest, "Account email is required.")
		return
	}

	account.Password = strings.TrimSpace(account.Password)
	if len(account.Password) == 0 {
		outJSONWithCode(ctx.W, http.StatusBadRequest, "Account password is required.")
		return
	}

	u, err := entries.GetAccountByEmail(account.Email)
	if err == nil || u != nil {
		outJSONWithCode(ctx.W, http.StatusConflict, fmt.Sprintf("Email [%s] has been used.", account.Email))
		return
	}

	salt := genSalt()
	err = entries.CreateAccount(&entries.Account{
		Role:     account.Role,
		Email:    account.Email,
		Salt:     salt,
		Status:   entries.UserActived,
		Password: encryptPassword(account.Password, salt),
	})
	if err != nil {
		outJSONWithCode(ctx.W, http.StatusBadRequest, fmt.Sprintf("Failed to create user: %s.", err.Error()))
		return
	}

	outJSONWithCode(ctx.W, http.StatusNoContent, nil)
}

func (this *Administrator) UpdateAccount(ctx *Context) {
	account := struct {
		Role        entries.Role       `json:"role"`
		OriginEmail string             `json:"originEmail"`
		Email       string             `json:"email"`
		Password    string             `json:"password"`
		Status      entries.UserStatus `json:"status"`
	}{}

	decoder := json.NewDecoder(ctx.R.Body)
	err := decoder.Decode(&account)
	if err != nil {
		outJSONWithCode(ctx.W, http.StatusBadRequest, err.Error())
		return
	}
	ctx.R.Body.Close()

	account.OriginEmail = strings.TrimSpace(account.OriginEmail)
	if account.OriginEmail == "" {
		outJSONWithCode(ctx.W, http.StatusBadRequest, "Account origin email is required.")
		return
	}

	if !account.Role.Defined() {
		outJSONWithCode(ctx.W, http.StatusBadRequest, "Account role undefined.")
		return
	}

	if !account.Status.Defined() {
		outJSONWithCode(ctx.W, http.StatusBadRequest, "Account status undefined.")
		return
	}

	account.Email = strings.TrimSpace(account.Email)
	if len(account.Email) == 0 {
		outJSONWithCode(ctx.W, http.StatusBadRequest, "Account email is required.")
		return
	}

	originAccount, err := entries.GetAccountByEmail(account.OriginEmail)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			outJSONWithCode(ctx.W, http.StatusNotFound, "Email not found.")
		} else {
			outJSONWithCode(ctx.W, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if originAccount.Unchangeable && originAccount.Email != ctx.Session.Email {
		outJSONWithCode(ctx.W, http.StatusForbidden, "You can not change this account.")
		return
	}

	var update = bson.M{}
	if !originAccount.Unchangeable {
		update = bson.M{
			"status": account.Status,
			"role":   account.Role,
		}
	}

	if account.Email != account.OriginEmail {
		update["email"] = account.Email
	}

	account.Password = strings.TrimSpace(account.Password)
	if len(account.Password) != 0 {
		salt := genSalt()
		update["salt"] = salt
		update["password"] = encryptPassword(account.Password, salt)
	}

	if len(update) == 0 {
		outJSONWithCode(ctx.W, http.StatusOK, nil)
		return
	}

	err = entries.UpdateAccount(bson.M{"email": account.OriginEmail}, update)
	if err != nil {
		outJSONWithCode(ctx.W, http.StatusBadRequest, fmt.Sprintf("Failed to update user: %s.", err.Error()))
		return
	}

	this.removeSession(account.Email)
	if ctx.Session.Email == originAccount.Email {
		ctx.Session.Email = ""
		delete(ctx.Session.Data, "role")
		ctx.Session.Store()

		outJSONWithCode(ctx.W, http.StatusUnauthorized, nil)
		return
	}

	outJSONWithCode(ctx.W, http.StatusOK, nil)
}

func (this *Administrator) removeSession(email string) {
	u, err := entries.GetAccountByEmail(email)
	if err != nil {
		log.Errorf("Failed to remove user [%s] session: %s", email, err.Error())
		return
	}

	if u.Session == "" {
		return
	}

	sessManager.CleanSeesionData(u.Session)
}
