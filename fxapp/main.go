package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/fx"
)

type Sender interface {
	Send()
}

type EmailSender struct{}

func (e *EmailSender) Send(to, body string) {
	fmt.Printf("📧 Отправлено письмо на %s: %s\n", to, body)
}

type UserHandler struct {
	sender *EmailSender
}

func NewEmailSender() *EmailSender {
	return &EmailSender{}
}

func NewUserHandler(sender *EmailSender) *UserHandler {
	return &UserHandler{sender: sender}
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	uh.sender.Send("user@example.com", "Добро пожаловать!")
	fmt.Fprintln(w, "Пользователь зарегистрирован.")
}

func (uh *UserHandler) Akn(w http.ResponseWriter, r *http.Request) {
	uh.sender.Send("user@example.com", "Вы замечены!")
	fmt.Fprintln(w, "Пользователь зашел в на сайт.")
}

func NewMux(uh *UserHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", uh.Register)
	mux.HandleFunc("/", uh.Akn)
	return mux
}

func NewHTTPServer(mux *http.ServeMux, lc fx.Lifecycle) *http.Server {
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("🚀 Сервис стартует на :8080")
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("🛑 Завершается...")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

func main() {
	app := fx.New(
		fx.Provide(
			NewMux,
			NewHTTPServer,
			NewEmailSender,
			NewUserHandler,
		),
		fx.Invoke(func(*http.Server) {}),
	)

	app.Run()
}
