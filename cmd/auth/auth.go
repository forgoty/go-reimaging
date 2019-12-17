package auth

import (
	"encoding/json"
	// "fmt"
	vkapi "github.com/SevereCloud/vksdk/5.92/api"
	"io/ioutil"
	"net/http"
)

const (
	serviceToken string = "66619e0066619e0066d3e34c266634f6666666166619e003ea8d033c12d1a3d08e6fd55"
	tokenUrl     string = "https://oauth.vk.com/access_token"
	clientSecret string = "VeWdmVclDCtn6ihuP1nt"
	clientID     string = "5597286"
)

type Token struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	UserID           int    `json:"user_id"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (t *Token) String() string {
	return t.AccessToken
}

//const AppId string = "5597286"

func GetClient(auth bool) *vkapi.VK {
	if auth {
		return GetUserClient()
	}
	return GetServiceClient()
}

func GetServiceClient() *vkapi.VK {
	client := vkapi.Init(serviceToken)
	return client
}

func GetUserClient() *vkapi.VK {
	userToken := getUserToken()

	client := vkapi.Init(userToken.String())
	return client
}

func getUserToken() Token {
	user, password := getUserCredentials()
	client := http.Client{}
	req, error := http.NewRequest("GET", tokenUrl, nil)
	if error != nil {
		return Token{}
	}
	query := req.URL.Query()
	query.Add("grant_type", "password")
	query.Add("client_id", clientID)
	query.Add("client_secret", clientSecret)
	query.Add("username", user)
	query.Add("password", password)
	query.Add("v", "5.92")
	req.URL.RawQuery = query.Encode()

	response, error := client.Do(req)
	if error != nil {
		return Token{}
	}
	defer response.Body.Close()

	body, error := ioutil.ReadAll(response.Body)
	if error != nil {
		return Token{}
	}

	var token Token
	json.Unmarshal(body, &token)
	if token.Error != "" {
		return Token{}
	}
	return token
}

func getUserCredentials() (login, password string) {
	return "", "password"
}
