# hy-smi-exporter

面向海光 DCU / HCU 加速器的 Prometheus 指标采集器。项目基于驱动提供的通用 `hy-smi` 命令接口适配设备；命令接口保持兼容时可共用当前采集实现，具体支持范围和真实设备验证状态见 [兼容性矩阵](docs/compatibility.md)。
默认通过轻量只读命令 `hy-smi --json` 采集每张卡的基础运行数据；JSON 不可用时自动回退到 `hy-smi` 文本输出。设备唯一标识和健康状态分别通过 `--showuniqueid`、`--healthcheck` 补充采集，并以 `hy_smi_` 前缀暴露在 `/metrics` 端点。

参考 [reloaded/prometheus_consumer_amdgpu_exporter](https://github.com/reloaded/prometheus_consumer_amdgpu_exporter)（AMD）和 [NVIDIA/dcgm-exporter](https://github.com/NVIDIA/dcgm-exporter)（NVIDIA）的架构，采用"每次抓取执行主查询 + prometheus.Collector"模式；当主查询不含显存统计时，按需补充只读显存查询。

English overview: [README.en.md](README.en.md). Community contributions, security
reports, and release history are covered in [CONTRIBUTING.md](CONTRIBUTING.md),
[SECURITY.md](SECURITY.md), and [CHANGELOG.md](CHANGELOG.md). Release steps are
documented in [RELEASING.md](RELEASING.md). Driver and release documentation is
available in [docs/compatibility.md](docs/compatibility.md) and
[docs/installation.md](docs/installation.md).

项目架构、支持边界、第三方依赖和正式发布表单分别见
[ARCHITECTURE.md](ARCHITECTURE.md)、[SUPPORT.md](SUPPORT.md)、
[THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md) 和
[`release-materials/`](release-materials/00-发布材料总目录.md)。后续方向见
[ROADMAP.md](ROADMAP.md)。

---

## 安全说明

**本程序只调用只读命令，绝不执行任何写入操作。** 默认每次抓取按需执行：

```
hy-smi --json
hy-smi                            # 仅当 JSON 不支持或解析失败时
hy-smi --showuniqueid --json      # 主查询缺少 Unique ID 时
hy-smi --healthcheck --json       # 主查询缺少健康状态时
hy-smi --showmeminfo all --json       # 仅当主查询缺少显存统计时
hy-smi --showmemavailable --json      # 仅当主查询缺少显存统计时
```

部分利用率和带宽参数会统计过去 1 秒的数据，在 8 卡等多卡节点上可能增加抓取耗时，因此默认不采集。需要这些指标时，显式设置 `--hy-smi.collect-1s-stats=true`；启用后主查询使用 `hy-smi -a --json`，必要时回退到 `hy-smi -a`。

不会调用 `--set*`、`--rasinject`、`--load*driver`、`--reset*` 等任何修改硬件配置的选项。

---

## 采集指标

所有指标以 `hy_smi_` 为前缀。设备级指标携带 `gpu`（设备序号）标签；
`instance` 由 Prometheus 抓取目标自动添加，exporter 不重复生成。

利用率及 PCIe、XHCL、DF、UMC 带宽指标属于 1 秒统计项，默认不输出；只有启用 `--hy-smi.collect-1s-stats=true` 后才会出现。

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `hy_smi_up` | gauge | 最近一次 `hy-smi` 采集是否成功（成功为 1，失败为 0） |
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

---

## 环境要求

- Go 1.23+（编译时需要）
- `hy-smi` 命令可执行（在 `$PATH` 中或通过 `--hy-smi.path` 指定）
- 目标设备应由对应驱动正常支持，且 `hy-smi` 命令输出与当前解析实现兼容；具体支持范围和实测状态见 [兼容性矩阵](docs/compatibility.md)
- 首版支持直接运行或 systemd；容器、Docker Compose 和 Kubernetes 不属于首版兼容范围

---

## 编译

```bash
cd hy-smi-exporter

# 编译二进制（输出到 bin/hy-smi-exporter）
make build

# 运行单元测试
make test

# 运行格式检查
make fmt
```

纯 Go 编译，无 cgo，无额外系统依赖。

---

## 运行

### 直接运行

```bash
./bin/hy-smi-exporter
```

默认监听 `:9584`，访问 `http://<HOST>:9584/metrics` 查看指标。

### 常用参数

```bash
./bin/hy-smi-exporter \
  --web.listen-address=:9584 \
  --hy-smi.path=/usr/bin/hy-smi \
  --hy-smi.timeout=5s \
  --hy-smi.collect-1s-stats=false \
  --log.level=info
```

### 全部参数说明

| 参数 | 默认值 | 环境变量 | 说明 |
|------|--------|----------|------|
| `--web.listen-address` | `:9584` | `HY_SMI_EXPORTER_LISTEN_ADDR` | 监听地址 |
| `--web.telemetry-path` | `/metrics` | `HY_SMI_EXPORTER_METRICS_PATH` | 指标路径 |
| `--log.level` | `info` | `HY_SMI_EXPORTER_LOG_LEVEL` | 日志级别（debug/info/warn/error） |
| `--hy-smi.path` | `hy-smi` | `HY_SMI_EXPORTER_BINARY_PATH` | hy-smi 二进制路径 |
| `--hy-smi.timeout` | `5s` | `HY_SMI_EXPORTER_TIMEOUT` | 单次命令超时时间 |
| `--hy-smi.collect-1s-stats` | `false` | `HY_SMI_EXPORTER_COLLECT_1S_STATS` | 是否采集耗时的 1 秒利用率和带宽统计；多卡节点建议保持关闭 |

## 部署

### systemd

仓库提供 `examples/hy-smi-exporter.service` 和
`examples/hy-smi-exporter.env`。将二进制安装到目标路径、确认
`hy-smi` 可执行后，再按目标发行版注册服务。

### 首版不支持的部署方式

Docker、Docker Compose、Kubernetes 和容器镜像不属于 `v1.0.0` 的兼容或
发布范围，仓库不提供相关部署文件。

---

## Prometheus 接入

在 Prometheus 配置文件中添加抓取任务：

```yaml
scrape_configs:
  - job_name: 'hy-smi-exporter'
    static_configs:
      - targets:
          - '<节点IP>:9584'
        labels:
          hostname: '<节点名>'
```

多节点场景建议使用文件服务发现，参考 `examples/prometheus.yml`。

---

## 目录结构

```
hy-smi-exporter/
├── cmd/hy-smi-exporter/main.go          # 入口：HTTP 服务 + /metrics
├── internal/
│   ├── config/config.go                 # 参数与环境变量配置
│   ├── hysmi/hysmi.go                   # hy-smi 命令调用 + JSON/文本解析
│   ├── hysmi/hysmi_test.go              # 解析和异常路径单元测试
│   ├── hysmi/testdata/*.json            # hy-smi --json 测试 fixture
│   └── exporter/
│       ├── collector.go                 # 每设备指标输出
│       ├── descs.go                     # Prometheus 指标描述符
│       └── collector_test.go            # 指标契约测试
├── examples/
│   ├── hy-smi-exporter.service          # systemd 服务示例
│   ├── hy-smi-exporter.env              # 环境变量示例
│   └── prometheus.yml                   # Prometheus 抓取配置示例
├── Makefile
├── go.mod / go.sum
├── LICENSE                              # Apache-2.0
├── THIRD_PARTY_NOTICES.md               # 第三方依赖及许可证清单
├── release-materials/                   # 发布确认、测试、验收和放行表单
└── README.md                            # 本文档
```

---

## 故障排查

| 现象 | 排查方法 |
|------|----------|
| `hy_smi_up` 为 0 | 检查日志及 `hy-smi --json`、`hy-smi` 是否可正常执行；确认驱动、工具路径和权限 |
| `hy_smi_up` 为 1 但设备数为 0 | 检查驱动模块、设备节点和厂商工具输出 |
| `hy-smi binary not found` 日志 | 确认 `hy-smi` 在 `$PATH` 中，或通过 `--hy-smi.path` 指定绝对路径 |
| `Permission denied` | 检查宿主机 `hy-smi` 权限、设备节点和运行用户 |
| 驱动未加载 | `lsmod \| grep -E "hydcu\|hycu"` 确认驱动已加载 |
| `hy_smi_info` 的 `unique_id` 或 `health` 为空 | 分别执行 `hy-smi --showuniqueid`、`hy-smi --healthcheck`，保留脱敏输出用于校准字段映射 |
| 利用率或带宽指标未出现 | 这些 1 秒统计项默认关闭；按需设置 `--hy-smi.collect-1s-stats=true`，并适当增大抓取超时 |

---

## 许可证

Apache-2.0，与 `node_exporter`、`dcgm-exporter`、`prometheus/client_golang` 保持一致。

---

## 参考资料

- [rookie-J/hy-smi](https://github.com/rookie-J/hy-smi) — hy-smi 命令完整参考（MIT）
- [reloaded/prometheus_consumer_amdgpu_exporter](https://github.com/reloaded/prometheus_consumer_amdgpu_exporter) — 架构参考（AMD Radeon）
- [NVIDIA/dcgm-exporter](https://github.com/NVIDIA/dcgm-exporter) — 指标命名规范参考
- [prometheus/client_golang](https://github.com/prometheus/client_golang) — Prometheus Go 客户端库
