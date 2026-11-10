# hy-smi-exporter

面向海光 DCU / HCU 加速器的 Prometheus 指标采集器，实际支持范围以 [`docs/compatibility.md`](docs/compatibility.md) 中完成真实设备验证的组合为准。
通过只读命令 `hy-smi -a --json` 采集每张卡的运行数据；对于不支持该 JSON 组合参数的驱动版本，自动回退到 `hy-smi -a` 文本输出，以 `hy_smi_` 前缀暴露在 `/metrics` 端点。


## 采集指标

所有指标以 `hy_smi_` 为前缀。设备级指标携带 `gpu`（设备序号）标签；
`instance` 由 Prometheus 抓取目标自动添加。

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `hy_smi_up` | gauge | 最近一次厂商工具采集是否成功（成功为 1，失败为 0） |
| `hy_smi_devices` | gauge | 最近一次成功采集发现的 HCU 设备总数 |
| `hy_smi_info` | gauge (=1) | 设备身份信息（device_id / serial / product_name / vbios / driver / bus / health / mode 等） |
| `hy_smi_temperature_celsius` | gauge | 各传感器温度（edge / junction / mem / hotspot） |
| `hy_smi_power_watts` | gauge | 平均功耗（W） |
| `hy_smi_power_cap_watts` | gauge | 当前功耗上限（W） |
| `hy_smi_power_max_watts` | gauge | 最大允许功耗（W） |
| `hy_smi_voltage_volts` | gauge | 核心电压（V） |
| `hy_smi_hcu_utilization_ratio` | gauge | HCU 活跃时间比率（0-1） |
| `hy_smi_cu_utilization_ratio` | gauge | CU 波前利用率 |
| `hy_smi_wave_utilization_ratio` | gauge | 波前驻留率 |
| `hy_smi_decoder_utilization_ratio` | gauge | 解码器利用率 |
| `hy_smi_encoder_utilization_ratio` | gauge | 编码器利用率 |
| `hy_smi_memory_utilization_ratio` | gauge | 内存控制器利用率 |
| `hy_smi_vram_utilization_ratio` | gauge | 显存利用率 |
| `hy_smi_vram_total_bytes` | gauge | 显存总量（字节） |
| `hy_smi_vram_used_bytes` | gauge | 已用显存（字节） |
| `hy_smi_vram_available_bytes` | gauge | 可用显存（字节） |
| `hy_smi_gtt_total_bytes` | gauge | GTT 内存总量（字节） |
| `hy_smi_gtt_used_bytes` | gauge | 已用 GTT（字节） |
| `hy_smi_clock_hertz` | gauge | 当前时钟频率（sclk / mclk / socclk，Hz） |
| `hy_smi_clock_min_hertz` | gauge | 最低可用时钟频率（Hz） |
| `hy_smi_clock_max_hertz` | gauge | 最高可用时钟频率（Hz） |
| `hy_smi_pcie_link_speed_gts` | gauge | PCIe 链路速率（GT/s） |
| `hy_smi_pcie_link_width` | gauge | PCIe 链路宽度（通道数） |
| `hy_smi_pcie_bandwidth_mbytes_per_second` | gauge | PCIe 带宽（MB/s） |
| `hy_smi_xhcl_bandwidth_mbytes_per_second` | gauge | xhcl 互联带宽（MB/s） |
| `hy_smi_df_bandwidth_mbytes_per_second` | gauge | DF Fabric 带宽（MB/s） |
| `hy_smi_umc_bandwidth_mbytes_per_second` | gauge | UMC 内存带宽（MB/s） |
| `hy_smi_xhcl_status_info` | gauge (=1) | xhcl 链路状态 |
| `hy_smi_health_info` | gauge (=1) | 设备健康状态 |
| `hy_smi_ras_errors_total` | counter | RAS 错误总数 |
| `hy_smi_ecc_correctable_errors_total` | counter | 可纠正 ECC 错误数 |
| `hy_smi_ecc_uncorrectable_errors_total` | counter | 不可纠正 ECC 错误数 |
| `hy_smi_bad_pages_retired_total` | counter | 已退役坏页数 |
| `hy_smi_bad_pages_pending` | gauge | 待处理坏页数 |
| `hy_smi_fan_speed_rpm` | gauge | 风扇转速（RPM） |
| `hy_smi_fan_speed_ratio` | gauge | 风扇 PWM 占空比（0-1） |
| `hy_smi_performance_level_info` | gauge (=1) | DPM 性能等级 |
| `hy_smi_boost_level_info` | gauge (=1) | Boost 等级 |
