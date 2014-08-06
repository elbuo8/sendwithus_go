# SendWithUs-Go

This is a simple package to interface with [SendWithUs](https://sendwithus.com) using Golang.

## Installation

```bash
$ go get github.com/elbuo8/sendwithus_go
```

## Example

This is a brief example on how to send 1 email. You can find more examples by looking at [the test cases](https://github.com/elbuo8/sendwithus_go/blob/master/swu_test.go).

```go
package main

import (
  "github.com/elbuo8/sendwithus_go"
  "fmt"
)

func main() {
	api := New("SWU_KEY")
	email := &SWUEmail{
		ID: "EMAIL_TEMPLATE_ID",
		Recipient: &SWURecipient{
			Address: "example@email.com",
		},
		EmailData: make(map[string]string),
	}
	err := api.Send(email)
	if err != nil {
      fmt.Println(err)
	}
}

```

## [Documentation (GoDoc)](https://github.com/elbuo8/sendwithus_go/blob/master/swu_test.go)

## MIT License

Enjoy! Feel free to send pull requests or submit issues :)
