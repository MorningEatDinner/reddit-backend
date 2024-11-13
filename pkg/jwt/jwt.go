package jwt

import (
	"errors"
	jwtpkg "github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

var mySecret = []byte("这是一个加盐句子")

type MyClaims struct {
	UserID   int64  `json:"userID"`
	Username string `json:"username"`
	jwtpkg.RegisteredClaims
}

func keyFunc(_ *jwtpkg.Token) (i interface{}, err error) {
	return mySecret, nil
}

// 生成JWT
// GenToken 生成JWT
func GenToken(userID int64, username string) (string, string, error) {
	// 创建一个我们自己的声明
	claims := MyClaims{
		userID,
		username, // 自定义字段
		jwtpkg.RegisteredClaims{
			ExpiresAt: jwtpkg.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "my-project", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	accessToken, err := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, claims).SignedString(mySecret)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	refreshToken, err := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, jwtpkg.RegisteredClaims{
		ExpiresAt: jwtpkg.NewNumericDate(time.Now().Add(time.Hour * 30)),
		Issuer:    "my-project",
	}).SignedString(mySecret)
	return accessToken, refreshToken, err
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	var mc = new(MyClaims)
	token, err := jwtpkg.ParseWithClaims(tokenString, mc, func(token *jwtpkg.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}

// RefreshToken: 刷新accessToken
func RefreshToken(at, rt string) (accessToken, refreshToken string, err error) {
	//  验证rt是否有效
	if _, err = jwtpkg.Parse(rt, keyFunc); err != nil {
		return
	}
	// 从旧的claim中解析出数据
	var claim MyClaims
	_, err = jwtpkg.ParseWithClaims(at, &claim, keyFunc)
	if strings.Contains(err.Error(), jwtpkg.ErrTokenExpired.Error()) { // 如果错误是过期所导致的， 如果不是则pass
		// 如果是过期了
		return GenToken(claim.UserID, claim.Username)
	}

	return
}
