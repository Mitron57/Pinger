package models

import "time"

type Machine struct {
    IP          string        `json:"ip"`
    PingTime    time.Duration `json:"ping_time"`
    Success     bool          `json:"success"`
    LastSuccess time.Time     `json:"last_success"`
}
