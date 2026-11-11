# Release Candidate 验收记录

RC 版本：`v1.0.0`  
Git commit：以 `v1.0.0` Tag 最终指向为准  
验收日期：2026-11-13  
验收负责人：Esteric Shen / 上海鼎茂信息技术有限公司发布团队

## 1 仓库和合规

| 检查项 | 结果 | 证据 / 备注 |
| --- | --- | --- |
| 公开仓库和 Go module 地址一致 | 通过 | GitHub 仓库及 `go.mod` 均为 `github.com/di-matrix-ai/hy-smi-exporter` |
| `LICENSE`、`NOTICE` 和版权主体正确 | 通过 | Apache-2.0；上海鼎茂信息技术有限公司中英文主体已确认 |
| 维护者、CODEOWNERS 和安全渠道有效 | 通过 | `MAINTAINERS.md`、`.github/CODEOWNERS`、`SECURITY.md` |
| 第三方许可证清单已审核 | 条件通过 | 公开清单已完成，内部合规确认在正式放行时留档 |
| 技术协作和非官方发布边界已留档 | 通过 | `02-collaboration-and-publication-boundary.md` |
| 敏感信息扫描无阻断问题 | 发布前确认 | Tag 前对最终提交和 GitHub 历史执行检查 |

## 2 工程验证

| 检查项 | 结果 | 证据 / 备注 |
| --- | --- | --- |
| `go mod verify` | 等待 Tag/CI | 正式 CI 链接 |
| 格式检查 | 等待 Tag/CI | 正式 CI 链接 |
| `go test -race ./...` | 等待 Tag/CI | 正式 CI 链接 |
| `go vet ./...` | 等待 Tag/CI | 正式 CI 链接 |
| Linux amd64 构建和启动 | 等待 Draft Release | 校验值 / 日志 |
| Linux arm64 构建和启动 | 等待 Draft Release | 校验值 / 日志 |
| Release 归档内容检查 | 等待 Draft Release | 文件清单及 `checksums.txt` |

## 3 真实设备和部署验证

| 检查项 | 结果 | 证据 / 备注 |
| --- | --- | --- |
| 真实设备测试报告完成 | 条件通过 | `HY-SMI-RDT-20260914-01`；BW100 历史实测，当前版本未重新真机回归 |
| 兼容性矩阵区分通用接口支持和真实设备验证 | 已完成 | `docs/compatibility.md` |
| 指标语义由上海鼎茂信息技术有限公司定稿 | 已完成 | `04-metric-semantics-confirmation.md` |
| 直接运行 / systemd | 历史验证通过 | `03-real-device-test-report.md` |
| Prometheus 抓取和稳定性 | 历史验证通过 | 历史记录约 1 天 22 小时 systemd 运行，Prometheus Target 为 UP |

## 4 对外材料

| 检查项 | 结果 | 证据 / 备注 |
| --- | --- | --- |
| README 中英文核心口径一致 | 通过 | `README.md`、`README.en.md` |
| 指标和安装文档准确 | 通过 | `metric.md`、`docs/installation.md` |
| Release Notes 无占位符 | 通过 | `docs/RELEASE_NOTES_v1.0.0.md` |
| 发布公告无占位符 | 通过 | `docs/COMMUNITY_ANNOUNCEMENT_v1.0.0.md` |
| 已知限制和支持边界已披露 | 通过 | 兼容性、容器范围、敏感标签及历史测试边界均已公开 |
| Prometheus 收录材料已准备 | 已准备 | GitHub Release 发布后提交并补充 PR 地址 |

## 5 验收结论

结论：文档、兼容性口径和历史验证材料条件通过；完成正式 CI、Draft Release 产物校验和内部放行后可发布。  
遗留问题及截止日期：CI、二进制、校验值、敏感信息扫描和内部签批，均应在公开 Release 前完成。  
鼎茂验收负责人：正式放行时签署  
技术协作复核人：可选，不作为发布前置条件  
确认日期：2026-11-13
