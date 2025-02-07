package app

import (
    "context"
    "fmt"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/client"
    "go.uber.org/zap"
    "log"
    "net/url"
    "os"
    "os/signal"
    "pinger/config"
    "pinger/internal/domain/interfaces/handlers"
    "pinger/internal/domain/interfaces/services"
    "pinger/internal/handlers/http"
    serviceImpl "pinger/internal/services"
    "runtime"
    "sync"
    "syscall"
    "time"
)

type Pinger struct {
    poller  services.Poller
    handler handlers.Streamer
    client  *client.Client
    logger  *zap.Logger
}

func initDockerClient() *client.Client {
    cli, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        log.Fatal(err)
    }
    return cli
}

func InitPinger(cfgPath string) *Pinger {
    cfg, err := config.ParseConfig(cfgPath)
    if err != nil {
        log.Fatal(err)
    }
    link, err := url.Parse(cfg.Pinger.Api)
    if err != nil {
        log.Fatal(err)
    }
    logger := zap.Must(zap.NewProduction())
    handler := http.NewHttpStreamer(link, logger)
    poller := serviceImpl.NewPoller(handler)
    dockerClient := initDockerClient()
    return &Pinger{
        poller:  poller,
        handler: handler,
        client:  dockerClient,
        logger:  logger,
    }
}

func (p *Pinger) Run() error {
    var wg sync.WaitGroup
    defer p.client.Close()
    defer wg.Wait()
    containers := make(chan *types.Container, 100)
    ctx, cancel := context.WithCancel(context.Background())
    go handleSignals(cancel)
    for range runtime.NumCPU() {
        wg.Add(1)
        go p.poll(ctx, containers, &wg)
    }
    return p.mainLoop(ctx, containers)
}

func (p *Pinger) poll(ctx context.Context, containers chan *types.Container, wg *sync.WaitGroup) {
    for {
        select {
        case <-ctx.Done():
            wg.Done()
            return
        default:
            if err := p.poller.Poll(containers); err != nil {
                p.logger.Error(err.Error())
            }
        }
    }
}

func (p *Pinger) mainLoop(ctx context.Context, containers chan *types.Container) error {
    ticker := time.Tick(time.Second * 10)
loop:
    for {
        select {
        case <-ticker:
            list, err := p.client.ContainerList(ctx, container.ListOptions{All: true})
            if err != nil {
                return err
            }
            for _, c := range list {
                containers <- &c
            }
        case <-ctx.Done():
            close(containers)
            break loop
        }
    }
    return nil
}

func handleSignals(cancel context.CancelFunc) {
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
    <-signals
    fmt.Println("\nGraceful shutdown, wait a few seconds...")
    cancel()
}
