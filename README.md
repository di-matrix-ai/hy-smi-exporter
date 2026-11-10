# hy-smi-exporter

面向海光 DCU / HCU 加速器的 Prometheus 指标采集器，实际支持范围以 [兼容性矩阵](docs/compatibility.md) 中完成真实设备验证的组合为准。
通过只读命令 `hy-smi -a --json` 采集每张卡的运行数据；对于不支持该 JSON 组合参数的驱动版本，自动回退到 `hy-smi -a` 文本输出，以 `hy_smi_` 前缀暴露在 `/metrics` 端点。

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

**本程序只调用只读命令，绝不执行任何写入操作。** 每次抓取依次尝试：

```
hy-smi -a --json
hy-smi -a   # 仅当上一个命令不支持 JSON 或解析失败时
hy-smi --showmeminfo all --json       # 仅当主查询缺少显存统计时
hy-smi --showmemavailable --json      # 仅当主查询缺少显存统计时
```

不会调用 `--set*`、`--rasinject`、`--load*driver`、`--reset*` 等任何修改硬件配置的选项。

---

## 采集指标

所有指标以 `hy_smi_` 为前缀，最少携带 `gpu`（设备序号）和 `instance`（主机名）标签。

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `hy_smi_devices_total` | gauge | 节点上 HCU 设备总数 |
| `hy_smi_info` | gauge (=1) | 设备身份信息（device_id / serial / product_name / vbios / driver / bus / health / mode 等） |
| `hy_smi_temperature_celsius` | gauge | 各传感器温度（edge / junction / mem / hotspot） |
| `hy_smi_power_watts` | gauge | 平均功耗（W） |
| `hy_smi_power_cap_watts` | gauge | 当前功耗上限（W） |
| `hy_smi_power_max_watts` | gauge | 最大允许功耗（W） |
| `hy_smi_voltage_volts` | gauge | 核心电压（V） |
| `hy_smi_hcu_utilization_percent` | gauge | HCU 活跃时间占比（0-100） |
| `hy_smi_cu_utilization_percent` | gauge | CU 波前利用率 |
| `hy_smi_wave_utilization_percent` | gauge | 波前驻留率 |
| `hy_smi_decoder_utilization_percent` | gauge | 解码器利用率 |
| `hy_smi_encoder_utilization_percent` | gauge | 编码器利用率 |
| `hy_smi_memory_utilization_percent` | gauge | 内存控制器利用率 |
| `hy_smi_vram_utilization_percent` | gauge | 显存利用率 |
| `hy_smi_vram_total_bytes` | gauge | 显存总量（字节） |
| `hy_smi_vram_used_bytes` | gauge | 已用显存（字节） |
| `hy_smi_vram_available_bytes` | gauge | 可用显存（字节） |
| `hy_smi_gtt_total_bytes` | gauge | GTT 内存总量（字节） |
| `hy_smi_gtt_used_bytes` | gauge | 已用 GTT（字节） |
| `hy_smi_clock_mhz` | gauge | 当前时钟频率（sclk / mclk / socclk，MHz） |
| `hy_smi_clock_min_mhz` | gauge | 最低可用时钟频率 |
| `hy_smi_clock_max_mhz` | gauge | 最高可用时钟频率 |
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
| `hy_smi_fan_speed_percent` | gauge | 风扇 PWM 占空比（0-100） |
| `hy_smi_performance_level_info` | gauge (=1) | DPM 性能等级 |
| `hy_smi_boost_level_info` | gauge (=1) | Boost 等级 |

---

## 环境要求

- Go 1.23+（编译时需要）
- `hy-smi` 或 `hcu-smi` 命令可执行（在 `$PATH` 中或通过 `--hy-smi.path` 指定）
- 驱动、固件、DTK 和 `hy-smi` / `hcu-smi` 版本以 [兼容性矩阵](docs/compatibility.md) 中双方实际验证的组合为准
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

默认监听 `:9393`，访问 `http://<HOST>:9393/metrics` 查看指标。

### 常用参数

```bash
./bin/hy-smi-exporter \
  --web.listen-address=:9393 \
  --hy-smi.path=/usr/bin/hy-smi \
  --hy-smi.timeout=5s \
  --log.level=info
```

### 全部参数说明

| 参数 | 默认值 | 环境变量 | 说明 |
|------|--------|----------|------|
| `--web.listen-address` | `:9393` | `HY_SMI_EXPORTER_LISTEN_ADDR` | 监听地址 |
| `--web.telemetry-path` | `/metrics` | `HY_SMI_EXPORTER_METRICS_PATH` | 指标路径 |
| `--log.level` | `info` | `HY_SMI_EXPORTER_LOG_LEVEL` | 日志级别（debug/info/warn/error） |
| `--hy-smi.path` | `hy-smi` | `HY_SMI_EXPORTER_BINARY_PATH` | hy-smi / hcu-smi 二进制路径 |
| `--hy-smi.timeout` | `5s` | `HY_SMI_EXPORTER_TIMEOUT` | 单次命令超时时间 |

`NODE_NAME` 环境变量可覆盖 `instance` 标签。

---

## 部署

### systemd

仓库提供 `examples/hy-smi-exporter.service` 和
`examples/hy-smi-exporter.env`。将二进制安装到目标路径、确认
`hy-smi` / `hcu-smi` 可执行后，再按目标发行版注册服务。

### 首版不支持的部署方式

Docker、Docker Compose、Kubernetes 和容器镜像不属于 `v0.1.0` 的兼容或
发布范围。仓库中的相关文件仅作为后续适配的工作底稿，不应作为已验证的
生产部署方案使用。

---

## Prometheus 接入

在 Prometheus 配置文件中添加抓取任务：

```yaml
scrape_configs:
  - job_name: 'hy-smi-exporter'
    static_configs:
      - targets:
          - '<节点IP>:9393'
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
│   ├── hysmi/hysmi_test.go              # 单元测试（5 个用例）
│   ├── hysmi/testdata/*.json            # hy-smi --json 测试 fixture
│   └── exporter/
│       ├── collector.go                 # 每设备指标输出
│       └── descs.go                     # Prometheus 指标描述符
├── examples/
│   ├── docker-compose.yml               # 后续容器适配工作底稿（首版不支持）
│   ├── daemonset.yaml                   # 后续容器适配工作底稿（首版不支持）
│   └── prometheus.yml                   # Prometheus 抓取配置示例
├── Dockerfile
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
| 指标全为空 / 设备数为 0 | 检查 `hy-smi -a --json` 和 `hy-smi -a` 是否可正常执行；检查 `/dev/kfd`、`/dev/dri` 是否已映射 |
| `hy-smi binary not found` 日志 | 确认 `hy-smi` 在 `$PATH` 中，或通过 `--hy-smi.path` 指定绝对路径 |
| `Permission denied` | 检查宿主机 `hy-smi` / `hcu-smi` 权限、设备节点和运行用户；首版不支持容器权限配置 |
| 驱动未加载 | `lsmod \| grep -E "hydcu\|hycu"` 确认驱动已加载 |
| 指标字段值为 0 | 在真实海光机器上执行 `hy-smi -a > sample.txt` 或 `hy-smi -a --json > sample.json`，将样本提供给开发者以校准字段映射 |

---

## 许可证

Apache-2.0，与 `node_exporter`、`dcgm-exporter`、`prometheus/client_golang` 保持一致。

---

## 参考资料

- [rookie-J/hy-smi](https://github.com/rookie-J/hy-smi) — hy-smi 命令完整参考（MIT）
- [reloaded/prometheus_consumer_amdgpu_exporter](https://github.com/reloaded/prometheus_consumer_amdgpu_exporter) — 架构参考（AMD Radeon）
- [NVIDIA/dcgm-exporter](https://github.com/NVIDIA/dcgm-exporter) — 指标命名规范参考
- [prometheus/client_golang](https://github.com/prometheus/client_golang) — Prometheus Go 客户端库
