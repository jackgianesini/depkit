[![Go](https://github.com/kitstack/depkit/actions/workflows/coverage.yml/badge.svg)](https://github.com/kitstack/depkit/actions/workflows/coverage.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kitstack/depkit)
[![Go Report Card](https://goreportcard.com/badge/github.com/kitstack/depkit)](https://goreportcard.com/report/github.com/kitstack/depkit)
[![codecov](https://codecov.io/gh/kitstack/depkit/branch/main/graph/badge.svg?token=3JRL5ZLSIH)](https://codecov.io/gh/kitstack/depkit)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/kitstack/depkit/blob/main/LICENSE)
[![Github tag](https://badgen.net/github/release/kitstack/depkit)](https://github.com/kitstack/depkit/releases)

# Overview
This package provides a simple way to manage dependencies in a Go application. It allows you to register and retrieve dependencies using Go interfaces.

It simplifies the process of writing unit tests by providing a simple and easy-to-use interface for managing dependencies between different components of your application. It allows you to register dependencies and callbacks, and retrieve them whenever needed, making it easy to test your code in isolation. This results in more maintainable and reliable tests, as well as a faster development process.

## Installation
To install this package, use the `go get` command:

```bash
go get github.com/kitstack/depkit
```

## 📚 Usage
To use this package, you must first register your dependency using the Register function :

```go
package example

import "github.com/kitstack/depkit"

type MyService interface {
	DoSomething()
}

type myServiceImpl struct{}

func (s *myServiceImpl) DoSomething() {
	// Do something here...
}

func init() {
	depkit.Register[MyService](new(myServiceImpl))
}
```

You can now retrieve your service using the `Get` function :

```go
package example

import "github.com/kitstack/depkit"

type MyService interface {
	DoSomething()
}

func main() {
	depkit.Get[MyService]().DoSomething()
}
```

You can also use the GetAfterRegister function to execute a callback once the service has been registered :

```go 
package example

import "github.com/kitstack/depkit"

type MyService interface {
	DoSomething()
}

func main() {
    depkit.GetAfterRegister[MyService](func(s MyService) {
        s.DoSomething()
    })
}
```

To reset all registered dependencies, use the Reset function :
```go
package example

import "github.com/kitstack/depkit"

func main() {
	depkit.Reset()
}
```

### Notes
- Services must be registered using interfaces or func, not concrete types.
- If you try to retrieve a dependency that has not been registered, an error will panic.

## ⚡️ Benchmark

```text
goos: darwin
goarch: arm64
pkg: github.com/kitstack/depkit
BenchmarkGet
BenchmarkGet-10                 	 6689894	       167.8 ns/op
BenchmarkGetAfterRegister
BenchmarkGetAfterRegister-10    	 5415369	       223.7 ns/op
```

## 🤝 Contributions
Contributors to the package are encouraged to help improve the code.

If you have any questions or comments, don't hesitate to contact me.

I hope this package is helpful to you !