package api

import (
	"context"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	console "github.com/alibabacloud-go/tea-console/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"math/rand"
	"myshop_api/user-web/forms"
	"net/http"
	"strings"
	"time"

	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"myshop_api/user-web/global"
)

const (
	mnsDomain = "1943695596114318.mns.cn-hangzhou.aliyuncs.com"
)

func GenerateSmsCode(witdh int) string {
	//生成width长度的短信验证码

	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < witdh; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

// 使用AK&SK初始化账号Client
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi.Client, _err error) {
	config := &openapi.Config{}
	config.AccessKeyId = accessKeyId
	config.AccessKeySecret = accessKeySecret
	_result = &dysmsapi.Client{}
	_result, _err = dysmsapi.NewClient(config)
	return _result, _err
}

func SendSms(ctx *gin.Context) {
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	smsCode := GenerateSmsCode(6)
	client, _err := CreateClient(&global.ServerConfig.AliSmsInfo.ApiKey, &global.ServerConfig.AliSmsInfo.ApiSecrect)

	if _err != nil {
		return
	}
	SignName := "广州慕学服务"
	TemplateCode := "SMS_474916068"
	TemplateParam := "{\"code\":" + smsCode + "}"
	// 1.发送短信
	sendReq := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  &sendSmsForm.Mobile,
		SignName:      &SignName,
		TemplateCode:  &TemplateCode,
		TemplateParam: &TemplateParam,
	}
	sendResp, _err := client.SendSms(sendReq)
	if _err != nil {
		return
	}
	code := sendResp.Body.Code
	if !tea.BoolValue(util.EqualString(code, tea.String("OK"))) {
		console.Log(tea.String("错误信息: " + tea.StringValue(sendResp.Body.Message)))
		return
	}

	//将验证码保存起来 - redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})

	rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
