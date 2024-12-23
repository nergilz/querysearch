package main

import (
	"context"
	"sync"
	"sync/atomic"
)

// Реализовать функцию, которая выполняет поиск query во всех переданных SearchFunc
// Когда получаем первый успешный результат - отдаем его сразу, не дожидаясь результата других SearchFunc
// Если все SearchFunc отработали с ошибкой - отдаем последнюю полученную ошибку

type Result struct{}

type SearchFunc func(ctx context.Context, query string) (*Result, error)

func main() {
	// TODO init search

}

func MultiSearch(ctx context.Context, query string, sfs []SearchFunc) (*Result, error) {
	wg := &sync.WaitGroup{}
	errPoint := atomic.Pointer[error]{}
	resCh := make(chan *Result)

	ctx, cancel := context.WithCancel(ctx)

	for _, f := range sfs {
		wg.Add(1)

		go func() {
			defer wg.Done()

			r, err := f(ctx, query)
			if err != nil {
				errPoint.Store(&err)
			} else {
				resCh <- r
				cancel()
			}
		}()
	}

	// ждать все горутины для получения последней ошибки, но при этом нужно получить первый результат
	wg.Wait()
	close(resCh)

	if err := *errPoint.Load(); err != nil {
		return nil, err
	}

	// канал буфиризированный - блокирует, и нам нужно только первое значение
	r := <-resCh

	return r, nil
}
