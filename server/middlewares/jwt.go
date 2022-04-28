package middlewares

import (
	"crypto/rsa"
	"fmt"
	"github.com/gin-gonic/gin"
	"im/global"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("w-token")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
				"msg": "请登录",
			})
			return
		}
		file,_ := os.Open(global.SrvConfig.JWTInfo.PublicKeyPath)
		bytes,_ := ioutil.ReadAll(file)
		publickey,_ := jwt.ParseRSAPublicKeyFromPEM(bytes)
		jwtVerier := &JWTTokenVerifier{PublicKey: publickey}
		claim,err := jwtVerier.Verify(token)
		if err != nil {
			fmt.Println(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
				"msg": "请登录",
			})
			return
		}
		ctx.Set("claim",claim)
		ctx.Set("name",claim.Subject)
		ctx.Next()
	}
}

func Auth(token string) (*MyClaim, error){
	if token == "" {
		return nil,fmt.Errorf("ws认证失败，token为空")
	}
	file,_ := os.Open(global.SrvConfig.JWTInfo.PublicKeyPath)
	bytes,_ := ioutil.ReadAll(file)
	publickey,_ := jwt.ParseRSAPublicKeyFromPEM(bytes)
	jwtVerier := &JWTTokenVerifier{PublicKey: publickey}
	return jwtVerier.Verify(token)
}

type JWTTokenVerifier struct {
	PublicKey *rsa.PublicKey
}

func (v *JWTTokenVerifier) Verify(token string) (*MyClaim, error) {
	t, err := jwt.ParseWithClaims(token, &MyClaim{},
		func(*jwt.Token) (interface{}, error) {
			return v.PublicKey, nil
		})

	if err != nil {
		return nil, fmt.Errorf("cannot parse token: %v", err)
	}

	if !t.Valid {
		return nil, fmt.Errorf("token not valid")
	}

	clm, ok := t.Claims.(*MyClaim)
	if !ok {
		return nil, fmt.Errorf("token claim is not MyClaim")
	}

	if err := clm.Valid(); err != nil {
		return nil, fmt.Errorf("claim not valid: %v", err)
	}
	return clm,nil

}

type JWTTokenGen struct {
	privateKey *rsa.PrivateKey
	issuer     string
	nowFunc    func() time.Time
}

func NewJWTTokenGen(issuer string, privateKey *rsa.PrivateKey) *JWTTokenGen {
	return &JWTTokenGen{
		issuer:     issuer,
		nowFunc:    time.Now,
		privateKey: privateKey,
	}
}

type MyClaim struct {
	Role int
	jwt.StandardClaims
}
func (t *JWTTokenGen) GenerateToken(userName string, expire time.Duration) (string, error) {
	nowSec := t.nowFunc().Unix()
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS512, &MyClaim{
		StandardClaims: jwt.StandardClaims{
			Issuer:    t.issuer,
			IssuedAt:  nowSec,
			ExpiresAt: nowSec + int64(expire.Seconds()),
			Subject:   userName,
		},
	})
	return tkn.SignedString(t.privateKey)
}