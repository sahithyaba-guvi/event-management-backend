package loginRespo

import (
	commonRespo "em_backend/responses/common"
)

type LoginResp struct {
	commonRespo.Response
	UserInfo commonRespo.LoginDetails `json:"userInfo"`
}
