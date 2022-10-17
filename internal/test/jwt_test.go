package test

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name string `json:"name""`
	jwt.StandardClaims
}

var myKey =[]byte("gin-gorm-oj-key")

//生成token
func TestGenerateToken(t *testing.T){
	UserClaim := &UserClaims{
		Identity: "user1",
		Name: "abc",
		StandardClaims:jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tokenString)
}

//解析token
func TestAnalyseToken(t *testing.T){
	tokenString :="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXIxIiwibmFtZSI6ImFiYyJ9.Jed6S6b14a7Fo9Vqyh9uakZecYfMTgWwmZJT50Rpbys"
	userClaim :=new(UserClaims)
	claims,err := jwt.ParseWithClaims(tokenString,userClaim,func(token *jwt.Token)(interface{},error){
		return myKey,nil
	})
	if err!=nil{
		t.Fatal(err)
	}
	if claims.Valid{
		fmt.Println(userClaim)
	}
}