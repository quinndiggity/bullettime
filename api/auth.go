package api

import (
	"fmt"
	"strings"

	"github.com/Rugvip/bullettime/service"

	"net/http"
)

type LoginType string

const (
	LoginTypePassword LoginType = "m.login.password"
	LoginTypeEmail              = "m.login.email.identity"
)

type AuthFlow struct {
	Stages []LoginType `json:"stages,omitempty"`
	Type   LoginType   `json:"type"`
}

type AuthFlows struct {
	Flows []AuthFlow `json:"flows"`
}

type authRequest struct {
	Type     LoginType `json:"type"`
	Username string    `json:"user"`
	Password string    `json:"password"`
}

type authResponse struct {
	UserId      string `json:"user_id"`
	AccessToken string `json:"access_token"`
}

var defaultRegisterFlows = AuthFlows{
	Flows: []AuthFlow{
		{
			Stages: []LoginType{ // not implemented
				LoginTypeEmail,
				LoginTypePassword,
			},
			Type: LoginTypeEmail,
		},
		{Type: LoginTypePassword},
	},
}

var defaultLoginFlows = AuthFlows{
	Flows: []AuthFlow{
		{Type: LoginTypePassword},
	},
}

func registerWithPassword(hostname string, body *authRequest) interface{} {
	if body.Username == "" {
		return BadJsonError("Missing or invalid user")
	}
	if body.Password == "" {
		return BadJsonError("Missing or invalid password")
	}
	userId := fmt.Sprintf("@%s:%s", body.Username, hostname)
	user, err := service.CreateUser(userId)
	if err != nil {
		return UserInUseError(err.Error())
	}
	if err := user.SetPassword(body.Password); err != nil {
		return ServerError(err.Error())
	}
	accessToken, err := service.NewAccessToken(userId)
	if err != nil {
		return ServerError(err.Error())
	}
	return authResponse{
		UserId:      userId,
		AccessToken: accessToken,
	}
}

var registerResource = Resource{
	Get: NewJsonHandler(func(req *http.Request) interface{} {
		return &defaultRegisterFlows
	}),
	Post: NewJsonHandler(func(req *http.Request, body *authRequest) interface{} {
		switch body.Type {
		case LoginTypePassword:
			hostname := strings.Split(req.Host, ":")[0]
			return registerWithPassword(hostname, body)
		}
		return BadJsonError(fmt.Sprintf("Missing or invalid login type: '%s'", body.Type))
	}),
}

var loginResource = Resource{
	Get: NewJsonHandler(func(req *http.Request) interface{} {
		return &defaultLoginFlows
	}),
}
