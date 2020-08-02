package domain

import (
	"encoding/json"
	"time"
)

type Entry struct {
	Time        time.Time       `json:"time"`
	NSec        int64           `json:"nsec,string"`
	Namespace   string          `json:"namespace"`
	Source      string          `json:"source"`
	Host        string          `json:"host"`
	Level       string          `json:"level"`
	TraceID     string          `json:"trace_id"`
	Message     string          `json:"message"`
	RemoteIP    string          `json:"remote_ip"`
	Params      json.RawMessage `json:"params"`
	BuildCommit string          `json:"build_commit"`
	ConfigHash  string          `json:"config_hash"`
}
