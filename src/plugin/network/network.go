package network

import (
	"plugin_interface"
	"gopkg.in/yaml.v1"
	"time"
	"io/ioutil"
	"fmt"
	"strings"
	"strconv"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

type Config struct {
	iface string
}

type netStats struct {
	tx uint64
	rx uint64
	old_tx uint64
	old_rx uint64
	old_ts time.Time
	ts time.Time
}

func New(config map[string]interface{}, events chan plugin_interface.Event, update chan plugin_interface.Update) {
	c := loadConfig(config)
	str, _ := yaml.Marshal(config)
	log.Warning(string(str))
	var stats netStats
	stats.old_ts = time.Now()
	stats.ts = time.Now()
	var ev plugin_interface.Update
	//send sth at start of plugin, in case we dont get anything useful (like interface with no traffic)
	ev.FullText= fmt.Sprintf("%s!!", c.iface)
	ev.Color = "#999999"
	update <- ev
	Update(update,c,&stats)
	for {
		select {
		case _ = (<-events):
			Update(update,c,&stats)
		case <-time.After(time.Second):
			Update(update,c,&stats)
		}
	}

}

func loadConfig(raw map[string]interface{}) Config {
	var c Config
	c.iface = `lo`
	for key, value := range raw {
		converted, ok := value.(string)
		if ok {
			switch {
			case key == `iface`:
				c.iface=converted
				log.Warning("-- %s %s--", key, c.iface)

			}
		} else {
			log.Warning("-- %s--", key)
			_ = ok
		}
	}
	return c
}


func Update(update chan plugin_interface.Update, cfg Config, stats *netStats) {
	var ev plugin_interface.Update
	ev.Color=`#ffffdd`
	stats.old_ts = stats.ts
	rx, tx := getStats(cfg.iface)
	stats.ts = time.Now()

	// TODO: do same on bigger time diff
	if (stats.ts.UnixNano() < stats.old_ts.UnixNano()) {
		// we are in time machine.. or ntp changed clock
		stats.old_ts = stats.ts
		return
	}

	// either interface never seen packets, or it got recreated, reset it
	if rx == 0 && tx == 0 {
		stats.rx = 0
		stats.tx = 0
		stats.old_rx = 0
		stats.old_tx = 0
		return
	}

	// counter flipped, or interface recreated, reset to current value
	if ( stats.rx > rx || stats.tx > tx ) {
		stats.rx = rx
		stats.tx = tx
		stats.old_rx = rx
		stats.old_tx = tx
		return
	}
	//  init on first probe on empty interface
	if (stats.rx == 0 && rx > 0) {
		stats.rx = rx
		stats.tx = tx
	}
	// should be only useful data left
	stats.old_rx = stats.rx
	stats.old_tx = stats.tx
	stats.rx = rx
	stats.tx = tx
	rx_diff := stats.rx - stats.old_rx
	tx_diff := stats.tx - stats.old_tx
	ts_diff := float64(stats.ts.UnixNano() - stats.old_ts.UnixNano())
	ts_diff = ts_diff / 1000000000 //float64(time.Duration(time.Second).Nanoseconds()) // normalize
	if (ts_diff < 0.01) {
		return ; // quicker probing doesnt make sense, no div by 0, should probably return an error...
	}
	rx_bw := float64(rx_diff) / ts_diff
	tx_bw := float64(tx_diff) / ts_diff
	divider, unit := getUnit(rx_bw + tx_bw)
	ev.FullText = fmt.Sprintf(`%s:%6.3g/%6.3g %s`, cfg.iface,  rx_bw/divider, tx_bw/divider, unit)
	if strings.Contains(ev.FullText,`+`) {
		log.Warning("%v %v %v %v", rx_bw, tx_bw, divider, unit)
	}
	ev.ShortText = fmt.Sprintf(`-%s-`, cfg.iface)
	switch {
	case (rx_bw + tx_bw) < 50 * 1024:
		ev.Color = "#aaaaff"
	case (rx_bw + tx_bw) < 150 * 1024:
		ev.Color = "#aa33ff"
	case (rx_bw + tx_bw) < 450 * 1024:
		ev.Color = "#00ffff"
	case (rx_bw + tx_bw) < 1024 * 1024:
		ev.Color = "#00ff00"
	case (rx_bw + tx_bw) < 2048 * 1024:
		ev.Color = "#99ff00"
	case (rx_bw + tx_bw) < 4096 * 1024:
		ev.Color = "#ffff00"
	default:
		ev.Color = "#ff4400"
	}
	update <- ev
}

func getStats(iface string) (uint64, uint64) {
    raw_rx, _ := ioutil.ReadFile(fmt.Sprintf(`/sys/class/net/%s/statistics/rx_bytes`,iface))
	raw_tx, _ := ioutil.ReadFile(fmt.Sprintf(`/sys/class/net/%s/statistics/tx_bytes`,iface))
	str_rx := strings.TrimSpace(string(raw_rx))
	str_tx := strings.TrimSpace(string(raw_tx))
	rx, _ := strconv.ParseUint(string(str_rx),10,64)
	tx, _ := strconv.ParseUint(string(str_tx),10,64)
	return rx, tx
}


func getUnit(bytes float64) (divider float64, unit string) {
	switch {
	case bytes < 125 * 1024:
		return 1024/8, `Kb`
	case bytes < 100 * 1024 * 1024:
		return 1024 * 1024 / 8, `Mb`
	default:
		return 1024 * 1024 * 1024 / 8, `Gb`
	}
}
