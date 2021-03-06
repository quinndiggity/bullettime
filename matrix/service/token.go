// Copyright 2015  Ericsson AB
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"encoding/base64"
	"fmt"
	"strings"

	ct "github.com/matrix-org/bullettime/core/types"
	"github.com/matrix-org/bullettime/matrix/interfaces"
	"github.com/matrix-org/bullettime/matrix/types"
	"github.com/matrix-org/bullettime/utils"
)

func CreateTokenService() (interfaces.TokenService, error) {
	return tokenService{}, nil
}

type tokenService struct{}

type tokenInfo struct {
	userId ct.UserId
}

func (t tokenInfo) String() string {
	encodedUserId := base64.RawURLEncoding.EncodeToString([]byte(t.userId.String()))
	return fmt.Sprintf("%s..%s", encodedUserId, utils.RandomString(16))
}

func (t tokenInfo) UserId() ct.UserId {
	return t.userId
}

func (t tokenService) NewAccessToken(userId ct.UserId) (interfaces.Token, types.Error) {
	return tokenInfo{userId}, nil
}

func (t tokenService) ParseAccessToken(token string) (interfaces.Token, types.Error) {
	splits := strings.Split(token, "..")
	if len(splits) != 2 {
		return nil, types.DefaultUnknownTokenError
	}
	userIdStr, err := base64.RawURLEncoding.DecodeString(splits[0])
	if err != nil {
		return nil, types.DefaultUnknownTokenError
	}
	userId, err := ct.ParseUserId(string(userIdStr))
	if err != nil {
		return nil, types.DefaultUnknownTokenError
	}
	return tokenInfo{userId}, nil
}
