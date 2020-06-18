# Golang IoC
This is a pure golang library to using inversion of control (ioc) in golang

## Installation
```bash
go get github.com/morilog/gioc
```

## Example
```golang
package main

import "guthub.com/morilog/gioc"
import "log"
import "fmt"

func main() {
    // Bind binds a resolver to abstract type and
    // the resolver called every time you needed to the type
    gioc.Bind(func () (Greeter, err) {
        return &GoodMorning{}
    })

    var g Greeter
    if err := gioc.Make(&g); err != nil {
        log.Fatal(err)
    }

    fmt.Println(g.SayHello()) // prints "good morning"


    gioc.Singleton(func (greeter Greeter, c *SimpleClient) (Mailer, error)) {
        return &simpleMailer(simpleClient: c, g: greeter), nil
    }

    var m Mailer
    if err := gioc.Make(&m); err != nil {
        log.Fatal(err)
    }

    m.Send("receiver@host.com", "don't reply me")
}


type Greeter interface{
    SayHello() string
}

type Morning struct{}

func (g Morning) SayHello() string {
    return "good morning"
}

type Mailer interface {
    Send(to string, msg string) error
}

type SimpleMailer struct{
    g Greeter
    simpleClient *client
}

func (s *SimpleMailer) Send(to string, msg string) error {
    msg = s.g.SayHello() + msg

    s.simpleClient.SendText(to, msg, "info@example.com")
}

```