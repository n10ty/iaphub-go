# Iaphub-go

`iaphub-go` is an unofficial Go client for [iaphub.com REST API](https://www.iaphub.com/docs/api/).

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
