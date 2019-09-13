package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/redis.v4"
)

var redisClient *redis.Client

func main() {
	initializeRedis()
	CreateAndStartServer()

}

// This starts the server on "APP_URL" environment variable and
// uses default *ServerMux for routing.
func CreateAndStartServer() {

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/getShortLink", onGetShortLink)

	http.HandleFunc("/getRedirectLink", onGetRedirectLink)
	http.HandleFunc("/getVisits", onGetVisits)
	http.HandleFunc("/registerNewKey", onRegisterNewKey)
	http.ListenAndServe(os.Getenv("APP_URL"), nil)

}

// This starts the redis server on "REDIS_URL"
func initializeRedis() {

	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})
	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)

}

// Homepage "/" hadler
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	responseToCLient(w, "Nothing is here... \n\nUse following apis\n\n/getShortLink\n/getRedirectLink\n/getVisits\n/registerNewKey")
}

// GET request handler on "/GetShortLink"
//
// QueryParam for /GetShortLink:
//  - "longURL" : "value"
//
// **NOTE** This functions requires a correct URL as value
func onGetShortLink(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	longURL, ok := values["longURL"]
	if ok {
		w.WriteHeader(http.StatusOK)
		if len(longURL) >= 1 {
			fmt.Println(longURL[0])
			_, err := url.ParseRequestURI(longURL[0])

			if err != nil {
				responseToCLient(w, "Please enter the correct and complete url, example - http://google.com")

			} else {

				shortUrl, ok := getShortURL(longURL[0])
				if ok {

					responseToCLient(w, "Your Short URL is : http://mydomain.com/"+shortUrl)
				} else {
					responseToCLient(w, "Please check the request parameters")

				}
			}
		} else {
			responseToCLient(w, "Please check the request parameters")
		}
	} else {
		responseToCLient(w, "No longURL found, Please check the request parameters")

	}

}

// GET request handler on "/GetRedirectLink"
//
// QueryParam for /GetRedirectLink:
//  - "shortURL" : "value"
//
// **NOTE** This functions requires a correct URL as value
func onGetRedirectLink(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	shortURL, ok := values["shortURL"]

	if ok {
		w.WriteHeader(http.StatusOK)
		if len(shortURL) >= 1 {
			correctURL, err := url.ParseRequestURI(shortURL[0])

			if err != nil {
				responseToCLient(w, "Please enter the correct and complete url, example - http://google.com")

			} else {

				if correctURL.Host != "mydomain.com" {
					responseToCLient(w, "Not the correct short link provided by mydomain.com")
				} else {
					a := correctURL.Path[1:]
					str, _ := getLongURL(a)
					responseToCLient(w, str)
				}
			}
		} else {
			responseToCLient(w, "Please check the request parameters")
		}
	} else {
		responseToCLient(w, "No shortURL found, Please check the request parameters")

	}

}

// GET request handler on "/GetShortLink"
//
// QueryParam for /GetShortLink:
//  - "longURL" : "URLvalue"
//	- "key" : "KeyValue"
//	- "userName" : "userNameValue"
//
// **NOTE** This functions requires a correct URL as URLvalue
func onGetVisits(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	shortURL, ok := values["shortURL"]

	key, ok1 := values["key"]
	userName, ok2 := values["userName"]

	if (ok1 == true && ok2 == true) && (len(key) >= 1) && (len(userName) >= 1) {

		if ValidateAPIKey(key[0], userName[0]) == false {
			responseToCLient(w, "Wrong or expired key")
			return

		}

	} else {
		responseToCLient(w, "Please check the request parameters")
		return

	}

	if ok {
		w.WriteHeader(http.StatusOK)
		if len(shortURL) >= 1 {
			correctURL, err := url.ParseRequestURI(shortURL[0])

			if err != nil {
				responseToCLient(w, "Please enter the correct and complete url, example - http://google.com")

			} else {
				fmt.Println("host " + correctURL.Host)
				if correctURL.Host != "mydomain.com" {
					responseToCLient(w, "Not the correct short link provided by mydomain.com")
				} else {
					a := correctURL.Path[1:] // **NOTE** This functions requires a correct URL as URLvalue
					str, _ := getCounter(a)
					responseToCLient(w, str)
				}
			}
		} else {
			responseToCLient(w, "Please check the request parameters")
		}
	} else {
		responseToCLient(w, "No shortURL found, Please check the request parameters")

	}

}

//This common response is used for every request
func responseToCLient(w http.ResponseWriter, str string) {

	w.Write([]byte(str))

}

// GET request handler on "/RegisterNewKey"
//
// QueryParam for /RegisterNewKey:
//	- "userName" : "userNameValue"
//
func onRegisterNewKey(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	name, ok := values["userName"]
	if ok {
		w.WriteHeader(http.StatusOK)
		if len(name) >= 1 {

			uName, key := createAPIKey(name[0])

			if uName != "" && key != "" {
				responseToCLient(w, "userName : "+uName+"\nkey : "+key)
			} else {
				responseToCLient(w, "key generation failed")

			}
		}
	}
}