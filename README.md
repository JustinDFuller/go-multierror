# go-multierror

Package multierror implements an `Append` function that adds two Go errors together.

```
go get github.com/justindfuller/go-multierror
```

## Example

```go
func main() {
  err1 := errors.New("my first error")
  err2 := errors.New("my second error")
  err := multierror.Append(err1, err2)
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

They also support common Go error methods.

* errors.Is
* errors.As

