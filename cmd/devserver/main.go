// Command devserver は、ドメインロジックをブラウザや curl から手で叩いて確認する
// ためのローカル専用 HTTP サーバーである。Lambda にはデプロイしない。
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultPort = "8080"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", hello)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = defaultPort
	}

	srv := &http.Server{
		// コンテナ内で 127.0.0.1 に閉じるとポートマッピングを通ってもホストから
		// 到達できないため、ホスト部を空にして 0.0.0.0 で待ち受ける。
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("devserver listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}

func hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}
