package main

import (
	"context"
	"errors"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/joseph0x45/sad"
	"github.com/joseph0x45/taskeroni/cli"
	"github.com/joseph0x45/taskeroni/db"
	"github.com/joseph0x45/taskeroni/tui"
)

var version = "dev"

const (
	host = "localhost"
	port = "23234"
)

func main() {
	dbPath := setup()
	conn := db.Connect(sad.DBConnectionOptions{
		EnableForeignKeys: true,
		DatabasePath:      dbPath,
	})
	cliFlag := flag.Bool("cli", false, "Use the CLI")
	flag.Parse()
	if *cliFlag {
		opts := &cli.CLIOptions{}
		cli.DispatchCommands(opts, conn)
		return
	}
	hostKeyPath := ".ssh/id_ed25519"
	if version != "dev" {
		hostKeyPath = ""
	}
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(hostKeyPath),
		wish.WithMiddleware(
			bubbletea.Middleware(func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
				renderer := bubbletea.MakeRenderer(sess)
				m := tui.InitApp(renderer, conn)
				return m, []tea.ProgramOption{tea.WithAltScreen()}
			}),
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
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}
