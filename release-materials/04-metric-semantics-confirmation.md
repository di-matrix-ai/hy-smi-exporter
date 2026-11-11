# 指标语义确认记录

适用版本：`v1.0.0`  
定稿责任方：上海鼎茂信息技术有限公司  
技术意见来源：代码实现、历史真实设备测试记录及 `hy-smi` / `hcu-smi` 命令输出  
文档状态：当前技术结论，可用于 `v1.0.0` 发布

## 1 通用规则

| 确认项 | v1.0.0 规则 | 结论 |
| --- | --- | --- |
| 设备序号 | `gpu` 使用厂商工具返回或解析后的设备序号 | 采用 |
| 实例标签 | Exporter 不生成 `instance`，由 Prometheus 根据抓取目标添加 | 采用 |
| 采集状态 | `hy_smi_up=1` 表示最近一次命令调用和解析成功，失败时为 0 | 采用 |
| 缺失字段 | 可选指标在字段不可用时可能不输出；部分基础指标按当前实现输出 0 | 采用并作为已知边界公开 |
| 利用率和 PWM | 厂商 0-100 百分值归一化为 0-1 ratio | 采用 |
| 温度 | 输出摄氏度，指标后缀为 `_celsius` | 采用 |
| 功耗 | 输出 W，指标后缀为 `_watts` | 采用 |
| 电压 | 输出 V，指标后缀为 `_volts` | 采用 |
| 容量 | 统一换算为 bytes，指标后缀为 `_bytes` | 采用 |
| 时钟 | 厂商 MHz 数值转换为 Hz，指标后缀为 `_hertz` | 采用 |
| PCIe 速率 | 保留厂商 GT/s 口径 | 采用 |
| 带宽 | 保留厂商工具提供的 MB/s 数值和口径，不转换为 bytes/s | 采用 |
| Info 指标 | 当前状态放在标签中，指标值恒为 1 | 采用 |

## 2 核心指标组结论

| 指标组 | 指标 | v1.0.0 语义结论 |
| --- | --- | --- |
| 采集状态和设备 | `hy_smi_up`, `hy_smi_devices`, `hy_smi_info` | 分别表示采集成功状态、发现的设备数量和设备身份信息 |
| 温度 | `hy_smi_temperature_celsius` | `sensor` 标签区分 edge、junction、mem 和 hotspot；实际输出取决于厂商命令字段 |
| 功耗电压 | `hy_smi_power_*`, `hy_smi_voltage_volts` | 直接采用厂商工具提供的当前值、上限值及电压值，并按指标单位输出 |
| 利用率 | `hy_smi_*_utilization_ratio` | 厂商百分值除以 100 后输出 0-1 ratio |
| 显存 | `hy_smi_vram_*_bytes`, `hy_smi_gtt_*_bytes` | total、used、available 和 GTT 字段按厂商命令含义采集并统一为 bytes |
| 时钟 | `hy_smi_clock_*_hertz` | `domain` 标签区分 sclk、mclk、socclk，当前值及最小/最大值统一为 Hz |
| PCIe | `hy_smi_pcie_*` | 输出厂商工具提供的链路速率、宽度和带宽 |
| 互联 | `hy_smi_xhcl_*`, `hy_smi_df_*`, `hy_smi_umc_*` | 输出厂商工具提供的互联带宽和状态；字段不存在时不强制生成 |
| 健康 | `hy_smi_health_info` | 厂商健康状态作为 `level` 标签，值恒为 1 |
| RAS/ECC | `hy_smi_ras_errors_total`, `hy_smi_ecc_*_total` | 按厂商累计计数字段作为 counter 输出；驱动重载或设备复位后允许计数重置 |
| 坏页 | `hy_smi_bad_pages_retired_total`, `hy_smi_bad_pages_pending` | retired 作为累计计数，pending 作为当前待处理数量 |
| 风扇 | `hy_smi_fan_speed_*` | 输出 RPM 和 PWM ratio；仅在设备及命令返回对应字段时提供 |
| 性能状态 | `hy_smi_performance_level_info`, `hy_smi_boost_level_info` | 厂商状态枚举放入 `level` 标签，值恒为 1 |

## 3 敏感信息和使用边界

- `serial`、`unique_id`、`bus` 等字段用于设备识别，但可能属于资产信息；生产环境必须限制 `/metrics` 访问，对外样例必须脱敏。
- 当 `hy_smi_up=0` 时，不应使用同一抓取周期的设备指标判断设备健康状态。
- 部分厂商字段缺失时，基础指标中的 0 可能表示真实零值，也可能表示当前命令未提供该字段；使用方应结合设备型号、驱动版本和原始命令输出判断。
- 设备、驱动、固件、DTK 或厂商工具升级后，如输出字段、单位或枚举变化，需要补充解析适配和回归测试。
- Exporter 面向通用 `hy-smi` / `hcu-smi` 命令接口，不针对每个硬件型号重复实现采集逻辑。

## 4 发布结论

总体结论：`v1.0.0` 指标名称、类型、单位和标签规则由上海鼎茂信息技术有限公司定稿，可用于社区发布。  
兼容结论：驱动正常支持目标设备且 `hy-smi` / `hcu-smi` 命令输出保持兼容时，当前指标实现适用于该设备。  
真实设备结论：真实设备测试状态在 `docs/compatibility.md` 中单独记录，不与通用接口支持范围混淆。  
后续修订：发现厂商字段语义或输出格式变化时，通过 Issue、Pull Request、解析样例和版本变更记录进行修订。
