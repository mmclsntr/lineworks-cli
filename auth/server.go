package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type CallbackHandler struct {
	Context       context.Context
	ContextCancel context.CancelFunc
	CallbackFunc  func(code string, state string)
}

func (handler *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryParts, _ := url.ParseQuery(r.URL.RawQuery)

	if _, isExists := queryParts["code"]; !isExists {
		fmt.Println("code does not exist.")
		fmt.Fprintf(w, "invalid request")
		return
	}
	code := queryParts["code"][0]
	state := queryParts["state"][0]

	// callback
	handler.CallbackFunc(code, state)

	defer handler.ContextCancel()

	// show succes page
	msg := "<p><strong>Success!</strong></p>"
	msg = msg + "<p>You are authenticated and can now return to the CLI.</p>"
	fmt.Fprintf(w, msg)

}

// Start callback local server
func StartCallbackServer(ctx context.Context, addr string, port string, path string, timeoutSec int, callback func(code string, state string)) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)

	handler := &CallbackHandler{
		Context:       ctx,
		ContextCancel: cancel,
		CallbackFunc:  callback,
	}
	http.Handle(path, handler)

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%s", addr, port),
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	defer cancel()

	select {
	case <-ctx.Done():
		// shutdown
		//if err := srv.Shutdown(ctx); err != nil {
		//    log.Printf("%+v", err)
		//}
		if err := ctx.Err(); errors.Is(err, context.Canceled) {
			// キャンセルされていた場合
			fmt.Println("Done")
		} else if errors.Is(err, context.DeadlineExceeded) {
			// タイムアウトだった場合
			fmt.Println("Time out")
		}
	}
}
