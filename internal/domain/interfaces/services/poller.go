package services

import (
    "github.com/docker/docker/api/types"
)

type Poller interface {
    Poll(containers <-chan *types.Container) error
}
