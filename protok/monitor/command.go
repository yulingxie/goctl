package monitor

import (
	"errors"

	"github.com/urfave/cli"
)

func TrafficMonitor(context *cli.Context) error {
	protocol := context.String("protocol")
	serverAddr := context.String("server")
	packets := context.StringSlice("filter")
	heartbeat := context.Bool("heartbeat")
	monitor := NewMonitor(
		WithProtocol(protocol),
		WithServerAddr(serverAddr),
		WithHeartbeat(heartbeat),
		WithPackets(packets),
	)
	device := monitor.SelectDevice()
	if len(device) == 0 {
		return errors.New("未选择监听IP")
	}
	monitor.MonitorDevice(device)
	return nil
}
