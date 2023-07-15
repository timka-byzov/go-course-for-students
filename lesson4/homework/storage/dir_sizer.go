package storage

import (
	"context"
	"sync"
	// "sync"
)

const (
	maxWorkersCount = 10
)

type Result struct {
	Size  int64
	Count int64
}

type DirSizer interface {
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	maxWorkersCount int
	workerCh        chan struct{}
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{0, make(chan struct{}, maxWorkersCount)}
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	ch := make(chan File, 1000)
	errCh := make(chan error, 1)
	totalSize := int64(0)
	count := int64(0)

	go func() {
		defer close(ch)
		a.RecursiveSize(ctx, d, ch, errCh)
	}()

	for {
		select {
		case <-ctx.Done():
			return Result{}, ctx.Err()
		case err := <-errCh:
			return Result{}, err
		case f, ok := <-ch:
			if !ok {
				return Result{Size: totalSize, Count: count}, nil
			}
			count++
			size, err := f.Stat(ctx)
			if err != nil {
				return Result{}, err
			}
			totalSize += size
		}
	}
}

func (a *sizer) RecursiveSize(ctx context.Context, d Dir, ch chan File, errCh chan error) {
	dirs, files, err := d.Ls(ctx)
	if err != nil {
		errCh <- err
		return
	}

	for _, f := range files {
		select {
		case ch <- f:
		case <-ctx.Done():
			return
		}
	}

	var wg sync.WaitGroup
	for _, d := range dirs {
		wg.Add(1)
		a.workerCh <- struct{}{}
		go func(dir Dir) {
			defer func() {
				wg.Done()
				<-a.workerCh
			}()
			a.RecursiveSize(ctx, dir, ch, errCh)
		}(d)
	}
	wg.Wait()
}
