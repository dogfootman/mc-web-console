package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/golang-jwt/jwt/v4"

	mcmodels "mc_web_console_common_models"
)

func McIamAuthMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		accessToken := c.Session().Get("Authorization")
		if accessToken == nil {
			c.Flash().Add("danger", "No session")
			c.Session().Clear()
			return c.Redirect(http.StatusSeeOther, "authLoginPath()")
		}
		jwtDecode, _ := jwtDecode(accessToken.(string))
		t := time.Unix(int64(jwtDecode["exp"].(float64)), 0)
		if t.Before(time.Now()) {
			fmt.Println(time.Since(t))
			c.Session().Clear()
			c.Flash().Add("danger", "Session Expired")
			return c.Redirect(http.StatusSeeOther, "authLoginPath()")
		}

		c.Set("name", jwtDecode["name"])

		return next(c)
	}
}

// userinfo, msg, err := getUserInfo(c)
//
//	if err != nil {
//		log.Println(err.Error())
//		c.Session().Clear()
//		c.Flash().Add("danger", msg)
//		return c.Redirect(http.StatusSeeOther, "authLoginPath()")
//	}
func getUserInfo(c buffalo.Context) (mcmodels.UserInfo, string, error) {
	userinfo := &mcmodels.UserInfo{}

	status, userinfoByte, err := CommonAPIGet(APIUserInfoPath, c)
	if err != nil {
		msg := status.Status + " Error Get Userinfo from MC-IAM-MANAGER"
		return *userinfo, msg, err
	}
	if status.StatusCode != 200 {
		msg := status.Status + " Error Get Userinfo from MC-IAM-MANAGER"
		return *userinfo, msg, errors.New(msg)
	}

	if err := json.Unmarshal([]byte(userinfoByte), &userinfo); err != nil {
		msg := "Authentication Info Error"
		return *userinfo, msg, err
	}

	return *userinfo, "", nil
}

func jwtDecode(jwtToken string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) { return "", nil })
	if err != nil {
		return claims, err
	}
	return claims, nil
}
