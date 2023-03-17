[![Go](https://github.com/lab210-dev/service/actions/workflows/go.yml/badge.svg)](https://github.com/lab210-dev/service/actions/workflows/go.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lab210-dev/service)
[![Go Report Card](https://goreportcard.com/badge/github.com/lab210-dev/service)](https://goreportcard.com/report/github.com/lab210-dev/service)
[![codecov](https://codecov.io/gh/lab210-dev/service/branch/main/graph/badge.svg?token=3JRL5ZLSIH)](https://codecov.io/gh/lab210-dev/service)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/lab210-dev/service/blob/main/LICENSE)
[![Github tag](https://badgen.net/github/release/lab210-dev/service)](https://github.com/lab210-dev/service/releases)

# Overview
This package provides a simple way to manage services in a Go application. It allows you to register and retrieve services using Go interfaces.

It simplifies the process of writing unit tests by providing a simple and easy-to-use interface for managing dependencies between different components of your application. It allows you to register services and callbacks, and retrieve them whenever needed, making it easy to test your code in isolation. This results in more maintainable and reliable tests, as well as a faster development process.

## Installation
To install this package, use the `go get` command:

```bash
go get github.com/lab210-dev/service
```

## 📚 Usage
To use this package, you must first register your service using the Register function :

```go
package example

import "github.com/lab210-dev/servicekit"

type MyService interface {
	DoSomething()
}

type myServiceImpl struct{}

func (s *myServiceImpl) DoSomething() {
	// Do something here...
}

func init() {
	servicekit.Register[MyService](new(myServiceImpl))
}
```

You can now retrieve your service using the `Get` function :

```go
package example

import "github.com/lab210-dev/servicekit"

type MyService interface {
	DoSomething()
}

func main() {
	servicekit.Get[MyService]().DoSomething()
}
```

You can also use the GetAfterRegister function to execute a callback once the service has been registered :

```go 
package example

import "github.com/lab210-dev/servicekit"

type MyService interface {
	DoSomething()
}

func main() {
    servicekit.GetAfterRegister[MyService](func(s MyService) {
        s.DoSomething()
    })
}
```

To reset all registered services, use the Reset function :
```go
package example

import "github.com/lab210-dev/servicekit"

func main() {
	servicekit.Reset()
}
```

### Notes
- Services must be registered using interfaces, not concrete types.
- If you try to retrieve a service that has not been registered, an error will panic.


## ⚡️ Benchmark

```text
goos: darwin
goarch: arm64
pkg: github.com/lab210-dev/servicekit
BenchmarkServiceGet
BenchmarkServiceGet-10                 	 7306159	       162.2 ns/op
BenchmarkServiceGetAfterRegister
BenchmarkServiceGetAfterRegister-10    	 5068532	       232.1 ns/op
```

## 🤝 Contributions
Contributors to the package are encouraged to help improve the code.

If you have any questions or comments, don't hesitate to contact me.

I hope this package is helpful to you !