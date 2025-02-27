package global

import (
	ut "github.com/go-playground/universal-translator"
	"myshop_api/user-web/config"
	"myshop_api/user-web/proto"
)

var (
	Trans ut.Translator

	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	UserSrvClient proto.UserClient
)
