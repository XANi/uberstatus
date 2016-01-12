package network

import (
	"github.com/XANi/uberstatus/uber"
	"gopkg.in/yaml.v1"
	"time"
	"io/ioutil"
	"fmt"
	"net"
	"strings"
	"strconv"
	"github.com/op/go-logging"
	"github.com/VividCortex/ewma"
)

var log = logging.MustGetLogger("main")

type Config struct {
	iface string
}

type netStats struct {
	ip string
	tx uint64
	rx uint64
	oldTx uint64
	oldRx uint64
	ewmaTx ewma.MovingAverage
	ewmaRx ewma.MovingAverage
	oldTs time.Time
	ts time.Time
}

const ShowFirstAddr = 0
const ShowSecondAddr = 1
const ShowAllAddr = -1


func Run(cfg uber.PluginConfig) {
	c := loadConfig(cfg.Config)
	str, _ := yaml.Marshal(cfg.Config)
	log.Warning(string(str))
	var stats netStats
	stats.ewmaRx = ewma.NewMovingAverage(5)
	stats.ewmaTx = ewma.NewMovingAverage(5)
	stats.oldTs = time.Now()
	stats.ts = time.Now()
	var ev uber.Update
	//send sth at start of plugin, in case we dont get anything useful (like interface with no traffic)
	ev.FullText= fmt.Sprintf("%s??", c.iface)
	ev.Color = "#999999"
	cfg.Update <- ev
	Update(cfg.Update,c,&stats)
	for {
		select {
		case ev := <-cfg.Events:
			if ev.Button == 1 {
				UpdateAddr(cfg.Update,c.iface,ShowFirstAddr)
			} else if ev.Button == 3 {
				UpdateAddr(cfg.Update, c.iface,ShowSecondAddr)
			} else {
				UpdateAddr(cfg.Update, c.iface,ShowAllAddr)
			}
			select {
			case _ = <-cfg.Events:
			case <-time.After(10 * time.Second):
			}
		case <-time.After(time.Second):
			Update(cfg.Update,c,&stats)
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



func UpdateAddr(update chan uber.Update, ifname string, addr_id int) {
	var ev uber.Update
	ev.Color=`#aaffaa`
	ifaces, _ := net.Interfaces()
end:
	for _,iface := range ifaces {
		if iface.Name == ifname {
			v, _ := iface.Addrs()
			if len(v) <= addr_id {
				break end
			}
			if addr_id < 0 {
				ev.FullText = fmt.Sprintf("%+v", v)
			} else {
				ev.FullText = fmt.Sprintf("%+v", v[addr_id])
			}
			update <- ev
			return
		}
	}

	ev.FullText = fmt.Sprintf("%s??",ifname)
	update <- ev
}


func Update(update chan uber.Update, cfg Config, stats *netStats) {
	var ev uber.Update
	ev.Color=`#ffffdd`
	ev.FullText= fmt.Sprintf("%s!!", cfg.iface)
	stats.oldTs = stats.ts
	rx, tx := getStats(cfg.iface)
	stats.ts = time.Now()

	// TODO: do same on bigger time diff
	if (stats.ts.UnixNano() < stats.oldTs.UnixNano()) {
		// we are in time machine.. or ntp changed clock
		stats.oldTs = stats.ts
		return
	}

	// either interface never seen packets, or it got recreated, reset it
	if rx == 0 && tx == 0 {
		stats.rx = 0
		stats.tx = 0
		stats.oldRx = 0
		stats.oldTx = 0
		update <- ev
		return
	}

	// counter flipped, or interface recreated, reset to current value
	if ( stats.rx > rx || stats.tx > tx ) {
		stats.rx = rx
		stats.tx = tx
		stats.oldRx = rx
		stats.oldTx = tx
		return
	}
	//  init on first probe on empty interface
	if (stats.rx == 0 && rx > 0) {
		stats.rx = rx
		stats.tx = tx
	}
	// should be only useful data left
	stats.oldRx = stats.rx
	stats.oldTx = stats.tx
	stats.rx = rx
	stats.tx = tx
	rxDiff := stats.rx - stats.oldRx
	txDiff := stats.tx - stats.oldTx
	tsDiff := float64(stats.ts.UnixNano() - stats.oldTs.UnixNano())
	tsDiff = tsDiff / 1000000000 //float64(time.Duration(time.Second).Nanoseconds()) // normalize
	if (tsDiff < 0.01) {
		return ; // quicker probing doesnt make sense, no div by 0, should probably return an error...
	}
	rxBw := float64(rxDiff) / tsDiff
	txBw := float64(txDiff) / tsDiff
	stats.ewmaRx.Add(rxBw)
	stats.ewmaTx.Add(txBw)
	rxAvg := stats.ewmaRx.Value()
	txAvg := stats.ewmaTx.Value()
	divider, unit := getUnit(rxAvg + txAvg)
	// if speed is very low alias it to 0
	if rxAvg < 0.1 {
		rxAvg = 0
	}
	if txAvg < 0.1 {
		txAvg = 0
	}
	ev.FullText = fmt.Sprintf(`%s:%6.3g/%6.3g %s`, cfg.iface,  rxAvg/divider, txAvg/divider, unit)
	ev.ShortText = fmt.Sprintf(`-%s-`, cfg.iface)
	switch {
	case (rxBw + txBw) < 150 * 1024:
		ev.Color = "#11aaff"
	case (rxBw + txBw) < 450 * 1024:
		ev.Color = "#00ffff"
	case (rxBw + txBw) < 4 * 1024 * 1024:
		ev.Color = "#00ff00"
	case (rxBw + txBw) <  8 * 1024 * 1024:
		ev.Color = "#99ff00"
	case (rxBw + txBw) <  16 * 1024 * 1024:
		ev.Color = "#ffff00"
	default:
		ev.Color = "#ff4400"
	}
	update <- ev
}

func getStats(iface string) (uint64, uint64) {
    rawRx, _ := ioutil.ReadFile(fmt.Sprintf(`/sys/class/net/%s/statistics/rx_bytes`,iface))
	rawTx, _ := ioutil.ReadFile(fmt.Sprintf(`/sys/class/net/%s/statistics/tx_bytes`,iface))
	strRx := strings.TrimSpace(string(rawRx))
	strTx := strings.TrimSpace(string(rawTx))
	rx, _ := strconv.ParseUint(string(strRx),10,64)
	tx, _ := strconv.ParseUint(string(strTx),10,64)
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
