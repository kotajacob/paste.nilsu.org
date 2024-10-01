// License: AGPL-3.0-only
// (c) 2024 Dakota Walsh <kota@nilsu.org>
package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"
	"time"
)

type application struct {
	Storage string
}

func main() {
	addr := flag.String("addr", ":2016", "HTTP Network Address")
	storage := flag.String("storage", "storage", "Path to storage location")
	flag.Parse()

	err := os.MkdirAll(*storage, 0o755)
	if err != nil {
		log.Fatalln(err)
	}
	app := application{
		Storage: *storage,
	}

	err = app.serve(*addr)
	if err != nil {
		log.Fatalln(err)
	}
}

func (app *application) serve(addr string) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Handle shutdown signals gracefully.
	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		log.Println("shutting down server:", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	log.Println("listening on", addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Println("stopped server")
	return nil
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", app.upload)
	return app.recoverPanic(app.logRequest(app.secureHeaders(mux)))
}

func (app *application) upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024 * 1024 * 5) // Ram cap, not total.
	if err != nil {
		fmt.Println(err)
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	var uploaded string
	file, fileHeader, err := r.FormFile("file")
	if err == nil {
		defer file.Close()

		if fileHeader.Size > (1024 * 1024 * 50) {
			http.Error(
				w,
				"File must be under 50MB",
				http.StatusBadRequest,
			)
			return
		}

		uploaded, err = app.store(file, fileHeader.Filename)
		if err != nil {
			trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
			_ = log.Output(2, trace) // Ignore failed error logging.
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}
	} else if !errors.Is(err, http.ErrMissingFile) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, uploaded, http.StatusSeeOther)
}

func (app *application) store(src io.Reader, name string) (string, error) {
	// Get extension of the uploaded file.
	ext := filepath.Ext(name)

	// Store file in a buffer.
	contents, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}
	buf := bytes.NewReader(contents)

	// Calculate md5sum for the uploaded file.
	h := md5.New()
	_, err = io.Copy(h, buf)
	if err != nil {
		return "", err
	}
	buf.Reset(contents)

	// Store uploaded file.
	f, err := os.Create(filepath.Join(
		app.Storage,
		fmt.Sprintf("%x%v", h.Sum(nil), ext),
	))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(f, buf)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/%x%v", h.Sum(nil), ext), f.Close()
}

// secureHeaders is a middleware which adds strict CSP and other headers.
func (app *application) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'none'; script-src 'none'; style-src 'none"+
				"'; img-src 'self' https: data:",
		)
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

// logRequest is a middleware that prints each request to the info log.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addr := r.RemoteAddr

		// Use correct address if behind proxy.
		if proxyAddr := r.Header.Get("X-Forwarded-For"); proxyAddr != "" {
			addr = proxyAddr
		}

		log.Printf(
			"%s - %s %s %s",
			addr,
			r.Proto,
			r.Method,
			r.URL.RequestURI(),
		)
		next.ServeHTTP(w, r)
	})
}

// recoverPanic is a middleware which recovers from a panic and logs the error.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				trace := fmt.Sprintf("%v\n%s", err, debug.Stack())
				_ = log.Output(2, trace) // Ignore failed error logging.
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
