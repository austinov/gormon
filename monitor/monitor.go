package monitor

import "io"

type Monitor interface {
	io.Closer
	Process(host, output string)
}
