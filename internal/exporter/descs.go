package exporter

import "github.com/prometheus/client_golang/prometheus"

const namespace = "hy_smi"

// collectorDescs holds every prometheus.Desc the exporter emits.
// Built once at startup and reused across scrapes.
type collectorDescs struct {
	up      *prometheus.Desc
	devices *prometheus.Desc
	info    *prometheus.Desc

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

	fanRPM   *prometheus.Desc
	fanRatio *prometheus.Desc

	perfLevel  *prometheus.Desc
	boostLevel *prometheus.Desc
}

func newCollectorDescs() collectorDescs {
	gpu := []string{"gpu"}
	gpuSensor := []string{"gpu", "sensor"}
	gpuFan := []string{"gpu", "fan"}
	gpuDomain := []string{"gpu", "domain"}
	gpuLevel := []string{"gpu", "level"}

	d := func(name, help string, labels []string) *prometheus.Desc {
		return prometheus.NewDesc(name, help, labels, nil)
	}

	return collectorDescs{
		up: d(namespace+"_up",
			"Whether the most recent hy-smi collection succeeded (1 for success, 0 for failure).",
			nil),

		devices: d(namespace+"_devices",
			"Number of Hygon HCU devices discovered on this node.",
			nil),

		info: d(namespace+"_info",
			"HCU device identity (always 1). Labels carry hy-smi identity fields.",
			[]string{"gpu", "hcu", "device_id", "unique_id", "serial",
				"product_name", "part_number", "vbios", "driver",
				"bus", "health", "mode"}),

		temperature: d(namespace+"_temperature_celsius",
			"Sensor temperature in degrees Celsius.",
			gpuSensor),

		power:    d(namespace+"_power_watts", "Average package power draw in Watts.", gpu),
		powerCap: d(namespace+"_power_cap_watts", "Currently configured power cap in Watts.", gpu),
		powerMax: d(namespace+"_power_max_watts", "Maximum allowed power cap in Watts.", gpu),
		voltage: d(namespace+"_voltage_volts", "Voltage rail reading in Volts.",
			[]string{"gpu", "rail"}),

		hcuUtil:  d(namespace+"_hcu_utilization_ratio", "HCU active-time utilization ratio (0-1).", gpu),
		cuUtil:   d(namespace+"_cu_utilization_ratio", "CU (wave-front) utilization ratio (0-1).", gpu),
		waveUtil: d(namespace+"_wave_utilization_ratio", "Wave residency utilization ratio (0-1).", gpu),
		decUtil:  d(namespace+"_decoder_utilization_ratio", "Decoder utilization ratio (0-1).", gpu),
		encUtil:  d(namespace+"_encoder_utilization_ratio", "Encoder utilization ratio (0-1).", gpu),
		memUtil:  d(namespace+"_memory_utilization_ratio", "Memory controller utilization ratio (0-1).", gpu),
		vramUtil: d(namespace+"_vram_utilization_ratio", "VRAM utilization ratio (0-1).", gpu),

		vramTotal: d(namespace+"_vram_total_bytes", "Total VRAM in bytes.", gpu),
		vramUsed:  d(namespace+"_vram_used_bytes", "Used VRAM in bytes.", gpu),
		vramAvail: d(namespace+"_vram_available_bytes", "Available VRAM in bytes.", gpu),
		gttTotal:  d(namespace+"_gtt_total_bytes", "Total GTT memory in bytes.", gpu),
		gttUsed:   d(namespace+"_gtt_used_bytes", "Used GTT memory in bytes.", gpu),

		clock:    d(namespace+"_clock_hertz", "Active clock frequency in hertz.", gpuDomain),
		clockMin: d(namespace+"_clock_min_hertz", "Lowest available clock frequency in hertz.", gpuDomain),
		clockMax: d(namespace+"_clock_max_hertz", "Highest available clock frequency in hertz.", gpuDomain),

		pcieSpeed: d(namespace+"_pcie_link_speed_gts", "Current PCIe link speed in GT/s.", gpu),
		pcieWidth: d(namespace+"_pcie_link_width", "Current PCIe link width (lanes).", gpu),
		pcieBw:    d(namespace+"_pcie_bandwidth_mbytes_per_second", "PCIe bandwidth in MB/s.", gpu),

		xhclBw:     d(namespace+"_xhcl_bandwidth_mbytes_per_second", "xhcl link bandwidth in MB/s.", gpu),
		dfBw:       d(namespace+"_df_bandwidth_mbytes_per_second", "DF fabric bandwidth in MB/s.", gpu),
		umcBw:      d(namespace+"_umc_bandwidth_mbytes_per_second", "UMC memory bandwidth in MB/s.", gpu),
		xhclStatus: d(namespace+"_xhcl_status_info", "xhcl link status (info gauge, =1).", gpuLevel),

		health: d(namespace+"_health_info", "Device health status (info gauge, =1).", gpuLevel),

		rasErrors:        d(namespace+"_ras_errors_total", "Total RAS error count reported by hy-smi (counter).", gpu),
		eccCorrectable:   d(namespace+"_ecc_correctable_errors_total", "Correctable ECC error count (counter).", gpu),
		eccUncorrectable: d(namespace+"_ecc_uncorrectable_errors_total", "Uncorrectable ECC error count (counter).", gpu),
		badPagesRetired:  d(namespace+"_bad_pages_retired_total", "Retired bad pages count (counter).", gpu),
		badPagesPending:  d(namespace+"_bad_pages_pending", "Pending bad pages count.", gpu),

		fanRPM:   d(namespace+"_fan_speed_rpm", "Fan tachometer reading in RPM.", gpuFan),
		fanRatio: d(namespace+"_fan_speed_ratio", "Fan PWM duty ratio (0-1).", gpuFan),

		perfLevel:  d(namespace+"_performance_level_info", "Active DPM performance level (info gauge, =1).", gpuLevel),
		boostLevel: d(namespace+"_boost_level_info", "Active boost level (info gauge, =1).", gpuLevel),
	}
}

func (d collectorDescs) describe(ch chan<- *prometheus.Desc) {
	for _, desc := range []*prometheus.Desc{
		d.up, d.devices, d.info,
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
		d.fanRPM, d.fanRatio,
		d.perfLevel, d.boostLevel,
	} {
		ch <- desc
	}
}
