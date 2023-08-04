package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Jhon-2801/course_user/internal/user"
	boostrap "github.com/Jhon-2801/course_user/pkg/bootstrap"
	"github.com/Jhon-2801/course_user/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db, err := boostrap.DBConnection()
	if err != nil {
		log.Fatal(err)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		log.Fatal("paginator limit default is required")
	}

	ctx := context.Background()
	userRepo := user.NewRepo(db)
	userSrv := user.NewService(userRepo)

	h := handler.NewUserHTTPServer(ctx, user.MakeEndpoints(userSrv, user.Config{LimPageDef: pagLimDef}))

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	errCh := make(chan error)
	go func() {
		fmt.Println("listen in ", address)
		errCh <- srv.ListenAndServe()
	}()

	err = <-errCh
	if err != nil {
		log.Fatal(err)
	}
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, HEAD, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
