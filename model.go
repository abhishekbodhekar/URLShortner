package main

type hashValue struct { // this model maps with the  json value regarding URL object stored in redis
	Link  string
	Count string
}
type APIKey struct { //this model maps with the json value regarding API object stored in redis
	UserName   string
	Key        string
	ExpiryTime string
}
