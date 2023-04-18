// SPDX-License-Identifier: GPL-3.0-or-later

package leezenflow_sensor

import (
	"bufio"
	"math/rand"
	"strconv"
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
			Serial: ConfigSerial{
				Port: "/dev/ttyUSB0",
				Baud: 115200,
			},
		},

		randInt:       func() int64 { return rand.Int63n(100) },
		collectedDims: make(map[string]bool),
	}
}

type (
	Config struct {
		Charts ConfigCharts `yaml:"charts"`
		Serial ConfigSerial `yaml:"serial"`
	}
	ConfigCharts struct {
		Type     string `yaml:"type"`
		Num      int    `yaml:"num"`
		Contexts int    `yaml:"contexts"`
		Dims     int    `yaml:"dimensions"`
		Labels   int    `yaml:"labels"`
	}
	ConfigSerial struct {
		Port string `yaml:"port"`
		Baud int    `yaml:"baud"`
	}
)

type LeezenflowSensor struct {
	module.Base // should be embedded by every module
	Config      `yaml:",inline"`

	randInt       func() int64
	charts        *module.Charts
	collectedDims map[string]bool

	lastTemperature, lastHumidty, lastVoltage, lastPressure int64
}

func (l *LeezenflowSensor) initSerial() error {

	l.Printf("Trying to open serial device '%s' with baud rate %d", l.Config.Serial.Port, l.Config.Serial.Baud)

	flagSerialDevice := l.Config.Serial.Port

	config := &serial.Config{
		Name: flagSerialDevice,
		Baud: l.Config.Serial.Baud,
		// ReadTimeout: 1,
		Size: 8,
	}

	stream, err := serial.OpenPort(config)
	if err != nil {
		l.Errorf("Could not open device '%s': %v\n", flagSerialDevice, err)
		return err
	}

	go func() {

		l.Printf("Listening for messages on '%s'", flagSerialDevice)

		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if line == "" {
				continue
			}

			// Hier Nachricht vom ESP parsen und zwischenspeichern

			// fmt.Fprintf(l.tempString, "%s,", line)
			l.Printf("Received: '%s'", line)
			// l.lastTemperature = 10

			// Split line
			dataFields := strings.Split(line, ":")
			if len(dataFields) != 2 {
				l.Warningf("Could not split '%s' into two fields", dataFields[1])
				continue
			}

			// Try to parse second field into float
			val, err := strconv.ParseFloat(strings.TrimSpace(dataFields[1]), 64)
			if err != nil {
				l.Warningf("Could not parse '%s' as float", dataFields[1])
				continue
			}

			if dataFields[0] == "C" {
				l.lastTemperature = int64(val * 10)
			} else if dataFields[0] == "H" {
				l.lastHumidty = int64(val)
			} else if dataFields[0] == "P" {
				l.lastPressure = int64(val)
			} else if dataFields[0] == "V" {
				l.lastTemperature = int64(val * 100)
			} else {
				l.Warningf("Unknown prefix '%s'", dataFields[0])
			}
		}
		if err := scanner.Err(); err != nil {
			l.Errorf("Error reading device %s: %v\n", flagSerialDevice, err)
		}
	}()

	return nil
}

func (l *LeezenflowSensor) Init() bool {
	err := l.initSerial()
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

	l.initDims()

	// first measurement seems to discarded anyway?
	l.lastTemperature = 10
	l.lastHumidty = 10
	l.lastVoltage = 10
	l.lastPressure = 10

	return true
}

func (l *LeezenflowSensor) Check() bool {
	return len(l.Collect()) > 0
}

func (l *LeezenflowSensor) Charts() *module.Charts {
	return l.charts
}

func (l *LeezenflowSensor) Collect() map[string]int64 {
	collected := make(map[string]int64)

	// Danger and not good practice
	// Hard coding these :(
	if l.lastTemperature != 0 {
		collected["temperature_temperature"] = l.lastTemperature
		l.lastTemperature = 0
	}

	if l.lastHumidty != 0 {
		collected["humidity_humidity"] = l.lastHumidty
		l.lastHumidty = 0
	}

	if l.lastVoltage != 0 {
		collected["voltage_voltage"] = l.lastVoltage
		l.lastVoltage = 0
	}

	if l.lastPressure != 0 {
		collected["pressure_pressure"] = l.lastPressure
		l.lastPressure = 0
	}

	return collected
}

func (l *LeezenflowSensor) Cleanup() {}
