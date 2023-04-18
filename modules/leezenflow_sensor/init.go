// SPDX-License-Identifier: GPL-3.0-or-later

package leezenflow_sensor

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (l *LeezenflowSensor) initCharts() (*module.Charts, error) {
	charts := &module.Charts{}

	if err := charts.Add(newTempChart()); err != nil {
		return nil, err
	}

	if err := charts.Add(newHumiChart()); err != nil {
		return nil, err
	}

	if err := charts.Add(newVoltageChart()); err != nil {
		return nil, err
	}

	if err := charts.Add(newPressureChart()); err != nil {
		return nil, err
	}

	return charts, nil
}

func (l *LeezenflowSensor) initDims() {
	for _, chart := range *l.Charts() {
		name := chart.ID
		id := fmt.Sprintf("%s_%s", chart.ID, name)

		dim := &module.Dim{ID: id, Name: name}
		if err := chart.AddDim(dim); err != nil {
			l.Warning(err)
		}
		chart.MarkNotCreated()
	}
}
