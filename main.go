package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Msg struct {
	msg_type  string `json:"msg_Type"`
	timestamp int64
	sign      string
	content   Content
}
type Content struct {
	text string
}

func sendMsg(apiUrl, msg string, secret string) {
	// json
	contentType := "application/json"
	// data
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
		body, err := ioutil.ReadAll(resp.Body)
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
func main() {
	webhook := os.Getenv("PLUGIN_WEBHOOK")
	secret := os.Getenv("PLUGIN_SECRET")
	text := os.Getenv("DRONE_BUILD_LINK")
	println(webhook, text, secret)
	sendMsg(webhook, text, secret)
}
