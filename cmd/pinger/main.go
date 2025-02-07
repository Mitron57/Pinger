package main

import (
    "flag"
    "log"
    "pinger/internal/app"
)

func main() {
    path := flag.String("c", "config/config.yaml", "Path to config.yaml")
    flag.Parse()
    p := app.InitPinger(*path)
    if err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
