# errx
golang errors with stack traces, based off [pkg/errors](https://github.com/pkg/errors) modified for use within my own projects.


errx implements error utilities which make adding context to errors more developer friendly while also preserving the original error.

## Create an error

```go
// using errx.New
if ok := someOperation(); !ok {
	return errx.New("someOperation failed")
}

// using errx.Errorf
if ok := someOperation(); !ok {
	return errx.Errorf("%s failed", "someOperation()")
}
```


## Add context to an error

```go
// using errx.Wrap                                                  
_, err := someOperationThatMayErr()                                 
if err != nil {                                                     
	return errx.Wrap(err, "someOperationThatMayErr() failed")       
}                                                                   
                                                                    
// using errx.Wrapf                                                 
_, err := someOperationThatMayErr()                                 
if err != nil {                                                     
	return errx.Wrapf(err, "%s failed", "someOperationThatMayErr()")
}                                                                   
```

## Print with stack trace

```go
err := errx.Wrap(errx.Wrap(errors.New("inner"), "middle"), "outer")
log.Printf("%v", err) // or use %s
```

yields

```
outer: middle: inner
  at main.main(examples/main.go:13)
  at runtime.main(runtime/proc.go:207)
  at runtime.goexit(runtime/asm_amd64.s:2362)
```

## Print without stack trace

```go
err := errx.Wrap(errx.Wrap(errors.New("inner"), "middle"), "outer")
log.Printf("%-v", err) // or use %-s
```

yields

```
outer: middle: inner
```


## Print top most error

```go
err := errx.Wrap(errx.Wrap(errors.New("inner"), "middle"), "outer")
if e, ok := err.(*errx.Error); ok {
	log.Print(e.Message)
}
```

yields

```
outer
```















