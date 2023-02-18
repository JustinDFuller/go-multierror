# go-multierror

[![Go Reference](https://pkg.go.dev/badge/github.com/justindfuller/go-multierror.svg)](https://pkg.go.dev/github.com/justindfuller/go-multierror)
[![Go Report Card](https://goreportcard.com/badge/github.com/justindfuller/go-multierror)](https://goreportcard.com/report/github.com/justindfuller/go-multierror)
[![License](https://img.shields.io/github/license/golangci/golangci-lint)](/LICENSE)

Package multierror implements a `Join` function that combines two or more Go errors.

```
go get github.com/justindfuller/go-multierror
```

## Examples

Joining two errors:

```go
func main() {
  err1 := errors.New("my first error")
  err2 := errors.New("my second error")

  err := multierror.Join(err1, err2)

  fmt.Println(err)
  // Found 2 errors:
  //  my first error
  //  my second error
}
```

Joining nil errors will result in `nil` so that you don't have to do extra `nil` checks before joining:

```go
func main() {
	err := multierror.Join(nil, nil, nil)
	fmt.Println(err)
	// <nil>
}
```

## Supported Interfaces

The resulting errors support many common Go interfaces.

* error
* Stringer
* Marshaler
* GoStringer
* GobEncoder
* BinaryMarshaler
* TextMarshaler

```go
func main() {
	err1 := errors.New("something bad happened")
	err2 := errors.New("something is broken")

	err := multierror.Join(err1, err2)
	b, _ := json.Marshal(err)

	fmt.Println(string(b))
	// output: "something bad happened, something is broken"
}
```

They also support common Go error methods.

* errors.Is
* errors.As
* errors.Unwrap

errors.Is:

```go
func main() {
	err1 := errors.New("something bad happened")
	err2 := errors.New("something is broken")
	err3 := errors.New("something is REALLY broken")

	err := multierror.Join(err1, err2)
	err = multierror.Join(err, err3)
	fmt.Println(errors.Is(err, err1))
	// output: true
}
```

errors.As:

```go
func main() {
  _, err := os.Open("non-existing")
	if err == nil {
		fmt.Println("No error")
	}

	err = multierror.Join(err, errSentinelOne)

	var pathError *fs.PathError
	fmt.Println(errors.As(err, &pathError))
	// output: true
}
```

## Why?

I've been unhappy with existing `go-multierror` implementations.
There are three that I am aware of:

1. The standard library [errors.Join](https://pkg.go.dev/errors#Join)
2. Hashicorp's [go-multierror](https://github.com/hashicorp/go-multierror)
3. Uber's [go-multierr](https://github.com/uber-go/multierr)

These libraries have the following problems (in no particular order, not all problems apply to all of them):

* They do not implement common interfaces such as Marshaler, so they don't work with JSON output. This applies to other interfaces and encoders as well.
* They all have different interfaces and methods.
* They expose their underlying error type.

This `go-multierror` solves these problems by:

* Implementing common interfaces (listed above).
* Aligning the interface with the Go standard library.
* Hiding the underlying error type.

