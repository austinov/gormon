package monitor

type Monitor interface {
	Process(host, output string)
}
