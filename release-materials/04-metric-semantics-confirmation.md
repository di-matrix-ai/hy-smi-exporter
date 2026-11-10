# 指标语义确认表

候选版本/commit：`TBD`  
海光确认人：`TBD`  
鼎茂记录人：`TBD`

## 1 通用规则

| 确认项 | 当前实现 | 海光结论 / 修订意见 |
| --- | --- | --- |
| 设备序号 | `gpu` 使用厂商工具返回或解析后的设备序号 | `TBD` |
| 实例标签 | `instance` 默认使用主机名，可由 `NODE_NAME` 覆盖 | `TBD` |
| 缺失字段 | 部分指标缺失时不输出，部分基础指标可能输出 0 | `TBD` |
| 百分比 | 使用 0-100 数值 | `TBD` |
| 温度 | 摄氏度 | `TBD` |
| 功耗 | W | `TBD` |
| 电压 | V | `TBD` |
| 容量 | bytes | `TBD` |
| 时钟 | MHz | `TBD` |
| PCIe 速率 | GT/s | `TBD` |
| 带宽 | MB/s；需确认 MB 定义和采样口径 | `TBD` |
| Info 指标 | 当前状态放在标签中，指标值恒为 1 | `TBD` |

## 2 核心指标组

| 指标组 | 指标 | 需要确认的关键语义 | 结论 |
| --- | --- | --- | --- |
| 设备 | `hy_smi_devices_total`, `hy_smi_info` | 设备数、身份字段、health/mode 枚举及敏感标签公开范围 | `TBD` |
| 温度 | `hy_smi_temperature_celsius` | edge/junction/mem/hotspot 的传感器含义 | `TBD` |
| 功耗电压 | `hy_smi_power_*`, `hy_smi_voltage_volts` | 平均功耗的采样窗口、cap/max 含义、rail 枚举 | `TBD` |
| 利用率 | `hy_smi_*_utilization_percent` | HCU/CU/Wave/编解码/内存/VRAM 的统计口径和时间窗口 | `TBD` |
| 显存 | `hy_smi_vram_*_bytes`, `hy_smi_gtt_*_bytes` | total/used/available 关系、单位换算和 GTT 含义 | `TBD` |
| 时钟 | `hy_smi_clock_*_mhz` | sclk/mclk/socclk 域、当前/最小/最大定义 | `TBD` |
| PCIe | `hy_smi_pcie_*` | link speed、width 和 bandwidth 的实时/理论口径 | `TBD` |
| 互联 | `hy_smi_xhcl_*`, `hy_smi_df_*`, `hy_smi_umc_*` | XHCL/DF/UMC 带宽和状态枚举 | `TBD` |
| 健康 | `hy_smi_health_info` | 健康状态枚举和告警建议 | `TBD` |
| RAS/ECC | `hy_smi_ras_errors_total`, `hy_smi_ecc_*_total` | 是否为启动以来累计值；RAS 是否已经包含 ECC，避免重复相加 | `TBD` |
| 坏页 | `hy_smi_bad_pages_retired_total`, `hy_smi_bad_pages_pending` | retired 是否单调累计，pending 的状态含义 | `TBD` |
| 风扇 | `hy_smi_fan_speed_*` | 风扇序号、RPM 和 PWM 适用型号 | `TBD` |
| 性能状态 | `hy_smi_performance_level_info`, `hy_smi_boost_level_info` | 枚举值及跨版本稳定性 | `TBD` |

## 3 必须重点确认的问题

1. `hy_smi_ras_errors_total` 当前由 RAS、可纠正 ECC 和不可纠正 ECC 计数相加生成。若厂商 RAS 总数已经包含 ECC，该实现会重复计数，需要在发布前修改。
2. `serial`、`unique_id`、`bus` 等标签可能属于敏感资产信息，也可能增加时序基数。请确认哪些字段允许默认公开。
3. 工具未返回部分数值时，当前代码对基础指标可能输出 0。请确认 0 表示真实值还是“未知”；若是未知，应改为不输出。
4. `MB/s` 的 MB 是十进制还是二进制口径，需要与厂商工具一致并写入 `metric.md`。

## 4 确认结果

总体结论：`通过 / 有条件通过 / 不通过 / TBD`  
需要修改的代码或文档：`TBD`  
对应 Issue / PR：`TBD`  
海光技术确认：`TBD`  
鼎茂维护者确认：`TBD`  
确认日期：`TBD`
