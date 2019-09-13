# Project Title

This is an implementation of URL shortning

# Description

This is a package contains implementation for shortning the URLs.
The golang is used on application layer. Redis is used for database.
A docker-compose is provided as well.


It inculdes every data stored in redis completely incrypted with AES ecryption.
Also, It stores the number of visits made to any short link.
An API key is required to get the # visits.
It provides an API to create a API key. Access to generate API key is currently unrestricted.


### Prerequisites

Docker with compose is needed to use with containerization.
Otherwise, only redis server is needed to be installed on host as compiled Go package itself is an executable.

## Implementation

HOW IT WORKS - 

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



## Running the tests

Explain how to run the automated tests for this system

### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```



## Deployment

Add additional notes about how to deploy this on a live system

## Built With

* [Dropwizard](http://www.dropwizard.io/1.0.2/docs/) - The web framework used
* [Maven](https://maven.apache.org/) - Dependency Management
* [ROME](https://rometools.github.io/rome/) - Used to generate RSS Feeds



## Authors

* **Abhishek Bodhekar** - (https://github.com/abhishekbodhekar)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc
