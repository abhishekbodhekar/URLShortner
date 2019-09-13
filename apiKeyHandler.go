package main

import (
	"encoding/json"
	"time"

	"gopkg.in/redis.v4"
)

// This function creates at API key for userNAme passed in the argument
// It stores the key with its expiryTime and username, everyting encrypted to redis with "key:"<username>
// as a key
//
// returns empty string on failure
func createAPIKey(userName string) (string, string) {

	keyElement := APIKey{
		UserName: userName,
	}
	keyElement.ExpiryTime = time.Now().AddDate(0, 0, 2).Format(time.RFC3339)

	tenMinDuration, _ := time.ParseDuration("0h10m")

	keyElement.ExpiryTime = time.Now().Add(tenMinDuration).Format(time.RFC3339)

	keyElement.Key = encodeToB64(userName + keyElement.ExpiryTime)
	simpleKey := keyElement.Key
	simpleName := keyElement.UserName
	keyElement.ExpiryTime, _ = encrypt(keyElement.ExpiryTime)
	keyElement.UserName, _ = encrypt(keyElement.UserName)
	keyElement.Key, _ = encrypt(keyElement.Key)

	byteArr, err := json.Marshal(keyElement)
	if err == nil {
		_, err2 := redisClient.Set("key:"+userName, string(byteArr), 0).Result()
		if err2 == nil {
			return simpleName, simpleKey

		}
	}

	return "", ""

}

// It validated the key for the user passed in the argumanet and return okidiom
// It decypted every value stored at redis at checks for timestap with time.Now(), if valid
// returns true
func ValidateAPIKey(key string, userName string) bool {

	if key == "" || userName == "" {
		return false
	} else {

		data, err := redisClient.Get("key:" + userName).Result()

		if err == redis.Nil {
			return false

		} else if err != nil {

			return false
		} else {

			api := APIKey{}
			err2 := json.Unmarshal([]byte(data), &api)
			if err2 == nil {

				decryptedName, _ := decrypt(api.UserName)
				decryptedKey, _ := decrypt(api.Key)

				if decryptedKey == key && decryptedName == userName {
					return true
				}
				return false
			}
			return false
		}

	}
}
