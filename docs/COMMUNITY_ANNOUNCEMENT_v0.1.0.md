# hy-smi-exporter v0.1.0 发布公告

发布日期：`2026-11-13（计划）`

鼎茂正式发布 `hy-smi-exporter` `v0.1.0`。该项目是一款面向海光 DCU / HCU 设备的开源 Prometheus exporter，由鼎茂负责主开发、仓库管理和版本维护，海光团队作为技术协作者参与设备、驱动、`hy-smi` / `hcu-smi` 兼容性和指标语义确认。

项目遵循 Apache License 2.0 发布。

## 项目能力

项目通过宿主机已有的 `hy-smi` / `hcu-smi` 只读命令采集设备运行数据，并通过 `/metrics` 提供 Prometheus 指标，覆盖：

- 设备身份和设备数量；
- 温度、功耗和电压；
- HCU、CU、Wave、编解码及显存利用率；
- VRAM、GTT 和可用显存；
- 核心、显存及 SoC 时钟；
- PCIe、XHCL、DF 和 UMC；
- 健康状态、RAS、ECC 和坏页；
- 风扇及性能等级。

项目首版支持直接运行和 systemd 部署。Docker、Docker Compose、Kubernetes
和容器镜像不属于 `v0.1.0` 的兼容和发布范围。

## 项目地址

- GitHub：`https://github.com/di-matrix-ai/hy-smi-exporter`
- Release：`https://github.com/di-matrix-ai/hy-smi-exporter/releases/tag/v0.1.0`
- 指标说明：`https://github.com/di-matrix-ai/hy-smi-exporter/blob/main/metric.md`
- 兼容性矩阵：`https://github.com/di-matrix-ai/hy-smi-exporter/blob/main/docs/compatibility.md`

## 首版兼容范围

`v0.1.0` 仅对双方实际验证过的设备、驱动、固件、DTK 和 `hy-smi` / `hcu-smi` 组合做兼容性说明。

已验证组合：

| 硬件型号 | 驱动版本 | DTK 版本 | 工具版本 | 结论 |
| --- | --- | --- | --- | --- |
| `<待填写>` | `<待填写>` | `<待填写>` | `<待填写>` | 已验证 / 有条件支持 |

海光驱动能够支持某一硬件型号，不自动代表本 exporter 已经完成该组合的验证。后续将通过真实设备测试、Issue 和 Pull Request 持续扩展兼容范围。

## 安全说明

项目只调用查询类命令，不执行设备配置修改、驱动加载、设备重置或 RAS 注入。项目不分发海光驱动、固件或厂商工具二进制。首版不发布容器镜像。

## 社区参与

欢迎通过 GitHub Issues 反馈问题，通过 Pull Request 贡献设备适配、解析样例、指标、文档和部署方案。提交日志和工具输出前，请删除序列号、Unique ID、主机名、IP 地址、PCIe 地址和客户信息。

项目后续将向 Prometheus exporter 目录提交收录申请。目录收录属于社区项目展示，不代表 Prometheus 官方认证或技术支持。

## 联系方式

- 鼎茂发布账号：`@di-matrix-ai`，联系邮箱 `esteric_shen@di-matrix.com`；
- 海光技术协作者：王忠勤（`@zqwangadv`）、代金龙（`@kinglong08`）；
- 安全问题反馈：`esteric_shen@di-matrix.com`。
