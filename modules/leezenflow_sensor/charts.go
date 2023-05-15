// SPDX-License-Identifier: GPL-3.0-or-later

package leezenflow_sensor

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

var chartTemplate = module.Chart{
	Title: "Leezenflow Sensor %s",
	Ctx:   "leezenflow_sensor.hardware",
	Type:  "line",
}

func newTempChart() *module.Chart {
	chart := chartTemplate.Copy()
	chart.ID = "temperature"
	chart.Units = "celsius"
	chart.Fam = "temperature"
	chart.Title = fmt.Sprintf(chart.Title, "Temperature")

	return chart
}

func newHumiChart() *module.Chart {
	chart := chartTemplate.Copy()
	chart.ID = "humidity"
	chart.Units = "percent"
	chart.Fam = "humidity"
	chart.Title = fmt.Sprintf(chart.Title, "Humidity")

	return chart
}

func newVoltageChart() *module.Chart {
	chart := chartTemplate.Copy()
	chart.ID = "voltage"
	chart.Units = "volt"
	chart.Fam = "voltage"
	chart.Title = fmt.Sprintf(chart.Title, "Voltage")

	return chart
}

func newPressureChart() *module.Chart {
	chart := chartTemplate.Copy()
	chart.ID = "pressure"
	chart.Units = "pascal"
	chart.Fam = "pressure"
	chart.Title = fmt.Sprintf(chart.Title, "Pressure")

	return chart
}
