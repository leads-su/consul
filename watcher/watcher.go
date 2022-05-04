package watcher

import (
	"errors"
	"github.com/cenkalti/backoff"
	consulAPI "github.com/hashicorp/consul/api"
	"sync"
	"time"
)

type Watcher struct {
	sync.Mutex
	Client				*consulAPI.Client
	Prefix				string
	UpdateChannel		chan <-consulAPI.KVPairs
	ErrorChannel		chan <-error
	QuiescencePeriod	time.Duration
	QuiescenceTimeout	time.Duration

	quitChannel			chan <-struct{}
	doneChannel			<-chan struct{}
}

func (watcher *Watcher) Start () {
	watcher.Lock()

	if watcher.doneChannel != nil {
		watcher.Unlock()
		return
	}

	quitChannel := make(chan struct{})
	doneChannel := make(chan struct{})
	watcher.quitChannel = quitChannel
	watcher.doneChannel = doneChannel
	watcher.Unlock()

	defer func() {
		watcher.Lock()
		defer watcher.Unlock()
		close(doneChannel)
		watcher.doneChannel = nil
	}()

	errorChannel, ok := watcher.errorChannel()

	if !ok {
		defer close(errorChannel)
	}

	if watcher.Prefix == "" {
		errorChannel <- errors.New("prefix cannot be empty")
		return
	}

	if watcher.Prefix[len(watcher.Prefix) - 1] != '/' {
		watcher.Prefix  += "/"
	}

	qscPeriod := watcher.QuiescencePeriod
	qscTimeout := watcher.QuiescenceTimeout

	if qscPeriod == 0 {
		qscPeriod = 500 * time.Millisecond
	}
	if qscTimeout == 0 {
		qscTimeout = 5 * time.Second
	}

	pairsChannel := make(chan consulAPI.KVPairs)

	go func() {
		var waitIndex uint64
		for {
			var pairs consulAPI.KVPairs
			var meta *consulAPI.QueryMeta

			queryOptions := &consulAPI.QueryOptions{
				WaitIndex:	waitIndex,
				WaitTime:	30 * time.Minute,
			}

			err := backoff.Retry(func() error {
				select {
				case <-quitChannel:
					return nil
				default:
				}

				var err error
				pairs, meta, err = watcher.Client.KV().List(watcher.Prefix, queryOptions)

				select {
				case <-quitChannel:
					return nil
				default:
				}

				if err != nil {
					errorChannel <- err
				}
				return err
			}, watcher.backOff())

			if err != nil {
				continue
			}

			select {
			case <-quitChannel:
				return
			default:
			}

			if meta.LastIndex == waitIndex {
				continue
			}
			waitIndex = meta.LastIndex
			pairsChannel <- pairs
		}
	}()

	init := false
	var pairs consulAPI.KVPairs
	var qscPeriodChannel, qscTimeoutChannel <-chan time.Time

	for {
		select {
		case <-quitChannel:
			return
		case pairs = <-pairsChannel:
			qscPeriodChannel = time.After(qscPeriod)
			if qscTimeoutChannel == nil {
				qscTimeoutChannel = time.After(qscTimeout)
			}
			if init {
				continue
			}
			init = true
		case <-qscPeriodChannel:
		case <-qscTimeoutChannel:
		}

		qscPeriodChannel = nil
		qscTimeoutChannel = nil

		watcher.UpdateChannel <-pairs
	}
}

func (watcher *Watcher) Stop() error {
	watcher.Lock()

	if watcher.doneChannel == nil {
		watcher.Unlock()
		return nil
	}

	if watcher.quitChannel != nil {
		close(watcher.quitChannel)
		watcher.quitChannel = nil
	}

	doneChannel := watcher.doneChannel
	watcher.Unlock()
	<-doneChannel
	return nil
}

func (watcher *Watcher) errorChannel() (chan <-error, bool) {
	errorChannel := watcher.ErrorChannel
	ok := true

	if errorChannel == nil {
		ok = false
		channel := make(chan error)
		errorChannel = channel
		go func() {
			for range channel {}
		}()
	}
	return errorChannel, ok
}

func (watcher *Watcher) backOff() backoff.BackOff {
	result := backoff.NewExponentialBackOff()
	result.InitialInterval = 1 * time.Second
	result.MaxInterval = 10 * time.Second
	result.MaxElapsedTime = 0
	return result
}