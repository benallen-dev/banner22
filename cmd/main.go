package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "embed"

	"banner22/internal/charmapp"

	tea "github.com/charmbracelet/bubbletea"
	redis "github.com/redis/go-redis/v9"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "localhost"
	port = "23234"
)

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	//pty, _, ok := s.Pty()
	// if !ok {
	// 	panic("Could not get PTY") // will crash the whole server but it's fine, right? Right?
	// }

	// MakeRenderer makes a renderer for the session client, instead of the system the server is running on
	renderer := bubbletea.MakeRenderer(s)
	textStyle    := renderer.NewStyle().Foreground(lipgloss.Color("252"))
	spinnerStyle := renderer.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle    := renderer.NewStyle().Foreground(lipgloss.Color("241"))
	// txtStyle := renderer.NewStyle().
	//  	Foreground(lipgloss.Color("#DDFFDD")).
	//  	Background(lipgloss.Color("#7D56F4")).
	//  	Padding(4)
//	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

	// bg := "light" // seems odd that this is just a rando string but let's just roll with it for now
	// if renderer.HasDarkBackground() {
	// 	bg = "dark"
	// }

	// m := TermModel{
	// 	term:      pty.Term,
	// 	profile:   renderer.ColorProfile().Name(),
	// 	width:     pty.Window.Width,
	// 	height:    pty.Window.Height,
	// 	bg:        bg,
	// 	txtStyle:  txtStyle,
	// 	quitStyle: quitStyle,
	// }

	m := charmapp.Spinmodel{
		TextStyle: textStyle,
		SpinnerStyle: spinnerStyle,
		HelpStyle: helpStyle,
	}
	m.ResetSpinner()

	return m, []tea.ProgramOption{tea.WithAltScreen()}

}

func main() {
	// Create redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "benisawesome",
		DB: 0,
	})

	// Test redis connection
	pong, err := rdb.Ping(context.Background()).Result()
	log.Info("Redis connection", "pong", pong, "error", err)

	// Literally all boilerplate, the magic lives in the teaHandler
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithPasswordAuth( func(ctx ssh.Context, password string) bool {
			log.Info("Password entered", "password", password)
			// Store the password in redis
			result := rdb.Incr(ctx, password)
			if result.Err() != nil {
				log.Warn("Failed to increment key in Redis", "error", result.Err().Error())
			}

			return false
		}),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)

	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}

}
