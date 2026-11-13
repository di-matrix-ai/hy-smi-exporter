# hy-smi-exporter v1.0.0 发布说明

本文档用于说明 `hy-smi-exporter` 的项目作用、服务工作方式、发布范围和 Prometheus 指标，可作为 `v1.0.0` 的对外发布文档及指标字典。

## 1 发布信息

| 项目 | 内容 |
| --- | --- |
| 项目名称 | `hy-smi-exporter` |
| 发布版本 | `v1.0.0` |
| 项目作用 | 将海光 DCU / HCU 设备运行数据转换为 Prometheus 标准指标 |
| 发布主体 | 上海鼎茂信息技术有限公司（Shanghai Ding Mao Information Technology Co., Ltd.） |
| 开源协议 | Apache License 2.0 |
| 代码仓库 | [github.com/di-matrix-ai/hy-smi-exporter](https://github.com/di-matrix-ai/hy-smi-exporter) |
| 默认监听地址 | `:9584` |
| 指标端点 | `/metrics` |
| 服务健康端点 | `/healthz` |
| 支持范围 | 以 [`docs/compatibility.md`](docs/compatibility.md) 为准 |

## 2 hy-smi-exporter 的作用

`hy-smi-exporter` 是连接海光 DCU / HCU 设备管理工具与 Prometheus 监控体系的指标采集服务。它调用宿主机已有的 `hy-smi` 只读命令，解析设备状态和运行数据，统一转换为以 `hy_smi_` 为前缀的 Prometheus 指标，并通过 HTTP `/metrics` 端点提供给 Prometheus 抓取。

它主要解决以下问题：

| 场景 | 作用 |
| --- | --- |
| 设备统一纳管 | 将不同节点上的 DCU / HCU 设备接入统一的 Prometheus 监控体系 |
| 运行状态监控 | 持续采集设备数量、温度、功耗、利用率、显存和时钟等运行数据 |
| 健康与故障监控 | 输出设备健康状态、链路状态、RAS、ECC 和坏页等故障相关指标 |
| 容量与性能观察 | 支持分析显存占用、设备负载、PCIe 及设备互联带宽情况 |
| 告警和可视化 | 指标可用于 Prometheus 告警规则以及 Grafana 等可视化平台 |
| 运维问题定位 | 通过设备、驱动、总线、健康状态和采集状态指标辅助定位故障 |

`hy-smi-exporter` 只负责采集和转换指标，不安装或分发驱动、固件及 `hy-smi` 厂商工具，也不修改设备配置。

## 3 服务工作方式

Prometheus 每次访问 `/metrics` 时，Exporter 按以下流程执行采集：

| 步骤 | 处理内容 |
| --- | --- |
| 1 | 默认执行轻量的 `hy-smi --json`，获取结构化设备基础数据 |
| 2 | JSON 参数不受支持或解析失败时，自动回退到 `hy-smi` 文本输出 |
| 3 | 主查询缺少唯一标识或健康状态时，按需执行 `--showuniqueid`、`--healthcheck` 补充查询 |
| 4 | 主查询缺少显存信息时，按需执行只读的显存补充查询 |
| 5 | 仅在显式启用时执行 `hy-smi -a`，采集过去 1 秒的利用率和带宽统计 |
| 6 | 解析每张设备卡的数据，并将单位统一转换为 Prometheus 推荐的基础单位 |
| 7 | 通过 `/metrics` 输出设备级指标，供 Prometheus 定时抓取 |
| 8 | 通过 `hy_smi_up` 标识最近一次必需主查询是否成功；可选补充查询失败不会把设备整体误报为不可用 |

显存补充查询包括：

```text
hy-smi --showmeminfo all --json
hy-smi --showmemavailable --json
```

服务还提供 `/healthz` 用于检查 Exporter HTTP 服务是否存活。`/healthz` 返回成功只代表服务进程可以响应；设备命令是否采集成功，应以 `hy_smi_up` 指标为准。

部分利用率和带宽参数统计过去 1 秒的数据，8 卡等多卡机型执行耗时会随设备数量增加。`v1.0.0` 默认关闭这些统计项，启用方式如下：

```bash
./hy-smi-exporter --hy-smi.collect-1s-stats=true
```

也可设置 `HY_SMI_EXPORTER_COLLECT_1S_STATS=true`。未启用时，对应的利用率以及 PCIe、XHCL、DF、UMC 带宽指标不会输出，避免用数值 0 误表示“已采集但无负载”。

## 4 v1.0.0 发布范围

| 项目 | v1.0.0 范围 |
| --- | --- |
| 运行方式 | 支持直接运行和 systemd 服务部署 |
| Prometheus 接入 | 支持通过 `/metrics` 标准端点抓取 |
| 容器场景 | 暂不兼容、不测试、不发布 Docker、Docker Compose、Kubernetes 和容器镜像 |
| 设备和驱动范围 | 以公开兼容性矩阵为准 |
| 厂商组件 | 不随项目分发驱动、固件或 `hy-smi` 二进制 |
| 操作安全 | 仅调用查询命令，不执行设备配置修改、驱动加载、复位或故障注入操作 |

实际兼容性取决于硬件型号、驱动、固件、DTK、`hy-smi` 版本、Linux 发行版和内核版本的组合。详细信息见 [`docs/compatibility.md`](docs/compatibility.md)。

## 5 Prometheus 指标

所有设备指标以 `hy_smi_` 为前缀。设备级指标携带 `gpu` 标签，用于区分同一节点上的不同设备；`instance` 标签由 Prometheus 抓取目标自动添加。

下表中利用率指标及 PCIe、XHCL、DF、UMC 带宽指标属于可选 1 秒统计项，默认不输出。

| 指标名 | 类型 | 说明 |
| --- | --- | --- |
| `hy_smi_up` | gauge | 最近一次 `hy-smi` 采集是否成功，成功为 1，失败为 0 |
| `hy_smi_devices` | gauge | 最近一次成功采集发现的 HCU 设备总数 |
| `hy_smi_info` | gauge (=1) | 设备身份信息，包括设备 ID、序列号、产品名称、VBIOS、驱动、总线、健康状态和工作模式等 |
| `hy_smi_temperature_celsius` | gauge | 各传感器温度，包括 edge、junction、mem 和 hotspot，单位为摄氏度 |
| `hy_smi_power_watts` | gauge | 平均功耗，单位为瓦 |
| `hy_smi_power_cap_watts` | gauge | 当前功耗上限，单位为瓦 |
| `hy_smi_power_max_watts` | gauge | 最大允许功耗，单位为瓦 |
| `hy_smi_voltage_volts` | gauge | 核心电压，单位为伏 |
| `hy_smi_hcu_utilization_ratio` | gauge | HCU 活跃时间比率，取值范围为 0-1 |
| `hy_smi_cu_utilization_ratio` | gauge | CU 波前利用率，取值范围为 0-1 |
| `hy_smi_wave_utilization_ratio` | gauge | 波前驻留率，取值范围为 0-1 |
| `hy_smi_decoder_utilization_ratio` | gauge | 解码器利用率，取值范围为 0-1 |
| `hy_smi_encoder_utilization_ratio` | gauge | 编码器利用率，取值范围为 0-1 |
| `hy_smi_memory_utilization_ratio` | gauge | 内存控制器利用率，取值范围为 0-1 |
| `hy_smi_vram_utilization_ratio` | gauge | 显存利用率，取值范围为 0-1 |
| `hy_smi_vram_total_bytes` | gauge | 显存总量，单位为字节 |
| `hy_smi_vram_used_bytes` | gauge | 已用显存，单位为字节 |
| `hy_smi_vram_available_bytes` | gauge | 可用显存，单位为字节 |
| `hy_smi_gtt_total_bytes` | gauge | GTT 内存总量，单位为字节 |
| `hy_smi_gtt_used_bytes` | gauge | 已用 GTT 内存，单位为字节 |
| `hy_smi_clock_hertz` | gauge | 当前时钟频率，包括 sclk、mclk 和 socclk，单位为赫兹 |
| `hy_smi_clock_min_hertz` | gauge | 最低可用时钟频率，单位为赫兹 |
| `hy_smi_clock_max_hertz` | gauge | 最高可用时钟频率，单位为赫兹 |
| `hy_smi_pcie_link_speed_gts` | gauge | PCIe 链路速率，单位为 GT/s |
| `hy_smi_pcie_link_width` | gauge | PCIe 链路宽度，数值为通道数 |
| `hy_smi_pcie_bandwidth_mbytes_per_second` | gauge | PCIe 带宽，单位为 MB/s |
| `hy_smi_xhcl_bandwidth_mbytes_per_second` | gauge | XHCL 互联带宽，单位为 MB/s |
| `hy_smi_df_bandwidth_mbytes_per_second` | gauge | DF Fabric 带宽，单位为 MB/s |
| `hy_smi_umc_bandwidth_mbytes_per_second` | gauge | UMC 内存带宽，单位为 MB/s |
| `hy_smi_xhcl_status_info` | gauge (=1) | XHCL 链路状态信息 |
| `hy_smi_health_info` | gauge (=1) | 设备健康状态信息 |
| `hy_smi_ras_errors_total` | counter | RAS 错误累计数量 |
| `hy_smi_ecc_correctable_errors_total` | counter | 可纠正 ECC 错误累计数量 |
| `hy_smi_ecc_uncorrectable_errors_total` | counter | 不可纠正 ECC 错误累计数量 |
| `hy_smi_bad_pages_retired_total` | counter | 已退役坏页累计数量 |
| `hy_smi_bad_pages_pending` | gauge | 待处理坏页数量 |
| `hy_smi_fan_speed_rpm` | gauge | 风扇转速，单位为 RPM |
| `hy_smi_fan_speed_ratio` | gauge | 风扇 PWM 占空比，取值范围为 0-1 |
| `hy_smi_performance_level_info` | gauge (=1) | DPM 性能等级信息 |
| `hy_smi_boost_level_info` | gauge (=1) | Boost 等级信息 |

Exporter 同时注册 Go 运行时和进程指标，用于观察 Exporter 自身的 CPU、内存、文件描述符和 Go 运行状态。

## 6 Prometheus 接入示例

```yaml
scrape_configs:
  - job_name: hy-smi-exporter
    static_configs:
      - targets:
          - '<节点IP>:9584'
```

接入后可使用以下查询检查采集情况：

| 检查目标 | PromQL 示例 |
| --- | --- |
| Exporter 采集是否成功 | `hy_smi_up` |
| 节点设备数量 | `hy_smi_devices` |
| 各设备温度 | `hy_smi_temperature_celsius` |
| 各设备 HCU 利用率 | `hy_smi_hcu_utilization_ratio` |
| 各设备已用显存 | `hy_smi_vram_used_bytes` |
| 不可纠正 ECC 错误增量 | `increase(hy_smi_ecc_uncorrectable_errors_total[5m])` |

具体告警阈值应根据设备型号、业务负载、厂商建议和实际运行基线确定，不应对所有型号使用未经验证的统一阈值。

## 7 安全与数据边界

- Exporter 只调用只读查询命令，不调用 `--set*`、`--rasinject`、`--load*driver`、`--reset*` 等修改设备状态的选项。
- `/metrics` 可能包含设备序列号、Unique ID、产品信息和 PCIe 地址等标签，生产环境应限制访问并在对外提供样例前完成脱敏。
- `hy_smi_up` 为 0 时，表示最近一次厂商工具调用或解析失败，应结合 Exporter 日志、驱动状态、工具路径和运行权限排查。
- 驱动、固件、DTK 或厂商工具升级后，应重新执行兼容性回归。

## 8 相关文档

| 文档 | 用途 |
| --- | --- |
| [`README.md`](README.md) | 编译、运行、参数、部署及故障排查 |
| [`docs/installation.md`](docs/installation.md) | 安装和部署说明 |
| [`docs/compatibility.md`](docs/compatibility.md) | 硬件、驱动和真实设备验证范围 |
| [`docs/RELEASE_NOTES_v1.0.0.md`](docs/RELEASE_NOTES_v1.0.0.md) | `v1.0.0` Release Notes |
| [`docs/COMMUNITY_ANNOUNCEMENT_v1.0.0.md`](docs/COMMUNITY_ANNOUNCEMENT_v1.0.0.md) | 社区发布公告 |
