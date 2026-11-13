# hy-smi-exporter v1.0.0 Release Notes

发布日期：`2026-11-13`

仓库：`https://github.com/di-matrix-ai/hy-smi-exporter`

## 版本定位

`v1.0.0` 是 `hy-smi-exporter` 的首个公开版本。项目通过宿主机已有的 `hy-smi` 只读命令采集海光 DCU / HCU 设备状态和运行数据，并通过 `/metrics` 暴露 Prometheus 指标。

本版本基于驱动提供的通用 `hy-smi` 命令接口适配设备。目标设备由对应驱动正常支持且命令输出保持兼容时，可使用本 Exporter；真实设备验证记录在兼容性矩阵中单独标明。

## 主要功能

- 默认使用 `hy-smi --json` 轻量主查询，JSON 不可用时回退到 `hy-smi` 文本输出；
- 主查询缺少字段时，通过 `--showuniqueid`、`--healthcheck` 补采设备唯一标识和健康状态；
- 主查询缺少显存数据时执行只读显存补充查询；
- 提供 `--hy-smi.collect-1s-stats` 配置，利用率和带宽等 1 秒统计项默认关闭；
- 输出设备身份、温度、功耗、利用率、显存、时钟和 PCIe 指标；
- 输出链路、健康、RAS、ECC、坏页、风扇和性能等级指标；
- 提供 `/metrics` 和 `/healthz` 端点；
- 提供 `hy_smi_up` 指标区分采集失败与正常的零设备结果；
- 支持直接运行和 systemd 部署；
- 提供 Linux `amd64`、`arm64` 二进制；
- 首版不发布、不承诺 Docker、Docker Compose、Kubernetes 或容器镜像兼容性。

## 支持与验证范围

| 硬件型号 | 驱动版本 | 固件版本 | DTK 版本 | `hy-smi` 版本 | Linux / 内核 | 结果 |
| --- | --- | --- | --- | --- | --- | --- |
| BW1100 / BW1000 / BW150 / BW100 / BW10 / K100-AI / K100 / Z100L / Z100 | `rock-6.3.30+` | 随具体设备和驱动版本确定 | 26.04 | 驱动提供的兼容命令集 | 以驱动支持的 Linux 环境为准 | 支持，基于通用命令接口 |
| BW100 | `6.3.28-V1.3.0b` | VBIOS `6.312.002400Q.984889` | 未单独记录 | 路径 `/opt/hyhal/bin/hy-smi`，版本号未单独记录 | Kylin Linux Advanced Server V10 / `4.19.90-89.11.v2401.ky10.x86_64` | 历史真机验证通过 |
| K100_AI（K100-AI） | `6.3.31-V1.5.1` | 未提供 | 未提供 | `hy-smi --showuniqueid` 可显示 UUID | 未提供 | 代码修复和协作者审核已完成，真机复测发布后补充 |

完整信息见 [`docs/compatibility.md`](compatibility.md)。

## 发布物

| 发布物 | 地址或校验信息 |
| --- | --- |
| Linux amd64 | `https://github.com/di-matrix-ai/hy-smi-exporter/releases/download/v1.0.0/hy-smi-exporter_v1.0.0_linux_amd64.tar.gz` |
| Linux arm64 | `https://github.com/di-matrix-ai/hy-smi-exporter/releases/download/v1.0.0/hy-smi-exporter_v1.0.0_linux_arm64.tar.gz` |
| SHA-256 | `https://github.com/di-matrix-ai/hy-smi-exporter/releases/download/v1.0.0/checksums.txt` |
| 源码 | `https://github.com/di-matrix-ai/hy-smi-exporter/tree/v1.0.0` |

## 安装和验证

```bash
./hy-smi-exporter \
  --web.listen-address=:9584 \
  --hy-smi.path=/usr/bin/hy-smi \
  --hy-smi.timeout=5s \
  --hy-smi.collect-1s-stats=false

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
- 项目不分发海光驱动、固件、`hy-smi` 二进制或商业授权文件；
- 生产环境应限制 `/metrics` 的网络访问；
- Issue 和日志中不得提交未脱敏的设备及基础设施信息。

## 已知限制

- 当前版本未提供 vHCU 维度的独立指标；Exporter 自身的 Go 运行时和进程指标由 Prometheus 客户端标准采集器提供；
- 不同驱动或工具版本可能存在字段、单位和枚举差异；
- 驱动、固件、DTK 或工具升级后需要重新执行兼容性回归；
- 多个 Prometheus 实例重复抓取会增加 `hy-smi` 调用次数；
- 利用率和 PCIe、XHCL、DF、UMC 带宽等 1 秒统计指标默认不输出，启用后多卡节点需要关注采集耗时；
- K100_AI / `6.3.31-V1.5.1` 的 Unique ID 和健康状态修复已经通过自动化测试和协作者代码审核，但尚未完成真机复测；该组合按“有条件支持、待复测”发布，复测结论将在发布后补充；
- 容器部署不属于 `v1.0.0` 支持范围；
- 未完成真机验证的型号按通用命令接口提供支持；若厂商工具输出发生变化，可能需要补充解析适配。

## 维护与反馈

- 一般问题：`https://github.com/di-matrix-ai/hy-smi-exporter/issues`；
- 安全问题：按照 [`SECURITY.md`](../SECURITY.md) 私密反馈；
- 鼎茂发布账号：`@estericShen`，联系邮箱 `esteric_shen@di-matrix.com`；
- 技术协作者/社区贡献者：`@zqwangadv`、`@kinglong08`，以个人身份参与，不代表厂商官方授权或认证。
