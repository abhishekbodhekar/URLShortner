package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"gopkg.in/redis.v4"
)

// This function finds for the "lastIndex" strored in redis,
//  increments it by 1, and returns it
//
// if not found, it creates it in redis and SET it to 1
//
// **NOTE** If "lastIndex" is SET other than int value, it DEL (deletes)
// it and SET "lastIndex" to 1
func getNewIndexToPutURL() int {

	lastIndex, err := redisClient.Get("lastIndex").Result()
	if err == redis.Nil {

		lastIndex = strconv.Itoa(1)
		err2 := redisClient.Set("lastIndex", lastIndex, 0).Err()
		if err2 != nil {
			return 0
		}
		return 1

	} else if err != nil {

		return 0
	} else {
		newIndex, err := strconv.Atoi(lastIndex)

		if err != nil {
			redisClient.Del("lastIndex")
			newIndex = 0
		}

		newIndex = newIndex + 1
		err3 := redisClient.Set("lastIndex", newIndex, 0).Err()
		if err3 != nil {
			return 0
		}

		return newIndex

	}

}

// This function puts an url to redis specified by index.
// It encrypts teh url value before putting it in redis.
//
//**NOTE** They key for the url is the combination of "link:"<index>
func putNewURL(url string, newIndex int) bool {

	_, err := redisClient.RPush("IdList", newIndex).Result()
	if err == nil {

		valueToPut := hashValue{
			Link:  url,
			Count: "0",
		}

		valueToPut.Count, _ = encrypt(valueToPut.Count)

		byteArr, newErr := json.Marshal(valueToPut)
		if newErr != nil {
			return false
		}
		newErr2 := redisClient.Set("link:"+strconv.Itoa(newIndex), string(byteArr), 0).Err()

		if newErr2 != nil {
			return false
		}

		return true
	}
	return false

}

//This function Check the url with all stored urls in the redis.
// It return "0" if the url is not found
func checkIfUrlAvailable(url string) string {

	idArr, err := redisClient.LRange("IdList", 0, -1).Result()

	if err == nil {
		for _, val := range idArr {

			strVal, err2 := redisClient.Get("link:" + val).Result()
			if err2 == nil {
				hashVal := hashValue{}
				err3 := json.Unmarshal([]byte(strVal), &hashVal)

				if err3 == nil {
					decryptedURL, _ := decrypt(hashVal.Link)

					if decryptedURL == url {
						return val
					}
				} else {
					return "0"
				}

			} else {
				return "0"
			}

		}
		return "0"
	}
	return "0"

}

//This fucntion gets short list of the actual url passed in the argument.
//
// It firtsly checks whether the corresponding longURL is present, get the index of it and encodeTo64 to get short url.
//
// Otherwise, It gets a new key comination with new index encode it and returns the encoded short URL.
//
//**NOTE** return of empty string with another return value as false means the functions has. failed
func getShortURL(longURL string) (string, bool) {

	encrytptedURL, err := encrypt(longURL)
	if err == nil {
		indexOfURL := checkIfUrlAvailable(longURL)
		if indexOfURL != "0" {
			fmt.Println(" available")
			shortURL := encodeToB64(indexOfURL)
			return shortURL, true
		} else {
			fmt.Println("not available")

			newIndex := getNewIndexToPutURL()
			ok := putNewURL(encrytptedURL, newIndex)
			if ok {
				newIndexStr := strconv.Itoa(newIndex)
				shortURL := encodeToB64(newIndexStr)
				return shortURL, true
			}
			return "", false
		}

	} else {
		fmt.Println("error here")

		return "", false
	}

}

// It returns the long url with true, of the short url provided in argument.
//  On succession, It gets the hasvVlaue stored at respective shortURL, decrypts it and sends shortURL contained in it.
func getLongURL(shortURL string) (string, bool) {
	index, ok := decodeFromB64(shortURL)

	if ok {
		hash, okk := getHashStoredAtIndex(index)
		if okk == false {
			return "Something went wrong OR no short link is not present yet", false

		} else if okk {
			decryptedValue, _ := decrypt(hash.Link)
			go incrementCounterAndSave(hash, index)
			return decryptedValue, true

		}
		return "Something went wrong OR no short link is not present yet", false
	} else {
		return "Something went wrong OR no short link is not present yet", false
	}
}

// It returns the counter (#visits) with true, of the short url provided in argument.
//  On succession, It gets the hasvVlaue stored at respective shortURL, decrypts it and sends cunter contained in it.
func getCounter(shortURL string) (string, bool) {
	index, ok := decodeFromB64(shortURL)

	if ok {
		hash, okk := getHashStoredAtIndex(index)
		if okk == false {
			return "Something went wrong OR no short link is not present yet", false

		} else if okk {
			decryptedValue, _ := decrypt(hash.Count)
			return decryptedValue, true

		}
		return "Something went wrong OR no short link is not present yet", false
	} else {
		return "Something went wrong OR no short link is not present yet", false
	}
}

// It gets the hashval (struct) and true stored at the index
// **NOTE** empty hashVal{} or flase in return statemet indicates that the
// function failed or no hash was present in specified location
func getHashStoredAtIndex(index string) (hashValue, bool) {

	hashVal, err := redisClient.Get("link:" + index).Result()
	if err == redis.Nil {

		return hashValue{}, false
	} else if err != nil {

		return hashValue{}, false

	} else {

		receievdHash := hashValue{}
		err2 := json.Unmarshal([]byte(hashVal), &receievdHash)
		if err2 == nil {

			return receievdHash, true
		}

		return hashValue{}, false
	}
}

//It increments the counter present in hashVal at specified index
func incrementCounterAndSave(hval hashValue, index string) {

	decryptedCount, _ := decrypt(hval.Count)
	newCounter, _ := strconv.Atoi(decryptedCount)
	newCounter = newCounter + 1
	hval.Count = strconv.Itoa(newCounter)
	hval.Count, _ = encrypt(hval.Count)

	byteArr, newErr := json.Marshal(hval)
	if newErr != nil {
		return
	}
	redisClient.Set("link:"+index, string(byteArr), 0).Err()

}
