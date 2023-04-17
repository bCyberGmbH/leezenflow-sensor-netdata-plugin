// SPDX-License-Identifier: GPL-3.0-or-later

package leezenflow_sensor

import (
	"bufio"
	"math/rand"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"

	"github.com/tarm/serial"
)

func init() {
	module.Register("leezenflow_sensor", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery:        module.UpdateEvery,
			AutoDetectionRetry: module.AutoDetectionRetry,
			Priority:           module.Priority,
			Disabled:           true,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *LeezenflowSensor {

	return &LeezenflowSensor{
		Config: Config{
			Charts: ConfigCharts{
				Num:  1,
				Dims: 4,
			},
		},

		randInt:       func() int64 { return rand.Int63n(100) },
		collectedDims: make(map[string]bool),
	}
}

type (
	Config struct {
		Charts ConfigCharts `yaml:"charts"`
	}
	ConfigCharts struct {
		Type     string `yaml:"type"`
		Num      int    `yaml:"num"`
		Contexts int    `yaml:"contexts"`
		Dims     int    `yaml:"dimensions"`
		Labels   int    `yaml:"labels"`
	}
)

type LeezenflowSensor struct {
	module.Base // should be embedded by every module
	Config      `yaml:",inline"`

	randInt       func() int64
	charts        *module.Charts
	collectedDims map[string]bool

	lastTemperature, lastHumidty, lastVoltage int64
}

func (l *LeezenflowSensor) initSerial() error {
	//
	//flagSerialDevice := "/dev/ttyUSB0"
	flagSerialDevice := "/dev/pts/3"

	config := &serial.Config{
		Name: flagSerialDevice,
		Baud: 9600,
		// ReadTimeout: 1,
		Size: 8,
	}

	stream, err := serial.OpenPort(config)
	if err != nil {
		l.Errorf("Could not open device %s: %v\n", flagSerialDevice, err)
		return err
	}

	go func() {

		// log.Printf("Listening for messages on %s", flagSerialDevice)

		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if line == "" {
				continue
			}

			// Hier Nachricht vom ESP parsen und zwischenspeichern

			// fmt.Fprintf(l.tempString, "%s,", line)
			l.Println(line)

		}
		if err := scanner.Err(); err != nil {
			l.Errorf("Error reading device %s: %v\n", flagSerialDevice, err)
		}
	}()

	return nil
}

func (l *LeezenflowSensor) Init() bool {
	err := l.validateConfig()
	if err != nil {
		l.Errorf("config validation: %v", err)
		return false
	}

	err = l.initSerial()
	if err != nil {
		l.Errorf("serial initialization: %v", err)
		return false
	}

	charts, err := l.initCharts()
	if err != nil {
		l.Errorf("charts init: %v", err)
		return false
	}
	l.charts = charts
	return true
}

func (l *LeezenflowSensor) Check() bool {
	return len(l.Collect()) > 0
}

func (l *LeezenflowSensor) Charts() *module.Charts {
	return l.charts
}

func (l *LeezenflowSensor) Collect() map[string]int64 {
	mx, err := l.collect()
	if err != nil {
		l.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (l *LeezenflowSensor) Cleanup() {}
