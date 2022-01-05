# Iaphub-go

`iaphub-go` is an unofficial Go client for [iaphub.com REST API](https://www.iaphub.com/docs/api/).

[![Go Report Card](https://goreportcard.com/badge/github.com/n10ty/iaphub-go)](https://goreportcard.com/report/github.com/n10ty/iaphub-go)
[![GoDoc](https://godoc.org/https://godoc.org/github.com/n10ty/iaphub-go?status.svg)](https://godoc.org/github.com/n10ty/iaphub-go)

## Installation

`go get github.com/n10ty/iaphub-go`

## Usage
```go
package main

import (
	"fmt"
	"github.com/n10ty/iaphub-go"
)

func main() {
    
        iaphubSecret := <Your secret>
        iaphubAppId := <Your app id>

	c, err := iaphub.NewClient(iaphubSecret, iaphubAppId)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	userRequest := iaphub.GetUserRequest{
		UserId:   userid,
		Platform: iaphub.PlatformIOS,
	}
	user, err := c.GetUser(getuser)
    
        if err != nil {
            fmt.Println(err)
            return
	}
	fmt.Println(user)
}
```

### Custom environment

```go
c, err := iaphub.NewClient(iaphubSecret, iaphubAppId, iaphub.UseEnv("sandbox"))
```


### Supported methods

* Get user
* Get user migrate
* Post user
* Post user receipt
* Get purchase
* Get purchases
* Get subscription
* Get receipt

## License

[MIT License](LICENSE)
