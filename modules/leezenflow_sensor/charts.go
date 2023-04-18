// SPDX-License-Identifier: GPL-3.0-or-later

package leezenflow_sensor

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

var chartTemplate = module.Chart{
	// ID:    "random_%d",
	Title: "Leezenflow Sensor %s",
	// Units: "cel",
	Fam: "family",
	Ctx: "leezenflow_sensor.hardware",
}

func newChart(num, ctx, labels int, typ module.ChartType) *module.Chart {
	chart := chartTemplate.Copy()
	chart.ID = fmt.Sprintf(chart.ID, num)
	chart.Type = typ
	if ctx > 0 {
		chart.Ctx += fmt.Sprintf("_%d", ctx)
	}
	for i := 0; i < labels; i++ {
		chart.Labels = append(chart.Labels, module.Label{
			Key:   fmt.Sprintf("leezenflow_sensor_name_%d", i),
			Value: fmt.Sprintf("leezenflow_sensor_value_%d_%d", num, i),
		})
	}
	return chart
}

func newTempChart() *module.Chart {
	chart := chartTemplate.Copy()
	chart.ID = "temperature"
	chart.Type = "line"
	chart.Units = "celsius"
	chart.Title = fmt.Sprintf(chart.Title, "Temperature")
	// if ctx > 0 {
	// 	chart.Ctx += fmt.Sprintf("_%d", ctx)
	// }
	// for i := 0; i < labels; i++ {
	// 	chart.Labels = append(chart.Labels, module.Label{
	// 		Key:   fmt.Sprintf("leezenflow_sensor_name_%d", i),
	// 		Value: fmt.Sprintf("leezenflow_sensor_value_%d_%d", num, i),
	// 	})
	// }

	return chart
}

func newHumiChart() *module.Chart {
	chart := chartTemplate.Copy()
	chart.ID = "humidity"
	chart.Type = "line"
	chart.Units = "percent"
	chart.Title = fmt.Sprintf(chart.Title, "Humidity")
	// if ctx > 0 {
	// 	chart.Ctx += fmt.Sprintf("_%d", ctx)
	// }
	// for i := 0; i < labels; i++ {
	// 	chart.Labels = append(chart.Labels, module.Label{
	// 		Key:   fmt.Sprintf("leezenflow_sensor_name_%d", i),
	// 		Value: fmt.Sprintf("leezenflow_sensor_value_%d_%d", num, i),
	// 	})
	// }

	return chart
}

func newVoltageChart() *module.Chart {
	chart := chartTemplate.Copy()
	chart.ID = "voltage"
	chart.Type = "line"
	chart.Units = "volt"
	chart.Title = fmt.Sprintf(chart.Title, "Voltage")
	// if ctx > 0 {
	// 	chart.Ctx += fmt.Sprintf("_%d", ctx)
	// }
	// for i := 0; i < labels; i++ {
	// 	chart.Labels = append(chart.Labels, module.Label{
	// 		Key:   fmt.Sprintf("leezenflow_sensor_name_%d", i),
	// 		Value: fmt.Sprintf("leezenflow_sensor_value_%d_%d", num, i),
	// 	})
	// }

	return chart
}

func newPressureChart() *module.Chart {
	chart := chartTemplate.Copy()
	chart.ID = "pressure"
	chart.Type = "line"
	chart.Units = "pascal"
	chart.Title = fmt.Sprintf(chart.Title, "Pressure")
	// if ctx > 0 {
	// 	chart.Ctx += fmt.Sprintf("_%d", ctx)
	// }
	// for i := 0; i < labels; i++ {
	// 	chart.Labels = append(chart.Labels, module.Label{
	// 		Key:   fmt.Sprintf("leezenflow_sensor_name_%d", i),
	// 		Value: fmt.Sprintf("leezenflow_sensor_value_%d_%d", num, i),
	// 	})
	// }

	return chart
}
