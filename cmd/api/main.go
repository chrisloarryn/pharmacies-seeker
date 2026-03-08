package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"

	pharmacyhttp "pharmacies-seeker/internal/pharmacies/adapters/http"
	"pharmacies-seeker/internal/pharmacies/adapters/provider/boostr"
	"pharmacies-seeker/internal/pharmacies/adapters/repository/memory"
	"pharmacies-seeker/internal/pharmacies/app"
	platformconfig "pharmacies-seeker/internal/platform/config"
	platformhttp "pharmacies-seeker/internal/platform/http"
	"pharmacies-seeker/internal/platform/scheduler"
)

type schedulerRunner interface {
	Start(context.Context, string, scheduler.Task)
}

type serverRunner interface {
	Run() error
}

type runtimeDeps struct {
	loadConfig    func(string) (platformconfig.Config, error)
	newRepository func() app.Repository
	newProvider   func(string, string, time.Duration) app.Provider
	newRunner     func(*log.Logger, time.Duration, time.Duration) schedulerRunner
	newServer     func(string, *fiber.App) serverRunner
	withTimeout   func(context.Context, time.Duration) (context.Context, context.CancelFunc)
}

var (
	backgroundContext           = context.Background
	exitFunc          func(int) = os.Exit
	mainStdout        io.Writer = os.Stdout
	depsFactory                 = newRuntimeDeps
)

func main() {
	exitFunc(runMain(backgroundContext(), mainStdout, depsFactory()))
}

func runMain(ctx context.Context, out io.Writer, deps runtimeDeps) int {
	logger := log.New(out, "", log.LstdFlags)
	if err := run(ctx, logger, deps); err != nil {
		logger.Print(err)
		return 1
	}
	return 0
}

func run(ctx context.Context, logger *log.Logger, deps runtimeDeps) error {
	cfg, err := deps.loadConfig("internal/platform/config")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	repository := deps.newRepository()
	provider := deps.newProvider(cfg.Provider.RegularURL, cfg.Provider.DutyURL, cfg.Sync.Timeout)
	catalogs := app.NewCatalogService(repository)
	syncService := app.NewSyncService(provider, repository)

	startupCtx, cancel := deps.withTimeout(ctx, cfg.Sync.Timeout)
	defer cancel()

	if err := syncService.SyncAll(startupCtx); err != nil {
		return fmt.Errorf("startup sync: %w", err)
	}

	runner := deps.newRunner(logger, cfg.Sync.Interval, cfg.Sync.Timeout)
	runner.Start(ctx, "refresh pharmacy catalogs", syncService.SyncAll)

	handler := pharmacyhttp.NewHandler(catalogs, syncService)
	server := deps.newServer(cfg.Server.Port, pharmacyhttp.NewRouter(handler))
	if err := server.Run(); err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}

func newRuntimeDeps() runtimeDeps {
	return runtimeDeps{
		loadConfig: platformconfig.Load,
		newRepository: func() app.Repository {
			return memory.NewRepository()
		},
		newProvider: func(regularURL, dutyURL string, timeout time.Duration) app.Provider {
			return boostr.NewClient(regularURL, dutyURL, timeout)
		},
		newRunner: func(logger *log.Logger, interval, timeout time.Duration) schedulerRunner {
			return scheduler.NewRunner(logger, interval, timeout)
		},
		newServer: func(port string, fiberApp *fiber.App) serverRunner {
			return platformhttp.NewServer(port, fiberApp)
		},
		withTimeout: context.WithTimeout,
	}
}
