package network

import (
	"fmt"
	"github.com/VividCortex/ewma"
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	"github.com/op/go-logging"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"
)

var log = logging.MustGetLogger("main")

type Config struct {
	iface string
	interval int
}

type netStats struct {
	ip     string
	tx     uint64
	rx     uint64
	oldTx  uint64
	oldRx  uint64
	ewmaTx ewma.MovingAverage
	ewmaRx ewma.MovingAverage
	oldTs  time.Time
	ts     time.Time
	nextTs time.Time
	iface string
	interval int
}

const ShowFirstAddr = 0
const ShowSecondAddr = 1
const ShowAllAddr = -1


func New(c uber.PluginConfig) (uber.Plugin, error) {
	stats := &netStats{}
	cfg := loadConfig(c.Config)
	stats.ewmaRx = ewma.NewMovingAverage(5)
	stats.ewmaTx = ewma.NewMovingAverage(5)
	stats.oldTs = time.Now()
	stats.ts = time.Now()
	stats.iface = cfg.iface
	stats.interval = cfg.interval

	return  stats, nil
}

func (s *netStats) Init() error {return nil}

func (s *netStats) GetUpdateInterval() int {
	return s.interval
}

func (s *netStats) UpdatePeriodic() uber.Update {
	ev, _ := s.Update()
	return ev
}

func (n *netStats) UpdateFromEvent(ev uber.Event) uber.Update {
	n.nextTs = time.Now().Add(time.Second * 3)
	if ev.Button == 1 {
		return n.UpdateAddr(ShowFirstAddr)
	} else if ev.Button == 3 {
		return n.UpdateAddr(ShowSecondAddr)
	} else {
		return n.UpdateAddr(ShowAllAddr)
	}
}

func loadConfig(raw map[string]interface{}) Config {
	var c Config
	c.iface = `lo`
	c.interval = 1000
	for key, value := range raw {
		converted, ok := value.(string)
		if ok {
			switch {
			case key == `iface`:
				c.iface = converted
				log.Warningf("-- %s %s--", key, c.iface)

			}
	} else {
			converted, ok := value.(int)
			if ok {
				switch {
				case key == `interval`:
					c.interval = converted
				default:
					log.Warningf("unknown config key: [%s]", key)
				}
			} else {
				log.Errorf("Cant interpret value of config key [%s]", key)
			}
		}
	}
	return c
}

func (stats *netStats) UpdateAddr(addr_id int) (ev uber.Update) {
	ev.Color = `#aaffaa`
	ifaces, _ := net.Interfaces()
end:
	for _, iface := range ifaces {
		if iface.Name == stats.iface {
			v, _ := iface.Addrs()
			if len(v) <= addr_id {
				break end
			}
			if addr_id < 0 {
				ev.FullText = fmt.Sprintf("%+v", v)
			} else {
				ev.FullText = fmt.Sprintf("%+v", v[addr_id])
			}
			return ev
		}
	}

	ev.FullText = fmt.Sprintf("%s??", stats.iface)
	return ev
}

func (stats *netStats) Update() (ev uber.Update, ok bool) {
	util.WaitForTs(&stats.nextTs)
	ev.Color = `#666666`
	ev.Markup = `pango`
	ev.FullText = fmt.Sprintf(`<span color="#666666">%s!!</span>`, stats.iface)
	stats.oldTs = stats.ts
	rx, tx := getStats(stats.iface)
	stats.ts = time.Now()

	// TODO: do same on bigger time diff
	if stats.ts.UnixNano() < stats.oldTs.UnixNano() {
		// we are in time machine.. or ntp changed clock
		stats.oldTs = stats.ts
		return ev, false
	}

	// either interface never seen packets, or it got recreated, reset it
	if rx == 0 && tx == 0 {
		stats.rx = 0
		stats.tx = 0
		stats.oldRx = 0
		stats.oldTx = 0
		return ev, true
	}

	// counter flipped, or interface recreated, reset to current value
	if stats.rx > rx || stats.tx > tx {
		stats.rx = rx
		stats.tx = tx
		stats.oldRx = rx
		stats.oldTx = tx
		return ev, false
	}
	//  init on first probe on empty interface
	if stats.rx == 0 && rx > 0 {
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
	if tsDiff < 0.01 {
		return ev,false // quicker probing doesnt make sense, no div by 0, should probably return an error...
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

	ev.FullText = fmt.Sprintf(`<span color="#aaffaa">%s</span>:<span color="%s">%6.3g</span>/<span color="%s">%6.3g</span><span color="%s"> %s</span>`,
		stats.iface,
		getBwColor(rxAvg),
		rxAvg/divider,
		getBwColor(txAvg),
		txAvg/divider,
		getBwColor(txAvg+rxAvg),
		unit,
	)
	ev.ShortText = fmt.Sprintf(`-%s-`, stats.iface)
	return ev,true
}

func getStats(iface string) (uint64, uint64) {
	rawRx, _ := ioutil.ReadFile(fmt.Sprintf(`/sys/class/net/%s/statistics/rx_bytes`, iface))
	rawTx, _ := ioutil.ReadFile(fmt.Sprintf(`/sys/class/net/%s/statistics/tx_bytes`, iface))
	strRx := strings.TrimSpace(string(rawRx))
	strTx := strings.TrimSpace(string(rawTx))
	rx, _ := strconv.ParseUint(string(strRx), 10, 64)
	tx, _ := strconv.ParseUint(string(strTx), 10, 64)
	return rx, tx
}

func getBwColor(bw float64) string {
	switch {
	case bw < 50*1024:
		return "#666666"
	case bw < 150*1024:
		return "#11aaff"
	case bw < 450*1024:
		return "#00ffff"
	case bw < 4*1024*1024:
		return "#00ff00"
	case bw < 8*1024*1024:
		return "#99ff00"
	case bw < 16*1024*1024:
		return "#ffff00"
	default:
		return "#ff4400"
	}
}

func getUnit(bytes float64) (divider float64, unit string) {
	switch {
	case bytes < 125*1024:
		return 1024 / 8, `Kb`
	case bytes < 100*1024*1024:
		return 1024 * 1024 / 8, `Mb`
	default:
		return 1024 * 1024 * 1024 / 8, `Gb`
	}
}
