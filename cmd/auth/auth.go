package auth

import (
	"bufio"
	"fmt"
	vkapi "github.com/SevereCloud/vksdk/5.92/api"
	"github.com/justblender/vk-api"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
)

const (
	serviceToken string = "66619e0066619e0066d3e34c266634f6666666166619e003ea8d033c12d1a3d08e6fd55"
)

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

	client := vkapi.Init(userToken)
	return client
}

func getUserToken() string {
	user, password := getUserCredentials()
	client, err := vk_api.NewClient(vk_api.DirectAuthentication{
		Username: user,
		Password: password,
		Device:   vk_api.ANDROID,
	})
	if err != nil {
		panic("Couldn't authenticate with provided credentials.")
	}
	return client.AccessToken
}

func getUserCredentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Printf("\n\n")
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password)
}
