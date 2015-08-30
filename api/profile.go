package api

import (
	"net/http"

	"github.com/Rugvip/bullettime/interfaces"
	"github.com/Rugvip/bullettime/types"

	"github.com/julienschmidt/httprouter"
)

type avatarUrlRequest struct {
	AvatarUrl *string `json:"avatar_url"`
}
type avatarUrlResponse struct {
	AvatarUrl string `json:"avatar_url"`
}
type displayNameRequest struct {
	DisplayName *string `json:"displayname"`
}
type displayNameResponse struct {
	DisplayName string `json:"displayname"`
}

func (e profileEndpoint) getDisplayName(params httprouter.Params) interface{} {
	user, err := e.userFromParams(params)
	if err != nil {
		return err
	}
	profile, err := e.profiles.Profile(user, user)
	if err != nil {
		return err
	}
	return displayNameResponse{profile.DisplayName}
}

func (e profileEndpoint) setDisplayName(req *http.Request, params httprouter.Params, body *displayNameRequest) interface{} {
	authedUser, err := readAccessToken(e.users, e.tokens, req)
	if err != nil {
		return err
	}
	user, err := e.userFromParams(params)
	if err != nil {
		return err
	}
	if body.DisplayName == nil {
		return types.BadJsonError("missing 'displayname'")
	}
	if _, err := e.profiles.UpdateProfile(user, authedUser, body.DisplayName, nil); err != nil {
		return err
	}
	return struct{}{}
}

func (e profileEndpoint) getAvatarUrl(params httprouter.Params) interface{} {
	user, err := e.userFromParams(params)
	if err != nil {
		return err
	}
	profile, err := e.profiles.Profile(user, user)
	if err != nil {
		return err
	}
	return avatarUrlResponse{profile.AvatarUrl}
}

func (e profileEndpoint) setAvatarUrl(req *http.Request, params httprouter.Params, body *avatarUrlRequest) interface{} {
	authedUser, err := readAccessToken(e.users, e.tokens, req)
	if err != nil {
		return err
	}
	user, err := e.userFromParams(params)
	if err != nil {
		return err
	}
	if body.AvatarUrl == nil {
		return types.BadJsonError("missing 'avatar_url'")
	}
	if _, err := e.profiles.UpdateProfile(user, authedUser, nil, body.AvatarUrl); err != nil {
		return err
	}
	return struct{}{}
}

func (e profileEndpoint) userFromParams(params httprouter.Params) (types.UserId, types.Error) {
	user, err := types.ParseUserId(params[0].Value)
	if err != nil {
		return types.UserId{}, types.BadJsonError(err.Error())
	}
	if err := e.users.UserExists(user, user); err != nil {
		return types.UserId{}, err
	}
	return user, nil
}

func (e profileEndpoint) Register(mux *httprouter.Router) {
	mux.GET("/profile/:userId/displayname", jsonHandler(e.getDisplayName))
	mux.PUT("/profile/:userId/displayname", jsonHandler(e.setDisplayName))
	mux.GET("/profile/:userId/avatar_url", jsonHandler(e.getAvatarUrl))
	mux.PUT("/profile/:userId/avatar_url", jsonHandler(e.setAvatarUrl))
}

type profileEndpoint struct {
	users    interfaces.UserService
	tokens   interfaces.TokenService
	profiles interfaces.ProfileService
}

func NewProfileEndpoint(
	users interfaces.UserService,
	tokens interfaces.TokenService,
	profiles interfaces.ProfileService,
) Endpoint {
	return profileEndpoint{
		users,
		tokens,
		profiles,
	}
}
