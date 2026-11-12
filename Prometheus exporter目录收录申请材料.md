# Prometheus Exporter 目录收录申请材料

## 1 提交结论

`hy-smi-exporter` 计划在首个公开 Release 发布后，向
[`prometheus/docs`](https://github.com/prometheus/docs) 提交 exporter 目录收录
Pull Request。目录收录是社区项目展示，不代表 Prometheus 官方认证、兼容性
背书或技术支持。

当前仓库已公开。正式 `v1.0.0` Release、候选版本复测、兼容矩阵和鼎茂正式
放行完成前，不提交目录 PR，也不申请默认端口登记。

## 2 项目信息

| 项目 | 内容 |
| --- | --- |
| 项目名 | `hy-smi-exporter` |
| 仓库 | `https://github.com/di-matrix-ai/hy-smi-exporter` |
| 维护方 | Shanghai Ding Mao Information Technology Co., Ltd. |
| 技术协作 | 相关设备厂商技术人员以个人 GitHub 账号参与，不代表厂商官方授权或认证 |
| 许可证 | Apache License 2.0 |
| 语言 | Go |
| 运行平台 | Linux `amd64`、Linux `arm64`，以 Release 验证结果为准 |
| 依赖 | 宿主机已安装并可执行的 `hy-smi` |
| 指标端点 | `/metrics` |
| 默认端口 | `9584`，仓库公开后须再次复核并登记或调整 |

英文简介：

> Prometheus exporter for Hygon DCU / HCU accelerators using read-only
> `hy-smi` queries.

## 3 Prometheus 文档改动

按当前社区结构，目标文件为：

```text
prometheus/docs/docs/instrumenting/exporters.md
```

在 hardware-related exporter 区域按相邻条目的格式和排序增加一条链接。提交时
只改目录条目，不把项目介绍、兼容矩阵或发布说明复制到 Prometheus 文档仓库。

建议条目：

```markdown
* [hy-smi-exporter](https://github.com/di-matrix-ai/hy-smi-exporter)
```

正式提交前必须查看目标文件最新内容，按当时的标题、项目符号、字母排序和
校验规则调整，不假设本草稿格式长期不变。

## 4 Pull Request 描述草稿

```markdown
## Add hy-smi-exporter

This PR adds hy-smi-exporter to the hardware-related exporter directory.

- Repository: https://github.com/di-matrix-ai/hy-smi-exporter
- License: Apache-2.0
- Maintainer: Shanghai Ding Mao Information Technology Co., Ltd.

The exporter collects Hygon DCU / HCU device metrics through host-provided,
read-only `hy-smi` commands. Supported hardware and driver
combinations are documented in the project's verified compatibility matrix.
```

## 5 提交前置条件

- [ ] 仓库已按批准范围转为公开，匿名访问正常；
- [ ] 对外 Git 历史不含内部审批表、个人联系方式和真实设备标识；
- [ ] Apache-2.0 `LICENSE`、`NOTICE`、README、安装和指标文档完整；
- [ ] 至少一个正式 GitHub Release 可下载，Linux 产物和校验和已验证；
- [ ] 技术协作和非官方发布边界已留档，真实设备报告、兼容矩阵和指标语义已经确认；
- [ ] Release Notes 和公告没有 `TBD`、示例地址或未确认承诺；
- [ ] GitHub Issues、维护者和私密安全报告渠道有效；
- [ ] 默认端口 `9584` 已在 Prometheus 端口登记页面复核，无冲突或已完成登记；
- [ ] PR 只陈述项目事实，不使用“官方”“认证”“推荐”等背书性表述。

## 6 DCO 和提交流程

Prometheus 文档仓库要求按其当前贡献规则提交。提交 commit 时使用 DCO
签署，例如：

```bash
git commit -s -m "Add hy-smi-exporter"
```

确认 commit message 中存在与提交身份一致的 `Signed-off-by:` 行。若 CI 或维护者
要求修改，后续 commit 同样需要 DCO 签署。

流程：

1. Fork `prometheus/docs`，从最新默认分支创建工作分支；
2. 仅在 `docs/instrumenting/exporters.md` 的 hardware-related 区域增加条目；
3. 按仓库说明运行文档格式或链接检查；
4. 使用 `git commit -s` 提交并创建 Pull Request；
5. 跟进社区审核，记录 PR URL、审核意见、合并 commit 和日期；
6. 将结果填写到 `release-materials/08-prometheus-submission-record.md`。

## 7 统一对外口径

`hy-smi-exporter` 由鼎茂独立开发、维护和发布。相关设备厂商技术人员可按个人
身份参与设备、驱动、厂商工具兼容性和指标语义协作，但不代表厂商官方授权、
联合发布或质量担保。支持范围以鼎茂公开兼容矩阵中的真实设备测试结果为准。
Prometheus exporter 目录收录仅表示项目链接被社区文档列出。
