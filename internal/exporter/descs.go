package exporter

import "github.com/prometheus/client_golang/prometheus"

const namespace = "hy_smi"

// collectorDescs holds every prometheus.Desc the exporter emits.
// Built once at startup and reused across scrapes.
type collectorDescs struct {
	devicesTotal *prometheus.Desc
	info         *prometheus.Desc

	temperature *prometheus.Desc

	power    *prometheus.Desc
	powerCap *prometheus.Desc
	powerMax *prometheus.Desc
	voltage  *prometheus.Desc

	hcuUtil  *prometheus.Desc
	cuUtil   *prometheus.Desc
	waveUtil *prometheus.Desc
	decUtil  *prometheus.Desc
	encUtil  *prometheus.Desc
	memUtil  *prometheus.Desc
	vramUtil *prometheus.Desc

	vramTotal *prometheus.Desc
	vramUsed  *prometheus.Desc
	vramAvail *prometheus.Desc
	gttTotal  *prometheus.Desc
	gttUsed   *prometheus.Desc

	clock    *prometheus.Desc
	clockMin *prometheus.Desc
	clockMax *prometheus.Desc

	pcieSpeed *prometheus.Desc
	pcieWidth *prometheus.Desc
	pcieBw    *prometheus.Desc

	xhclBw     *prometheus.Desc
	dfBw       *prometheus.Desc
	umcBw      *prometheus.Desc
	xhclStatus *prometheus.Desc

	health           *prometheus.Desc
	rasErrors        *prometheus.Desc
	eccCorrectable   *prometheus.Desc
	eccUncorrectable *prometheus.Desc
	badPagesRetired  *prometheus.Desc
	badPagesPending  *prometheus.Desc

	fanRPM    *prometheus.Desc
	fanPWMPct *prometheus.Desc

	perfLevel  *prometheus.Desc
	boostLevel *prometheus.Desc
}

func newCollectorDescs() collectorDescs {
	gpuInst := []string{"gpu", "instance"}
	gpuSensorInst := []string{"gpu", "sensor", "instance"}
	gpuFanInst := []string{"gpu", "fan", "instance"}
	gpuDomainInst := []string{"gpu", "domain", "instance"}
	gpuLevelInst := []string{"gpu", "level", "instance"}

	d := func(name, help string, labels []string) *prometheus.Desc {
		return prometheus.NewDesc(name, help, labels, nil)
	}

	return collectorDescs{
		devicesTotal: d(namespace+"_devices_total",
			"Number of Hygon HCU devices discovered on this node.",
			[]string{"instance"}),

		info: d(namespace+"_info",
			"HCU device identity (always 1). Labels carry hy-smi identity fields.",
			[]string{"gpu", "instance", "hcu", "device_id", "unique_id", "serial",
				"product_name", "part_number", "vbios", "driver",
				"bus", "health", "mode"}),

		temperature: d(namespace+"_temperature_celsius",
			"Sensor temperature in degrees Celsius.",
			gpuSensorInst),

		power:    d(namespace+"_power_watts", "Average package power draw in Watts.", gpuInst),
		powerCap: d(namespace+"_power_cap_watts", "Currently configured power cap in Watts.", gpuInst),
		powerMax: d(namespace+"_power_max_watts", "Maximum allowed power cap in Watts.", gpuInst),
		voltage: d(namespace+"_voltage_volts", "Voltage rail reading in Volts.",
			[]string{"gpu", "rail", "instance"}),

		hcuUtil:  d(namespace+"_hcu_utilization_percent", "HCU active-time utilization (0-100).", gpuInst),
		cuUtil:   d(namespace+"_cu_utilization_percent", "CU (wave-front) utilization (0-100).", gpuInst),
		waveUtil: d(namespace+"_wave_utilization_percent", "Wave residency utilization (0-100).", gpuInst),
		decUtil:  d(namespace+"_decoder_utilization_percent", "Decoder utilization (0-100).", gpuInst),
		encUtil:  d(namespace+"_encoder_utilization_percent", "Encoder utilization (0-100).", gpuInst),
		memUtil:  d(namespace+"_memory_utilization_percent", "Memory controller utilization (0-100).", gpuInst),
		vramUtil: d(namespace+"_vram_utilization_percent", "VRAM utilization (0-100).", gpuInst),

		vramTotal: d(namespace+"_vram_total_bytes", "Total VRAM in bytes.", gpuInst),
		vramUsed:  d(namespace+"_vram_used_bytes", "Used VRAM in bytes.", gpuInst),
		vramAvail: d(namespace+"_vram_available_bytes", "Available VRAM in bytes.", gpuInst),
		gttTotal:  d(namespace+"_gtt_total_bytes", "Total GTT memory in bytes.", gpuInst),
		gttUsed:   d(namespace+"_gtt_used_bytes", "Used GTT memory in bytes.", gpuInst),

		clock:    d(namespace+"_clock_mhz", "Active clock frequency in MHz.", gpuDomainInst),
		clockMin: d(namespace+"_clock_min_mhz", "Lowest available clock frequency in MHz.", gpuDomainInst),
		clockMax: d(namespace+"_clock_max_mhz", "Highest available clock frequency in MHz.", gpuDomainInst),

		pcieSpeed: d(namespace+"_pcie_link_speed_gts", "Current PCIe link speed in GT/s.", gpuInst),
		pcieWidth: d(namespace+"_pcie_link_width", "Current PCIe link width (lanes).", gpuInst),
		pcieBw:    d(namespace+"_pcie_bandwidth_mbytes_per_second", "PCIe bandwidth in MB/s.", gpuInst),

		xhclBw:     d(namespace+"_xhcl_bandwidth_mbytes_per_second", "xhcl link bandwidth in MB/s.", gpuInst),
		dfBw:       d(namespace+"_df_bandwidth_mbytes_per_second", "DF fabric bandwidth in MB/s.", gpuInst),
		umcBw:      d(namespace+"_umc_bandwidth_mbytes_per_second", "UMC memory bandwidth in MB/s.", gpuInst),
		xhclStatus: d(namespace+"_xhcl_status_info", "xhcl link status (info gauge, =1).", gpuLevelInst),

		health: d(namespace+"_health_info", "Device health status (info gauge, =1).", gpuLevelInst),

		rasErrors:        d(namespace+"_ras_errors_total", "Total RAS error count (counter).", gpuInst),
		eccCorrectable:   d(namespace+"_ecc_correctable_errors_total", "Correctable ECC error count (counter).", gpuInst),
		eccUncorrectable: d(namespace+"_ecc_uncorrectable_errors_total", "Uncorrectable ECC error count (counter).", gpuInst),
		badPagesRetired:  d(namespace+"_bad_pages_retired_total", "Retired bad pages count (counter).", gpuInst),
		badPagesPending:  d(namespace+"_bad_pages_pending", "Pending bad pages count.", gpuInst),

		fanRPM:    d(namespace+"_fan_speed_rpm", "Fan tachometer reading in RPM.", gpuFanInst),
		fanPWMPct: d(namespace+"_fan_speed_percent", "Fan PWM duty (0-100).", gpuFanInst),

		perfLevel:  d(namespace+"_performance_level_info", "Active DPM performance level (info gauge, =1).", gpuLevelInst),
		boostLevel: d(namespace+"_boost_level_info", "Active boost level (info gauge, =1).", gpuLevelInst),
	}
}

func (d collectorDescs) describe(ch chan<- *prometheus.Desc) {
	for _, desc := range []*prometheus.Desc{
		d.devicesTotal, d.info,
		d.temperature,
		d.power, d.powerCap, d.powerMax, d.voltage,
		d.hcuUtil, d.cuUtil, d.waveUtil, d.decUtil, d.encUtil, d.memUtil, d.vramUtil,
		d.vramTotal, d.vramUsed, d.vramAvail, d.gttTotal, d.gttUsed,
		d.clock, d.clockMin, d.clockMax,
		d.pcieSpeed, d.pcieWidth, d.pcieBw,
		d.xhclBw, d.dfBw, d.umcBw, d.xhclStatus,
		d.health,
		d.rasErrors, d.eccCorrectable, d.eccUncorrectable,
		d.badPagesRetired, d.badPagesPending,
		d.fanRPM, d.fanPWMPct,
		d.perfLevel, d.boostLevel,
	} {
		ch <- desc
	}
}
