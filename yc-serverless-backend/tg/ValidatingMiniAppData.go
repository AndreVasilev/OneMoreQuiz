package tg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
	"yc-serverless-backend/models"
)

func Validate(initData string) *models.WebAppInitData {
	query, err := url.ParseQuery(initData)
	if err != nil {
		log.Printf(err.Error())
		return nil
	}

	dataCheckPairs := []string{}
	var hashValue string
	for k, v := range query {
		if k == "hash" {
			hashValue = v[0]
		} else {
			pair := k + "=" + v[0]
			dataCheckPairs = append(dataCheckPairs, pair)
		}
	}

	sort.Strings(dataCheckPairs)
	payload := strings.Join(dataCheckPairs, "\n")
	botToken := os.Getenv("BOT_TOKEN")
	hash := calcHash("WebAppData", botToken, payload)

	if hash == hashValue {
		return NewWebAppInitData(query)
	} else {
		return nil
	}
}

func calcHash(key string, botToken string, payload string) string {
	secretHmac := hmac.New(sha256.New, []byte(key))
	secretHmac.Write([]byte(botToken))
	secret := secretHmac.Sum(nil)
	hashHmac := hmac.New(sha256.New, secret)
	hashHmac.Write([]byte(payload))
	hash := hex.EncodeToString(hashHmac.Sum(nil))
	return hash
}
