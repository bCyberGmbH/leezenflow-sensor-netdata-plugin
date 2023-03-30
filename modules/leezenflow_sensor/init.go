// SPDX-License-Identifier: GPL-3.0-or-later

package leezenflow_sensor

import (
	"errors"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (l *LeezenflowSensor) validateConfig() error {
	if l.Config.Charts.Num <= 0 && l.Config.HiddenCharts.Num <= 0 {
		return errors.New("'charts->num' or `hidden_charts->num` must be > 0")
	}
	if l.Config.Charts.Num > 0 && l.Config.Charts.Dims <= 0 {
		return errors.New("'charts->dimensions' must be > 0")
	}
	if l.Config.HiddenCharts.Num > 0 && l.Config.HiddenCharts.Dims <= 0 {
		return errors.New("'hidden_charts->dimensions' must be > 0")
	}
	return nil
}

func (l *LeezenflowSensor) initCharts() (*module.Charts, error) {
	charts := &module.Charts{}

	var ctx int
	v := calcContextEvery(l.Config.Charts.Num, l.Config.Charts.Contexts)
	for i := 0; i < l.Config.Charts.Num; i++ {
		if i != 0 && v != 0 && ctx < (l.Config.Charts.Contexts-1) && i%v == 0 {
			ctx++
		}
		chart := newChart(i, ctx, l.Config.Charts.Labels, module.ChartType(l.Config.Charts.Type))

		if err := charts.Add(chart); err != nil {
			return nil, err
		}
	}

	ctx = 0
	v = calcContextEvery(l.Config.HiddenCharts.Num, l.Config.HiddenCharts.Contexts)
	for i := 0; i < l.Config.HiddenCharts.Num; i++ {
		if i != 0 && v != 0 && ctx < (l.Config.HiddenCharts.Contexts-1) && i%v == 0 {
			ctx++
		}
		chart := newHiddenChart(i, ctx, l.Config.HiddenCharts.Labels, module.ChartType(l.Config.HiddenCharts.Type))

		if err := charts.Add(chart); err != nil {
			return nil, err
		}
	}

	return charts, nil
}

func calcContextEvery(charts, contexts int) int {
	if contexts <= 1 {
		return 0
	}
	if contexts > charts {
		return 1
	}
	return charts / contexts
}
