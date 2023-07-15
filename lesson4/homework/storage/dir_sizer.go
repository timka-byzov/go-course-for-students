package storage

import (
	"context"
	"sync"
	// "sync"
)

const (
	maxWorkersCount = 3
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
	workerPool      chan struct{}
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{0, make(chan struct{}, maxWorkersCount)}
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	ch := make(chan File, 1000)

	totalSize := int64(0)
	count := int64(0)

	go func() {
		defer close(ch)
		a.RecursiveSize(ctx, d, ch)
	}()

	for f := range ch {
		select {
		case <-ctx.Done():
			return Result{}, ctx.Err()
		default:
			count++
			size, err := f.Stat(ctx)
			if err != nil {
				return Result{}, err
			}
			totalSize += size
		}
	}

	return Result{Size: totalSize, Count: count}, nil
}

func (a *sizer) RecursiveSize(ctx context.Context, d Dir, ch chan File) {
	dirs, files, err := d.Ls(ctx)
	if err != nil {
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
		a.workerPool <- struct{}{} // Запрашиваем доступ к горутине из пула
		go func(dir Dir) {
			defer func() {
				wg.Done()
				<-a.workerPool // Освобождаем горутину в пуле
			}()
			a.RecursiveSize(ctx, dir, ch)
		}(d)
	}
	wg.Wait()
}
