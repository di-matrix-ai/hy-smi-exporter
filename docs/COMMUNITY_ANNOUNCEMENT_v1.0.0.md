# hy-smi-exporter v1.0.0 发布公告

发布日期：`2026-11-13`

鼎茂正式发布 `hy-smi-exporter` `v1.0.0`。该项目是一款面向海光 DCU / HCU 设备的开源 Prometheus exporter，由鼎茂独立开发、管理仓库、维护版本并承担发布责任。相关设备厂商技术人员及社区贡献者可按个人身份参与测试、兼容性沟通和技术建议；本项目不是海光官方产品，不代表海光官方认证、联合发布或质量承诺。

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
和容器镜像不属于 `v1.0.0` 的兼容和发布范围。

## 项目地址

- GitHub：`https://github.com/di-matrix-ai/hy-smi-exporter`
- Release：`https://github.com/di-matrix-ai/hy-smi-exporter/releases/tag/v1.0.0`
- 指标说明：`https://github.com/di-matrix-ai/hy-smi-exporter/blob/main/metric.md`
- 兼容性矩阵：`https://github.com/di-matrix-ai/hy-smi-exporter/blob/main/docs/compatibility.md`

## 首版兼容范围

`v1.0.0` 基于驱动提供的通用 `hy-smi` / `hcu-smi` 命令接口适配设备。目标设备由对应驱动正常支持且命令输出保持兼容时，可使用本 Exporter；真实设备验证状态单独列明。

支持与验证组合：

| 硬件型号 | 驱动版本 | DTK 版本 | 工具版本 | 结论 |
| --- | --- | --- | --- | --- |
| BW1100 / BW1000 / BW150 / BW100 / BW10 / K100-AI / K100 / Z100L / Z100 | `rock-6.3.30+` | 26.04 | 驱动提供的兼容命令集 | 支持，基于通用命令接口 |
| BW100 | `6.3.28-V1.3.0b` | 未单独记录 | `/opt/hyhal/bin/hy-smi`，版本号未单独记录 | 历史真机验证通过 |

上述支持范围以驱动及厂商命令接口可正常工作为前提。“支持”表示共用当前命令采集和解析实现，“历史真机验证通过”表示已有可追溯的真实设备测试记录。后续将通过 Issue、Pull Request 和新增测试记录持续完善兼容性信息。

## 安全说明

项目只调用查询类命令，不执行设备配置修改、驱动加载、设备重置或 RAS 注入。项目不分发海光驱动、固件或厂商工具二进制。首版不发布容器镜像。

## 社区参与

欢迎通过 GitHub Issues 反馈问题，通过 Pull Request 贡献设备适配、解析样例、指标、文档和部署方案。提交日志和工具输出前，请删除序列号、Unique ID、主机名、IP 地址、PCIe 地址和客户信息。

项目后续将向 Prometheus exporter 目录提交收录申请。目录收录属于社区项目展示，不代表 Prometheus 官方认证或技术支持。

## 联系方式

- 鼎茂发布账号：`@estericShen`，联系邮箱 `esteric_shen@di-matrix.com`；
- 技术协作者/社区贡献者：王忠勤（`@zqwangadv`）、代金龙（`@kinglong08`），以个人身份参与；
- 安全问题反馈：`esteric_shen@di-matrix.com`。
