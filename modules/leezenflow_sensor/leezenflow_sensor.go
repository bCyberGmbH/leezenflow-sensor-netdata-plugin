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

			l.Printf("Received: '%s'", line)

			// Split line into measurements
			dataFields := strings.Split(line, ",")
			if len(dataFields) != 4 {
				l.Warningf("Could not split '%s' into 4 fields", dataFields[1])
				continue
			}

			// Try to parse first field into int (Celsius * 10)
			val, err := strconv.ParseInt(strings.TrimSpace(dataFields[0]), 10, 64)
			if err != nil {
				l.Warningf("Could not parse '%s' as int", dataFields[0])
			} else {
				l.lastTemperature = val
			}

			// Try to parse second field into int (Humidity)
			val, err = strconv.ParseInt(strings.TrimSpace(dataFields[1]), 10, 64)
			if err != nil {
				l.Warningf("Could not parse '%s' as int for humidity", dataFields[1])
			} else {
				l.lastHumidty = val
			}

			// Try to parse third field into int (Pressure)
			val, err = strconv.ParseInt(strings.TrimSpace(dataFields[2]), 10, 64)
			if err != nil {
				l.Warningf("Could not parse '%s' as int for pressure", dataFields[2])
			} else {
				l.lastPressure = val
			}

			// Try to parse fourth field into int (Voltage)
			val, err = strconv.ParseInt(strings.TrimSpace(dataFields[3]), 10, 64)
			if err != nil {
				l.Warningf("Could not parse '%s' as float for voltage", dataFields[3])
			} else {
				l.lastVoltage = val
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
	if l.lastTemperature != 0 && l.lastHumidty != 0 && l.lastPressure != 0 {
		collected["temperature_temperature"] = l.lastTemperature
		collected["humidity_humidity"] = l.lastHumidty
		collected["voltage_voltage"] = l.lastVoltage
		collected["pressure_pressure"] = l.lastPressure
	}

	return collected
}

func (l *LeezenflowSensor) Cleanup() {}
