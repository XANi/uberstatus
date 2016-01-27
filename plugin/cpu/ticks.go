package cpu

type cpuTicks struct {
	// non-system-specific
	total uint64 // generic used stat, in case that gets ever posrted for non-linux systems
	idle uint64 // generic idle stat
	// Linux-specific
	system uint64
	user uint64
	nice uint64
	iowait uint64 // 2.5.41
	irq uint64 // 2.6
	softirq uint64 // 2.6
	steal uint64 // 2.6.11
	guest uint64 // 2.6.24
	guestNice uint64 // 2.6.33
}

func (c cpuTicks) GetCpuUsagePercent () float64 {
	if (c.total == 0) { return 0 } // empty object
	return ( float64(c.idle) / float64(c.total) ) * 100

}

func (c cpuTicks) Sub(c2 cpuTicks) cpuTicks {
	var out cpuTicks
	out.total = c.total - c2.total
	out.idle = c.idle - c2.idle
	out.system = c.system - c2.system
	out.user = c.user - c2.user
	out.nice = c.nice - c2.nice
	out.iowait = c.iowait - c2.iowait
	out.irq = c.irq - c2.irq
	out.softirq = c.softirq - c2.softirq
	out.steal = c.steal - c2.steal
	out.guest = c.guest - c2.guest
	out.guestNice = c.guestNice - c2.guestNice
	return out
}
