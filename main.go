package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Msg struct {
	msg_type  string
	timestamp int64
	sign      string
	content   Content
}
type Content struct {
	text string
}

func sendMsg(apiUrl, msg string, secret string) {
	contentType := "application/json"
	currentTime := time.Now().Unix()
	sign, _ := GenSign(secret, currentTime)
	sendData := `{
		"timestamp": ` + strconv.FormatInt(currentTime, 10) + `,
		"sign": "` + sign + `",
		"msg_type": "text",
		"content": {"text": "` + msg + `"}
	}`
	resp, err := http.Post(apiUrl, contentType, strings.NewReader(sendData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		jsonStr := string(body)
		fmt.Println("Response: ", jsonStr)

	} else {
		fmt.Println("Get failed with error: ", resp.Status)
	}
}
func GenSign(secret string, timestamp int64) (string, error) {
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret

	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
func sendCard(apiUrl, secret string) {
	contentType := "application/json"
	currentTime := time.Now().Unix()
	sign, _ := GenSign(secret, currentTime)
	DroneBuildLink := os.Getenv("DRONE_BUILD_LINK")
	DroneCommitLink := os.Getenv("DRONE_COMMIT_LINK")
	DroneRepoLink := os.Getenv("DRONE_REPO_LINK")
	DroneRepo := os.Getenv("DRONE_REPO")
	sendData := `{
	"timestamp": ` + strconv.FormatInt(currentTime, 10) + `,
	"sign": "` + sign + `",
	"msg_type": "interactive",
    "card":{
        "config":{
            "wide_screen_mode":true
        },
        "elements":[
            {
                "tag":"hr"
            },
            {
                "tag":"div",
                "text":{
                    "content":"- [` + DroneRepo + `](` + DroneRepoLink + `)\n- [commit](` + DroneCommitLink + `)\n- [drone](` + DroneBuildLink + `)",
                    "tag":"lark_md"
                }
            }
        ],
        "header":{
            "template":"blue",
            "title":{
                "content":"drone构建通知",
                "tag":"plain_text"
            }
        }
    }
}`
	resp, err := http.Post(apiUrl, contentType, strings.NewReader(sendData))
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	if resp.StatusCode == http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		jsonStr := string(body)
		fmt.Println("Response: ", jsonStr)
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Get failed with error: ", resp.Status, string(body))
	}
}
func main() {
	webhook := os.Getenv("PLUGIN_WEBHOOK")
	secret := os.Getenv("PLUGIN_SECRET")
	text := os.Getenv("PLUGIN_TEXT")
	if webhook == "" {
		fmt.Printf("请配置%s\n", "webhook")
		os.Exit(1)
	}
	if secret == "" {
		fmt.Printf("请配置%s\n", "secret")
		os.Exit(1)
	}
	if text != "" {
		sendMsg(webhook, text, secret)
	} else {
		sendCard(webhook, secret)
	}
}
