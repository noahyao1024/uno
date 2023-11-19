package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"uno/pkg/database"
	"uno/pkg/setting"

	gohttpclient "github.com/bozd4g/go-http-client"
	"github.com/patrickmn/go-cache"
	"github.com/tidwall/gjson"
)

var lc = cache.New(5*time.Minute, 10*time.Minute)

type Entry struct {
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
}

func (e *Entry) TableName() string {
	return "subscriber"
}

func (e *Entry) token() (string, error) {
	apiCfg, ok := setting.AppInstance.APIs["access_token"]
	if !ok {
		return "", fmt.Errorf("access_token api not found")
	}

	key := fmt.Sprintf("access_token:%s", apiCfg.AppID)
	if accessToken, exists := lc.Get(key); exists {
		return accessToken.(string), nil
	}

	ctx := context.Background()
	client := gohttpclient.New(apiCfg.Host)

	body, _ := json.Marshal(map[string]interface{}{
		"app_key":    apiCfg.AppID,
		"app_secret": apiCfg.AppSecret,
	})

	data := make(map[string]interface{}, 0)

	response, err := client.Post(ctx, apiCfg.Path, gohttpclient.WithBody(body))
	if err != nil {
		return "", err
	}

	response.Unmarshal(&data)

	accessToken, ok := data["data"].(map[string]interface{})["access_token"]
	if !ok {
		return "", fmt.Errorf("access_token not found")
	}

	token := fmt.Sprintf("%v", accessToken)

	lc.Set(key, token, time.Second*3500)

	return token, nil
}

func (e *Entry) Detail() error {
	dbEntry := &Entry{}
	database.GetWriteDB().Where("user_id = ?", e.UserID).First(dbEntry)
	if dbEntry.Email != "" {
		e.Email = dbEntry.Email
		return nil
	}

	token, err := e.token()
	if err != nil {
		return err
	}

	apiCfg, ok := setting.AppInstance.APIs["subscriber_info"]
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

			if database.GetWriteDB().Create(e).Error != nil {
				return fmt.Errorf("create subscriber failed")
			}

			// Only get the first email.
			break
		}
	}

	return nil
}
