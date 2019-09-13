# URL Shortner

This is an implementation of URL shortning

# Description

This is a package contains implementation for shortning the URLs.

The golang is used on application layer. Redis is used for database.

A docker-compose is provided as well.


It inculdes every data stored in redis completely incrypted with AES ecryption.

Also, It stores the number of visits made to any short link.

An API key is required to get the # visits.

It provides an API to generate an API key. Access to generate API key is currently unrestricted.


### Prerequisites

Docker with compose is needed to use with containerization.
Otherwise, only redis server is needed to be installed on host as compiled Go package itself is an executable.

## Implementation

HOW IT WORKS - 


### Get a long (Redirect) URL -> store in redis with some id -> encode the id -> get a short URL  

### Get a short URL -> decode to get the id -> in redis, get the value at id which is long URL

# Brief - 
 * Storing the longURL for the first time and fetching the ShortURL
1) The lastIndex is stored at redis. It is basically the identification index of the last inserted URL.
2) When a new URL request comes, The lastIndex if fetched. It is then incremented with 1.
3) It is then pushed to a list (IdList) in redis. (That is how, we have all indices)
3) A new key is made ("key:"+lasIndex).
4) The longURL with visitCount as zero is encrypted and stored at above key.
5) Now, the lastIndex (Which is associated with the current longURL) is encoded with base64.
6) This encoded string would be the path of our ShortURL
7) adding domain to this ShortURL path gives complete shortURL

* Getting shortURL of already present URL

1) The list containg all indices is iterated
2) All the indices are checked with forming a key ("key:"+currentIndex)
3) The values (as hash) on those keys are GET and the encrypted URL contained in it is valueated.
4) The longURL value we got is decrypted and checked against the longURL we have (which is to be converted to short).
5) If matched, the index is returned.
6) Now, the index is encoded with base64.
7) This encoded string would be the path of our ShortURL
8) adding domain to this ShortURL path gives complete shortURL

* Getting Redirect (LongURL) from short

1) Firstly, the short URL is decoded with base64 to get the index
2) The key is formed ("key"+index) to get the value (hash)
3) This value is decrypted to get the value of longURL

### Contents

1) vendor - For external dependancy for redis driver

2) docker-compose.yml

3) aesEncrption.go - This file intents to include all encryption and decryption methods. A symetric AESalgorithm is used.The Key is also provided statically inthe file.

4) apiKeyHandler.go - This file intents to include all API key creations and validation methods. Keys are stored in redis.

5) encoder.go - This file consists encoding and decoding to base64 methods.

6) model.go - This file contents models to map with redis for accessing urls and keys

7) REAMME.md - this is a readme file.

8) server.go - This file handles server handling operations. It consists of methods which starts application server and redis server as well. All the web routing and responses are handled in this file.

9) shortner.go - This file consits of the logic which drives this application. Fetching, storing to redis, operating on data and their values can be seen here.

10) Encryption_test.go - The test file. It tests the decryption algorithm.

11) URLShortner (executable) - This is the executable for aplication.

### Usage

at root, run with 
$ docker-composer up

* 1) localhost:5899/ 
    - This is the homepage
        Noting is here, just a static page

* 2) /getShortLink/

    - Desciprion : GET reuest, This prints short URL for provided long URL.

        queryParam : 

        1) "longURL"

    example - localhost:5899/getShortLink?longURL=http://google.com

    (prints short URL)

* 3) /getRedirectLink/

    - Desciption : GET request, This prints Redirect (long) URL for provided short URL

        queryPatam :

        1) "shortURL"

    example - localhost:5899/getRedirectLink?shortURL=http://mydomain.com/NQ== 

    (prints long (redirect) URL)

* 4) /registerNewKey/

    - Desciption : GET request, This prints the API key for the userName provided. REMEMBER, this key is valid for only 10 minutes. You can create new key for the same userName anytime.

        queryParam : 

        1) "userName"

    example - localhost:5899/registerNewKey?userName=ggsdsdf

    (prints userName and API key)

* 5) /getVisits/

    - Desciption : GET request, This prints the number of visists made to shortURL. The key and userName along with shortURL must be passed. REMEMBER, get the API key from /registerNewKey/

        queryParam :

        1) "shortURL"

        2) "userName"

        3) "key"

    example - http://localhost:5899/getVisits?shortURL=http://mydomain.com/MQ==&userName=gg&key=Z2cyMDE5LTA5LTEzVDA3OjUxOjA1Wg==

    (prints count)

## Built With

* [gopkg](https://gopkg.in/redis.v4) - The Redis driver is used
* [stupidbodo](https://gist.github.com/stupidbodo/601b68bfef3449d1b8d9) - AES encryption code is used
 


## Authors

* **Abhishek Bodhekar** - (https://github.com/abhishekbodhekar)



