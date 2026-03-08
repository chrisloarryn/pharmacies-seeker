package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pharmacies-seeker/internal/pharmacies/adapters/repository/memory"
	"pharmacies-seeker/internal/pharmacies/app"
	"pharmacies-seeker/internal/pharmacies/domain"
	platformconfig "pharmacies-seeker/internal/platform/config"
	"pharmacies-seeker/internal/platform/scheduler"
)

type providerSpy struct {
	data  map[domain.CatalogKind][]domain.Pharmacy
	errs  map[domain.CatalogKind]error
	calls []domain.CatalogKind
}

func (s *providerSpy) Fetch(_ context.Context, kind domain.CatalogKind) ([]domain.Pharmacy, error) {
	s.calls = append(s.calls, kind)
	if err := s.errs[kind]; err != nil {
		return nil, err
	}
	return append([]domain.Pharmacy(nil), s.data[kind]...), nil
}

type runnerSpy struct {
	ctx   context.Context
	name  string
	task  scheduler.Task
	start int
}

func (s *runnerSpy) Start(ctx context.Context, name string, task scheduler.Task) {
	s.ctx = ctx
	s.name = name
	s.task = task
	s.start++
}

type serverSpy struct {
	runErr error
	runs   int
}

func (s *serverSpy) Run() error {
	s.runs++
	return s.runErr
}

func TestNewRuntimeDepsBuildsDefaultDependencies(t *testing.T) {
	deps := newRuntimeDeps()

	assert.NotNil(t, deps.loadConfig)
	assert.NotNil(t, deps.newRepository)
	assert.NotNil(t, deps.newProvider)
	assert.NotNil(t, deps.newRunner)
	assert.NotNil(t, deps.newServer)
	assert.NotNil(t, deps.withTimeout)

	repository := deps.newRepository()
	provider := deps.newProvider("https://example.com/regular", "https://example.com/duty", time.Second)
	runner := deps.newRunner(log.New(io.Discard, "", 0), time.Minute, time.Second)
	server := deps.newServer("0", fiber.New())
	ctx, cancel := deps.withTimeout(context.Background(), time.Millisecond)
	cancel()

	assert.NotNil(t, repository)
	assert.NotNil(t, provider)
	assert.NotNil(t, runner)
	assert.NotNil(t, server)
	assert.NotNil(t, ctx)
}

func TestRunSuccessStartsSchedulerAndRunsServer(t *testing.T) {
	cfg := platformconfig.Config{
		Server:   platformconfig.ServerConfig{Port: "8080"},
		Provider: platformconfig.ProviderConfig{RegularURL: "https://example.com/regular", DutyURL: "https://example.com/duty"},
		Sync:     platformconfig.SyncConfig{Interval: time.Minute, Timeout: time.Second},
	}
	provider := &providerSpy{
		data: map[domain.CatalogKind][]domain.Pharmacy{
			domain.CatalogRegular: {{ID: "regular"}},
			domain.CatalogDuty:    {{ID: "duty"}},
		},
		errs: map[domain.CatalogKind]error{},
	}
	runner := &runnerSpy{}
	server := &serverSpy{}
	ctx := context.WithValue(context.Background(), "request-id", "abc")

	var (
		loadPath       string
		regularURL     string
		dutyURL        string
		timeout        time.Duration
		serverPort     string
		serverApp      *fiber.App
		cancelCalled   bool
		timeoutParent  context.Context
		timeoutApplied time.Duration
	)

	err := run(ctx, log.New(io.Discard, "", 0), runtimeDeps{
		loadConfig: func(path string) (platformconfig.Config, error) {
			loadPath = path
			return cfg, nil
		},
		newRepository: func() app.Repository {
			return memory.NewRepository()
		},
		newProvider: func(gotRegularURL, gotDutyURL string, gotTimeout time.Duration) app.Provider {
			regularURL = gotRegularURL
			dutyURL = gotDutyURL
			timeout = gotTimeout
			return provider
		},
		newRunner: func(_ *log.Logger, interval, gotTimeout time.Duration) schedulerRunner {
			assert.Equal(t, cfg.Sync.Interval, interval)
			assert.Equal(t, cfg.Sync.Timeout, gotTimeout)
			return runner
		},
		newServer: func(port string, fiberApp *fiber.App) serverRunner {
			serverPort = port
			serverApp = fiberApp
			return server
		},
		withTimeout: func(parent context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
			timeoutParent = parent
			timeoutApplied = duration
			return parent, func() {
				cancelCalled = true
			}
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "internal/platform/config", loadPath)
	assert.Equal(t, cfg.Provider.RegularURL, regularURL)
	assert.Equal(t, cfg.Provider.DutyURL, dutyURL)
	assert.Equal(t, cfg.Sync.Timeout, timeout)
	assert.Equal(t, cfg.Server.Port, serverPort)
	assert.NotNil(t, serverApp)
	assert.Same(t, ctx, timeoutParent)
	assert.Equal(t, cfg.Sync.Timeout, timeoutApplied)
	assert.True(t, cancelCalled)
	assert.Equal(t, []domain.CatalogKind{domain.CatalogRegular, domain.CatalogDuty}, provider.calls)
	assert.Equal(t, 1, runner.start)
	assert.Same(t, ctx, runner.ctx)
	assert.Equal(t, "refresh pharmacy catalogs", runner.name)
	assert.NotNil(t, runner.task)
	assert.Equal(t, 1, server.runs)
}

func TestRunReturnsLoadConfigError(t *testing.T) {
	err := run(context.Background(), log.New(io.Discard, "", 0), runtimeDeps{
		loadConfig: func(string) (platformconfig.Config, error) {
			return platformconfig.Config{}, errors.New("boom")
		},
	})

	require.EqualError(t, err, "load config: boom")
}

func TestRunReturnsStartupSyncError(t *testing.T) {
	cfg := platformconfig.Config{
		Server:   platformconfig.ServerConfig{Port: "8080"},
		Provider: platformconfig.ProviderConfig{RegularURL: "regular", DutyURL: "duty"},
		Sync:     platformconfig.SyncConfig{Interval: time.Minute, Timeout: time.Second},
	}
	provider := &providerSpy{
		data: map[domain.CatalogKind][]domain.Pharmacy{},
		errs: map[domain.CatalogKind]error{
			domain.CatalogRegular: errors.New("boom"),
		},
	}

	err := run(context.Background(), log.New(io.Discard, "", 0), runtimeDeps{
		loadConfig: func(string) (platformconfig.Config, error) {
			return cfg, nil
		},
		newRepository: func() app.Repository {
			return memory.NewRepository()
		},
		newProvider: func(string, string, time.Duration) app.Provider {
			return provider
		},
		newRunner: func(*log.Logger, time.Duration, time.Duration) schedulerRunner {
			t.Fatal("runner should not start when startup sync fails")
			return nil
		},
		newServer: func(string, *fiber.App) serverRunner {
			t.Fatal("server should not start when startup sync fails")
			return nil
		},
		withTimeout: context.WithTimeout,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "startup sync")
	assert.Contains(t, err.Error(), "regular catalog fetch: boom")
}

func TestRunReturnsServerError(t *testing.T) {
	cfg := platformconfig.Config{
		Server:   platformconfig.ServerConfig{Port: "8080"},
		Provider: platformconfig.ProviderConfig{RegularURL: "regular", DutyURL: "duty"},
		Sync:     platformconfig.SyncConfig{Interval: time.Minute, Timeout: time.Second},
	}
	provider := &providerSpy{
		data: map[domain.CatalogKind][]domain.Pharmacy{
			domain.CatalogRegular: {{ID: "regular"}},
			domain.CatalogDuty:    {{ID: "duty"}},
		},
		errs: map[domain.CatalogKind]error{},
	}
	server := &serverSpy{runErr: errors.New("listen failed")}

	err := run(context.Background(), log.New(io.Discard, "", 0), runtimeDeps{
		loadConfig: func(string) (platformconfig.Config, error) {
			return cfg, nil
		},
		newRepository: func() app.Repository {
			return memory.NewRepository()
		},
		newProvider: func(string, string, time.Duration) app.Provider {
			return provider
		},
		newRunner: func(*log.Logger, time.Duration, time.Duration) schedulerRunner {
			return &runnerSpy{}
		},
		newServer: func(string, *fiber.App) serverRunner {
			return server
		},
		withTimeout: context.WithTimeout,
	})

	require.EqualError(t, err, "run server: listen failed")
	assert.Equal(t, 1, server.runs)
}

func TestRunMainReturnsNonZeroAndLogsError(t *testing.T) {
	var out bytes.Buffer

	code := runMain(context.Background(), &out, runtimeDeps{
		loadConfig: func(string) (platformconfig.Config, error) {
			return platformconfig.Config{}, errors.New("boom")
		},
	})

	assert.Equal(t, 1, code)
	assert.Contains(t, out.String(), "load config: boom")
}

func TestMainExitsWithRunMainStatus(t *testing.T) {
	previousExit := exitFunc
	previousOut := mainStdout
	previousDeps := depsFactory
	previousBackground := backgroundContext

	t.Cleanup(func() {
		exitFunc = previousExit
		mainStdout = previousOut
		depsFactory = previousDeps
		backgroundContext = previousBackground
	})

	provider := &providerSpy{
		data: map[domain.CatalogKind][]domain.Pharmacy{
			domain.CatalogRegular: {{ID: "regular"}},
			domain.CatalogDuty:    {{ID: "duty"}},
		},
		errs: map[domain.CatalogKind]error{},
	}
	cfg := platformconfig.Config{
		Server:   platformconfig.ServerConfig{Port: "8080"},
		Provider: platformconfig.ProviderConfig{RegularURL: "regular", DutyURL: "duty"},
		Sync:     platformconfig.SyncConfig{Interval: time.Minute, Timeout: time.Second},
	}

	var exitCode int
	exitFunc = func(code int) {
		exitCode = code
	}
	mainStdout = io.Discard
	backgroundContext = context.Background
	depsFactory = func() runtimeDeps {
		return runtimeDeps{
			loadConfig: func(string) (platformconfig.Config, error) {
				return cfg, nil
			},
			newRepository: func() app.Repository {
				return memory.NewRepository()
			},
			newProvider: func(string, string, time.Duration) app.Provider {
				return provider
			},
			newRunner: func(*log.Logger, time.Duration, time.Duration) schedulerRunner {
				return &runnerSpy{}
			},
			newServer: func(string, *fiber.App) serverRunner {
				return &serverSpy{}
			},
			withTimeout: context.WithTimeout,
		}
	}

	main()

	assert.Equal(t, 0, exitCode)
}
