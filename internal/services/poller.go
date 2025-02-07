package services

import (
    "github.com/docker/docker/api/types"
    "os/exec"
    "pinger/internal/domain/interfaces/handlers"
    "pinger/internal/domain/interfaces/services"
    "pinger/internal/domain/models"
    "time"
)

type Poller struct {
    broadcast handlers.Streamer
}

func NewPoller(broadcast handlers.Streamer) services.Poller {
    return Poller{broadcast}
}

func (p Poller) Poll(containers <-chan *types.Container) error {
    for c := range containers {
        ip := ipOfContainer(c)
        if ip == "" {
            continue
        }
        exit, pingTime, err := ping(ip)
        if err != nil && exit == 0 {
            return err
        }
        machine := models.Machine{
            IP:          ip,
            Success:     exit == 0,
            PingTime:    pingTime,
            LastSuccess: time.Now(),
        }
        err = p.broadcast.Send(&machine)
        if err != nil {
            return err
        }
    }
    return nil
}

func ipOfContainer(c *types.Container) string {
    var ip string
    for _, net := range c.NetworkSettings.Networks {
        if net.IPAddress != "" {
            ip = net.IPAddress
        }
    }
    return ip
}

func ping(ip string) (int, time.Duration, error) {
    start := time.Now()
    cmd := exec.Command("ping", "-c", "4", ip)
    err := cmd.Run()
    duration := time.Now().Sub(start)
    exit := cmd.ProcessState.ExitCode()
    return exit, duration, err
}
