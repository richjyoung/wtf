# wtf.IsThis

Golang is great until at runtime you are dealing with an interface for which the underlying implementation is buried in a package somewhere.

Yes it's easier than other languages to click through and find it, or...

```golang
wtf.IsThis(1)                              // int
wtf.IsThis("Hello, World!")                // string
wtf.IsThis(AStruct{})                      // pkg.AStruct
wtf.IsThis(&AStruct{})                     // *pkg.AStruct
wtf.IsThis([]*AStruct{})                   // []*pkg.AStruct
wtf.IsThis(map[int]*AStruct{})             // map[int]*pkg.AStruct
wtf.IsThis(make(<-chan int))               // <-chan int
wtf.IsThis(func(int) error { return nil }) // func (int) error {}

// In cases where the type is nil or contains an empty interface, wtf.NoIdea is returned.

wtf.IsThis(nil)                            // ¯\_(ツ)_/¯
wtf.IsThis(map[string]interface{}{})       // map[string]¯\_(ツ)_/¯

// Provides more detail for errors which may implement Unwrap

e1 := fmt.Errorf("error 1")
e2 := fmt.Errorf("error 2 - %w", e1)
e3 := fmt.Errorf("error 3 - %w", e2)

wtf.IsThisError(e3) // Returns:
// *fmt.wrapError[error 3 - error 2 - error 1]
// └─*fmt.wrapError[error 2 - error 1]
//   └─*errors.errorString[error 1]
```

More examples in the [test](./wtf_test.go).
