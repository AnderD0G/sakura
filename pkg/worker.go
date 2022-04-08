package pkg

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

type Finish struct {
	IsDone bool
	Err    error
}

type Worker interface {
	Work(ctx context.Context, finishChan chan<- Finish)
}

func Run(timeout time.Duration, ctx context.Context, workers ...Worker) error {
	cancel, cancelFunc := context.WithTimeout(ctx, timeout)
	finishSignal := make(chan Finish)
	errs := make([]error, 0)
	defer cancelFunc()
	for _, worker := range workers {
		go worker.Work(cancel, finishSignal)
	}
	for i := 0; i < len(workers); i++ {
		if f := <-finishSignal; !f.IsDone {
			errs = append(errs, f.Err)
			if len(errs) == 1 {
				cancelFunc()
			}
		}
	}
	close(finishSignal)
	if len(errs) != 0 {
		return errs[0]
	}
	return nil
}

func SafeSend(ch chan<- Finish, value Finish) {
	defer func() {
		if recover() != nil {

		}
	}()
	ch <- value // panic if ch is closed
}

func Watcher(ctx context.Context, ch chan<- Finish) {
	go func() {
	loop:
		for {
			select {
			//无论主动还是被动推出。总会
			case <-ctx.Done():
				SafeSend(ch, Finish{
					IsDone: false,
					Err:    errors.New(""),
				})
				break loop
			}
		}
	}()
}
