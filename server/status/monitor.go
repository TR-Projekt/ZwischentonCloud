package status

import (
	"time"
)

type MonitorNode struct {
	Type     string
	Location string
	LastSeen time.Time
}
