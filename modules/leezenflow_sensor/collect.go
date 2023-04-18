// SPDX-License-Identifier: GPL-3.0-or-later

package leezenflow_sensor

// import (
// 	"fmt"

// 	"github.com/netdata/go.d.plugin/agent/module"
// )

// func (l *LeezenflowSensor) collect() (map[string]int64, error) {
// 	collected := make(map[string]int64)

// 	for _, chart := range *l.Charts() {
// 		l.collectChart(collected, chart)
// 	}
// 	return collected, nil
// }

// func (l *LeezenflowSensor) collectChart(collected map[string]int64, chart *module.Chart) {
// 	num := l.Config.Charts.Dims

// 	for i := 0; i < num; i++ {
// 		name := fmt.Sprintf("random%d", i)
// 		id := fmt.Sprintf("%s_%s", chart.ID, name)

// 		if !l.collectedDims[id] {
// 			l.collectedDims[id] = true

// 			dim := &module.Dim{ID: id, Name: name}
// 			if err := chart.AddDim(dim); err != nil {
// 				l.Warning(err)
// 			}
// 			chart.MarkNotCreated()
// 		}
// 		if i%2 == 0 {
// 			collected[id] = l.randInt()
// 		} else {
// 			collected[id] = -l.randInt()
// 		}
// 	}
// }
