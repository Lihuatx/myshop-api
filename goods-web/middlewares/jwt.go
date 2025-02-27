package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"myshop_api/goods-web/global"
	"myshop_api/goods-web/models"
	"net/http"
	"time"
)

// JWTAuth 是一个Gin框架的中间件函数，用于JWT（JSON Web Token）认证。
// 它返回一个Gin的HandlerFunc类型的函数，该函数在每次HTTP请求时执行JWT认证逻辑。
//
// 参数:
//
//	无
//
// 返回值:
//
//	gin.HandlerFunc: Gin框架的中间件函数
//
// 在这个中间件函数中，首先从请求的Header中获取名为"x-token"的字段，该字段包含JWT token。
// 如果token为空，则返回401 Unauthorized状态码和提示信息"请登录"。
// 如果token不为空，则尝试解析token。如果token解析失败，则根据错误类型返回不同的响应。
// 如果token解析成功，则将解析得到的claims信息存储在Gin的上下文中，并继续处理请求。
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// JWT鉴权逻辑开始
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		if token == "" {
			// 如果token为空，则返回未授权状态码和提示信息，并中止后续处理
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "请登录",
			})
			c.Abort()
			return
		}
		// 实例化JWT对象
		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			// 如果解析token出现错误
			if err == TokenExpired {
				// 如果错误类型为token过期
				if err == TokenExpired {
					// 返回授权已过期状态码和提示信息，并中止后续处理
					c.JSON(http.StatusUnauthorized, map[string]string{
						"msg": "授权已过期",
					})
					c.Abort()
					return
				}
			}

			// 如果错误不是token过期，则返回未登录状态码和提示信息，并中止后续处理
			c.JSON(http.StatusUnauthorized, "未登陆")
			c.Abort()
			return
		}
		// 将解析得到的claims信息存储在Gin的上下文中
		c.Set("claims", claims)
		// 将用户的ID也存储在Gin的上下文中
		c.Set("userId", claims.ID)
		// 继续处理后续的中间件或处理器
		c.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConfig.JWTInfo.SigningKey), //可以设置过期时间
	}
}

// 创建一个token
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}

}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
