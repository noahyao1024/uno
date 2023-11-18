package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"uno/pkg/setting"

	gohttpclient "github.com/bozd4g/go-http-client"
	"github.com/tidwall/gjson"
)

type Entry struct {
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
}

func (e *Entry) token() (string, error) {
	apiCfg, ok := setting.AppInstance.APIs["access_token"]
	if !ok {
		return "", fmt.Errorf("access_token api not found")
	}

	ctx := context.Background()
	client := gohttpclient.New(apiCfg.Host)

	body, _ := json.Marshal(map[string]interface{}{
		"app_key":    apiCfg.AppID,
		"app_secret": apiCfg.AppSecret,
	})

	data := make(map[string]interface{}, 0)

	response, err := client.Post(ctx, apiCfg.Path, gohttpclient.WithBody(body))
	response.Unmarshal(&data)
	fmt.Println(err, data)

	accessToken, ok := data["data"].(map[string]interface{})["access_token"]
	if !ok {
		return "", fmt.Errorf("access_token not found")
	}

	token := fmt.Sprintf("%v", accessToken)

	// TODO add cache.

	return token, nil
}

func (e *Entry) Detail() error {
	token, err := e.token()
	fmt.Println(token, err)
	if err != nil {
		return err
	}

	apiCfg, ok := setting.AppInstance.APIs["subscriber_info"]
	fmt.Println(apiCfg)
	if !ok {
		return fmt.Errorf("subscriber_info api not found")
	}

	ctx := context.Background()
	client := gohttpclient.New(apiCfg.Host)

	userID, _ := strconv.ParseInt(e.UserID, 10, 64)

	body, _ := json.Marshal(map[string]interface{}{
		"user_id": userID,
	})

	response, err := client.Post(ctx, apiCfg.Path, gohttpclient.WithBody(body), gohttpclient.WithHeader("Authorization", token))

	for _, item := range gjson.GetBytes(response.Body(), "data.bind_account_list").Array() {
		switch gjson.Get(item.Raw, "identify_type").String() {
		case "2":
			email := gjson.Get(item.Raw, "identifier").String()
			if len(email) == 0 {
				return fmt.Errorf("email not found")
			}

			e.Email = email

			// Only get the first email.
			break
		}
	}

	return nil
}
