<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# watcher

```go
import "github.com/leads-su/consul/watcher"
```

## Index

- [type Watcher](<#type-watcher>)
  - [func (watcher *Watcher) Start()](<#func-watcher-start>)
  - [func (watcher *Watcher) Stop() error](<#func-watcher-stop>)


## type Watcher

```go
type Watcher struct {
    sync.Mutex
    Client            *consulAPI.Client
    Prefix            string
    UpdateChannel     chan<- consulAPI.KVPairs
    ErrorChannel      chan<- error
    QuiescencePeriod  time.Duration
    QuiescenceTimeout time.Duration
    // contains filtered or unexported fields
}
```

### func \(\*Watcher\) Start

```go
func (watcher *Watcher) Start()
```

### func \(\*Watcher\) Stop

```go
func (watcher *Watcher) Stop() error
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
