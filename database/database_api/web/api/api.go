package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/creatorkostas/KeyDB/database/database_core/conf"
	internal "github.com/creatorkostas/KeyDB/database/database_core/conf"
	database "github.com/creatorkostas/KeyDB/database/database_core/core"
	"github.com/creatorkostas/KeyDB/database/database_core/users"
	db_utils "github.com/creatorkostas/KeyDB/database/database_core/utils"
	"github.com/gin-gonic/gin"
	stats "github.com/semihalev/gin-stats"
)

func GetValue(c *gin.Context) {

	var key, key_found = c.GetQuery("key")

	var acc = c.MustGet("Account").(*users.Account)

	var result any
	if key_found {
		result = database.Get_value(acc.Username, key)
	} else {
		result = database.Get_value(acc.Username, "table.get.all.data")
	}

	var res = &Responce{C: c, ErrorMessage: "Key does not exist", Result: result, OkCode: http.StatusOK, ErrorCode: http.StatusBadRequest}
	res.sendResponce()

}

// func SetValues(key string, value_type string, value interface{}) {
func SetValues(c *gin.Context) {

	var acc = c.MustGet("Account").(*users.Account)
	var key, key_found = c.GetQuery("key")
	var value_type, value_type_found = c.GetQuery("type")
	var data, data_found = c.GetQuery("value")

	var error_message string = "Something went wrong!"
	var error_code int = http.StatusInternalServerError
	var err error = nil

	if key_found && value_type_found && data_found {
		err = database.Add_value(acc.Username, key, value_type, data)
	} else {
		error_message = fmt.Sprintf(
			"Missings parameters!! key found: %s , value_type found: %s , data found: %s",
			strconv.FormatBool(key_found),
			strconv.FormatBool(value_type_found),
			strconv.FormatBool(data_found))
		error_code = http.StatusBadRequest
	}

	var res = &Responce{C: c, ErrorMessage: error_message, Result: true, OkCode: http.StatusOK, ErrorCode: error_code, Result_error: err}
	res.sendResponce()

}

func Register(c *gin.Context) {

	var username, username_found = c.GetQuery("username")
	var email, email_found = c.GetQuery("email")
	var password, password_found = c.GetQuery("password")
	var acc_type, acc_type_found = c.GetQuery("type")

	var error_message string = "Something went wrong!"
	var error_code int = http.StatusInternalServerError
	var acc *users.Account = nil
	var err error = errors.New("something went wrong")

	if !conf.Web_Enable_admin_register && acc_type == "Admin" {
		error_message = "Admin register is disabled. Please contact the administrator!!"
		error_code = http.StatusUnauthorized
	} else if username_found && email_found && password_found && acc_type_found {
		acc = users.Create_account(username, acc_type, email, password)
		// log.Println("yesss")
		if acc != nil {
			err = nil
		}
	} else {
		error_message = fmt.Sprintf(
			"Missings parameters!! username found: %s , email found: %s , password found: %s , acc_type found: %s",
			strconv.FormatBool(username_found),
			strconv.FormatBool(email_found),
			strconv.FormatBool(password_found),
			strconv.FormatBool(acc_type_found))
		error_code = http.StatusBadRequest
	}

	var res = &Responce{C: c, ErrorMessage: error_message, Result: acc, OkCode: http.StatusOK, ErrorCode: error_code, Result_error: err}
	res.sendResponce()
}

func ChangeApiKey(c *gin.Context) {
	var acc = c.MustGet("Account").(*users.Account)
	acc.ChangeApiKey()
	var res = &Responce{C: c, ErrorMessage: "Something went wrong!", Result: acc, OkCode: http.StatusOK, ErrorCode: http.StatusBadRequest}
	res.sendResponce()
}

func ChangePassword(c *gin.Context) {
	var acc = c.MustGet("Account").(*users.Account)
	var new_pass, found = c.GetQuery("password")
	var error_message string = "Something went wrong!"
	var error_code int = http.StatusInternalServerError
	if found {
		acc.ChangePassword(new_pass)
	} else {
		error_message = "password parameter missing"
		error_code = http.StatusBadRequest
		acc = nil
	}
	var res = &Responce{C: c, ErrorMessage: error_message, Result: acc, OkCode: http.StatusOK, ErrorCode: error_code}
	res.sendResponce()
}

func GetStats(c *gin.Context) {
	c.JSON(http.StatusOK, stats.Report())
	c.Request.Context().Done()
}

func Save(c *gin.Context) {
	db_utils.SaveDB(internal.DB_filename)
	users.SaveAccounts(internal.Accounts_filename)

	c.JSON(http.StatusOK, JsonResponce{"ok"})
}

func Load(c *gin.Context) {
	db_utils.LoadDB(internal.DB_filename)
	users.LoadAccounts(internal.Accounts_filename)

	c.JSON(http.StatusOK, JsonResponce{"ok"})
}

func DisableAdmin(c *gin.Context) {
	conf.Web_Enable_admin_register = false
	c.JSON(http.StatusOK, JsonResponce{"Admin register disabled"})
}

func EnableAdmin(c *gin.Context) {
	conf.Web_Enable_admin_register = true
	c.JSON(http.StatusOK, JsonResponce{"Admin register enabled"})
}
