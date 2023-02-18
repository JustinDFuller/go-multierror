# go-multierror

Package multierror implements a `Join` function that adds two Go errors together.

```
go get github.com/justindfuller/go-multierror
```

## Example

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

The resulting errors support several common Go interfaces.

* Error
* Stringer
* Marshaler
* GoStringer
* GobEncode

They also support common Go error methods.

* errors.Is
* errors.As


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

