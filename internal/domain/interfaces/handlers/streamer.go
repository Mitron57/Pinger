package handlers

import "pinger/internal/domain/models"

type Streamer interface {
    Send(machine *models.Machine) error
}
