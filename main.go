/*
* @Author: 李澳华
* @Date:   2021/6/13 22:14 
*/

package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
)


func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	group, errctx := errgroup.WithContext(ctx)
	srv := &http.Server{
		Addr: ":1234",
	}

	//打开Server
	group.Go(func() error {
		return OpenServer(srv)
	})

	//后台挂起关闭的线程
	group.Go(func() error {
		<-errctx.Done()
		return srv.Shutdown(errctx)
	})

	//linux signal处理
	ch := make(chan os.Signal, 1)
	signal.Notify(ch)

	group.Go(func() error {
		for {
			select {
			case <-errctx.Done():
					return errctx.Err()
			case <-ch:
				cancel()
			}
		}
		return nil
	})
	err := group.Wait()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("all done")
}

func OpenServer(srv *http.Server) (err error) {
	http.HandleFunc("/",Response)
	err = srv.ListenAndServe()
	return err
}


func Response(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
