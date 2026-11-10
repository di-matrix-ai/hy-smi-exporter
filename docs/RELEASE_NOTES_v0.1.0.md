# hy-smi-exporter v0.1.0 Release Notes

发布日期：`2026-11-13（计划）`

仓库：`https://github.com/di-matrix-ai/hy-smi-exporter`

## 版本定位

`v0.1.0` 是 `hy-smi-exporter` 的首个公开版本。项目通过宿主机已有的 `hy-smi` / `hcu-smi` 只读命令采集海光 DCU / HCU 设备状态和运行数据，并通过 `/metrics` 暴露 Prometheus 指标。

本版本仅对兼容性矩阵中实际完成验证的设备、驱动、固件、DTK 和工具组合负责，不默认承诺覆盖全部海光设备和驱动版本。

## 主要功能

- 支持 `hy-smi -a --json` 主查询；
- JSON 不可用或解析失败时回退到 `hy-smi -a` 文本输出；
- 主查询缺少显存数据时执行只读显存补充查询；
- 输出设备身份、温度、功耗、利用率、显存、时钟和 PCIe 指标；
- 输出链路、健康、RAS、ECC、坏页、风扇和性能等级指标；
- 提供 `/metrics` 和 `/healthz` 端点；
- 提供 `hy_smi_up` 指标区分采集失败与正常的零设备结果；
- 支持直接运行和 systemd 部署；
- 提供 Linux `amd64`、`arm64` 二进制；
- 首版不发布、不承诺 Docker、Docker Compose、Kubernetes 或容器镜像兼容性。

## 已验证兼容范围

| 硬件型号 | 驱动版本 | 固件版本 | DTK 版本 | `hy-smi` / `hcu-smi` 版本 | Linux / 内核 | 结果 |
| --- | --- | --- | --- | --- | --- | --- |
| `<待填写>` | `<待填写>` | `<待填写>` | `<待填写>` | `<待填写>` | `<待填写>` | 已验证 / 有条件支持 |

完整信息见 [`docs/compatibility.md`](compatibility.md)。

## 发布物

| 发布物 | 地址或校验信息 |
| --- | --- |
| Linux amd64 | `<Release asset URL>` |
| Linux arm64 | `<Release asset URL>` |
| SHA-256 | `<checksums.txt URL>` |
| 源码 | `https://github.com/di-matrix-ai/hy-smi-exporter/tree/v0.1.0` |

## 安装和验证

```bash
./hy-smi-exporter \
  --web.listen-address=:9584 \
  --hy-smi.path=/usr/bin/hy-smi \
  --hy-smi.timeout=5s

curl http://127.0.0.1:9584/healthz
curl http://127.0.0.1:9584/metrics | grep hy_smi_
```

安装前请阅读：

- [`README.md`](../README.md)；
- [`docs/installation.md`](installation.md)；
- [`docs/compatibility.md`](compatibility.md)；
- [`metric.md`](../metric.md)。

## 安全边界

- exporter 只调用查询类命令；
- 不执行驱动加载、设备重置、RAS 注入或设备配置修改；
- 项目不分发海光驱动、固件、`hy-smi` / `hcu-smi` 二进制或商业授权文件；
- 生产环境应限制 `/metrics` 的网络访问；
- Issue 和日志中不得提交未脱敏的设备及基础设施信息。

## 已知限制

- 当前版本不提供进程级指标和 vHCU 指标；
- 不同驱动或工具版本可能存在字段、单位和枚举差异；
- 驱动、固件、DTK 或工具升级后需要重新执行兼容性回归；
- 多个 Prometheus 实例重复抓取会增加 `hy-smi` / `hcu-smi` 调用次数；
- 容器部署不属于 `v0.1.0` 支持范围；
- 其他限制：待双方真实设备测试后补充。

## 维护与反馈

- 一般问题：`https://github.com/di-matrix-ai/hy-smi-exporter/issues`；
- 安全问题：按照 [`SECURITY.md`](../SECURITY.md) 私密反馈；
- 鼎茂发布账号：`@estericShen`，联系邮箱 `esteric_shen@di-matrix.com`；
- 海光技术协作者：`@zqwangadv`、`@kinglong08`。
