# multiproc

A simple multiplexer for messages coming in on a channel that shuts down on error.

```go
err := multiproc.ProcessConcurrent(
    messages,
    func(m interface{}) error {
        switch message := m.(type) {
        case string:
            fmt.Println(message)
        default
            return fmt.Errorf("No strings here")
        }
    },
    done
)
if err != int {
    log.Fatalf(err)
}
```
