package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Реализовать функцию, которая выполняет поиск query во всех переданных SearchFunc
// Когда получаем первый успешный результат - отдаем его сразу, не дожидаясь результата других SearchFunc
// Если все SearchFunc отработали с ошибкой - отдаем последнюю полученную ошибку

type Result struct {
	Data string
}

type SearchFunc func(ctx context.Context, query string) (*Result, error)

func main() {
	sfs := []SearchFunc{foo1, foo2, foo3}

	ctx := context.Background()
	query := "tmp"

	r, err := MultiSearch(ctx, query, sfs)
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	} else {
		fmt.Println("RESULT:", r)
	}

}

func MultiSearch(ctx context.Context, query string, sfs []SearchFunc) (*Result, error) {
	wg := &sync.WaitGroup{}
	errPoint := atomic.Pointer[error]{}
	resCh := make(chan *Result, len(sfs))

	ctx, cancel := context.WithCancel(ctx)

	for _, f := range sfs {
		wg.Add(1)

		go func() {
			defer wg.Done()

			r, err := f(ctx, query)
			if err != nil {
				errPoint.Store(&err)
			}

			resCh <- r
			cancel()

		}()
	}

	// ждать все горутины для получения последней ошибки, но при этом нужно получить первый результат
	wg.Wait()
	close(resCh)

	// если error нет то читаем первый результат из канала
	if err := errPoint.Load(); err != nil {
		return nil, *err
	}

	r := <-resCh

	return r, nil
}

func foo1(ctx context.Context, query string) (*Result, error) {
	// n := rand.Intn(5)
	time.Sleep(time.Duration(2) * time.Second)
	return &Result{
		Data: "first",
	}, nil
}

func foo2(ctx context.Context, query string) (*Result, error) {
	// n := rand.Intn(5)
	time.Sleep(time.Duration(3) * time.Second)
	return &Result{
		Data: "second",
	}, nil
}

// func foo3(ctx context.Context, query string) (*Result, error) {
// 	// n := rand.Intn(10)
// 	time.Sleep(time.Duration(1) * time.Second)
// 	return nil, errors.New("foo3 error")
// }

func foo3(ctx context.Context, query string) (*Result, error) {
	// n := rand.Intn(5)
	time.Sleep(time.Duration(1) * time.Second)
	return &Result{
		Data: "third",
	}, nil
}
