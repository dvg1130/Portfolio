package authdb

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dvg1130/Portfolio/secure-auth/authsdk"
	"github.com/dvg1130/Portfolio/secure-backend/config"
	_ "github.com/go-sql-driver/mysql"
)

func AuthService() (sso *authsdk.SDKClient, err error) {
	// var base = config
	sso = authsdk.NewClient(config.AuthService.SSO_ADDR)

	resp, err := http.Get(config.AuthService.SSO_ADDR + "/health")
	if err != nil {
		fmt.Println("error reaching service at: ", config.AuthService.SSO_ADDR)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Println("Auth service is healthy")
	} else {
		log.Println("Auth service returned status:", resp.StatusCode)
	}
	return sso, nil

}
