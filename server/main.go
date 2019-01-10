package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kechako/sigctx"
)

func main() {
	var addr string
	var dir string
	flag.StringVar(&addr, "addr", ":5555", "server address:port")
	flag.StringVar(&dir, "dir", ".", "static contents directory")
	flag.Parse()

	log.Print("Start HTTP server...")

	ctx, cancel := sigctx.WithCancelBySignal(context.Background(), os.Interrupt)

	srv := &http.Server{
		Addr:    addr,
		Handler: http.FileServer(http.Dir(dir)),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		<-ctx.Done()

		log.Print("Shutdown HTTP server...")
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctxShutdown); err != nil {
			log.Print("Failed to gracefully shutdown HTTP server : ", err)
			// force close
			if err := srv.Close(); err != nil {
				log.Print("[Error] Failed to close HTTP server : ", err)
			}
		} else {
			log.Print("Done.")
		}
		close(idleConnsClosed)
	}()

	log.Printf("listening on %q...", addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Print("[Error] ", err)
	}

	cancel()
	<-idleConnsClosed
}
