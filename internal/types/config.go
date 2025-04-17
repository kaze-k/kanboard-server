package types

import (
	"time"
)

type Env struct {
	Mode string
}

type Server struct {
	Name    string
	Host    string
	Port    string
	Timeout time.Duration
}

type JWT struct {
	Secret                  string
	KanboardTokenExpiration time.Duration
	AdminTokenExpiration    time.Duration
}

type File struct {
	Path   string
	Static string
}
