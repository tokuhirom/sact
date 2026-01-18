package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tokuhirom/sact/internal"
)

func buildLogger(logPath string) (*slog.Logger, error) {
	if logPath == "" {
		logPath = os.DevNull // not portable
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return logger, nil
}

func initLogger(logPath string) error {
	logger, err := buildLogger(logPath)
	if err != nil {
		return fmt.Errorf("failed to build logger: %w", err)
	}

	slog.SetDefault(logger)
	return nil
}

func main() {
	logPath := flag.String("log", "", "Path to log file")
	flag.Parse()

	if err := initLogger(*logPath); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		if err != nil {
			slog.Error("Failed to write to stderr", slog.Any("error", err))
		}
		os.Exit(1)
	}

	slog.Info("Starting sact", slog.String("log_path", *logPath))

	opts, defaultZone, err := internal.LoadProfileAndZone()
	if err != nil {
		slog.Error("Failed to load profile", slog.Any("error", err))
		_, err := fmt.Fprintf(os.Stderr, "Error loading profile: %v\n", err)
		if err != nil {
			slog.Error("Failed to write to stderr", slog.Any("error", err))
		}
		os.Exit(1)
	}
	slog.Info("Profile loaded", slog.String("default_zone", defaultZone))

	client, err := internal.NewSakuraClient(opts, defaultZone)
	if err != nil {
		slog.Error("Failed to create client", slog.Any("error", err))
		_, err := fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		if err != nil {
			slog.Error("Failed to write to stderr", slog.Any("error", err))
		}
		os.Exit(1)
	}
	slog.Info("Client created", slog.String("zone", defaultZone))

	p := tea.NewProgram(internal.InitialModel(client, defaultZone))
	if _, err := p.Run(); err != nil {
		slog.Error("Program failed", slog.Any("error", err))
		_, err := fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		if err != nil {
			slog.Error("Failed to write to stderr", slog.Any("error", err))
			return
		}
		os.Exit(1)
	}

	slog.Info("Program exited normally")
}
