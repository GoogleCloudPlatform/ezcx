package ezcx

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// func TestServer(t *testing.T) {
// 	parent := context.Background()
// 	lg := log.Default()

// 	server := NewServer(parent, ":8082", lg, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
// 	mux, err := server.ServeMux()
// 	if err != nil {
// 		t.Log(err)
// 	}
// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintln(w, "Hello World!")
// 	})
// 	server.ListenAndServe(parent)
// }

func TestCxHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(sample))
	w := httptest.NewRecorder()
	handler := HandlerFunc(CxHandler)
	handler.ServeHTTP(w, req)
	resp := w.Result()
	io.Copy(os.Stdout, resp.Body)
	t.Log(resp)

}

func CxHandler(res *WebhookResponse, req *WebhookRequest) error {
	res.AddTextResponse("With much technolove from Yvan J. Aquino - I wrote this!")
	return nil
}
