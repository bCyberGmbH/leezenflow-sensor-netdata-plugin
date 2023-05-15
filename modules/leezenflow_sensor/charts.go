// SPDX-License-Identifier: GPL-3.0-or-later

package leezenflow_sensor

import (
	"github.com/netdata/go.d.plugin/agent/module"
)

// var chartTemplate = module.Chart{
// 	Title: "Leezenflow Sensor %s",
// 	Ctx:   "leezenflow_sensor.hardware",
// 	Type:  "line",
// }

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Chart is an alias for module.Chart
	Chart = module.Chart
	// Dims is an alias for module.Dims
	Dims = module.Dims
	// Dim is an alias for module.Dim
	Dim = module.Dim
)

var leezenflowCharts = Charts{
	{
		ID:    "temperature_case",
		Title: "Temperature inside the case",
		Units: "celsius",
		Fam:   "temperature",
		Ctx:   "leezenflow.temperature_case",
		Dims: Dims{
			{ID: "temperature", Name: "Temperature", Div: 10},
		},
	},
	{
		ID:    "pressure_case",
		Title: "Barometric pressure inside the case",
		Units: "hectopascal",
		Fam:   "pressure",
		Ctx:   "leezenflow.pressure_case",
		Dims: Dims{
			{ID: "pressure", Name: "Pressure", Div: 100},
		},
	},
	{
		ID:    "humidity_case",
		Title: "Humidity inside the case",
		Units: "percent",
		Fam:   "humidity",
		Ctx:   "leezenflow.humidity_case",
		Dims: Dims{
			{ID: "humidity", Name: "Humidity"},
		},
	},
	{
		ID:    "voltage_source",
		Title: "Voltage of the power source",
		Units: "volt",
		Fam:   "voltage",
		Ctx:   "leezenflow.voltage_source",
		Dims: Dims{
			{ID: "voltage", Name: "Voltage", Div: 1000},
		},
	},
}
