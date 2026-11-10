# Release Candidate 验收记录

RC 版本：`TBD`  
Git commit：`TBD`  
验收日期：`TBD`  
验收负责人：`TBD`

## 1 仓库和合规

| 检查项 | 结果 | 证据 / 备注 |
| --- | --- | --- |
| 公开仓库和 Go module 地址一致 | `TBD` | `TBD` |
| `LICENSE`、`NOTICE` 和版权主体正确 | `TBD` | `TBD` |
| 维护者、CODEOWNERS 和安全渠道有效 | `TBD` | `TBD` |
| 第三方许可证清单已审核 | `TBD` | `TBD` |
| 海光公开授权已留档 | `TBD` | `TBD` |
| 敏感信息扫描无阻断问题 | `TBD` | `TBD` |

## 2 工程验证

| 检查项 | 结果 | 证据 / 备注 |
| --- | --- | --- |
| `go mod verify` | `TBD` | CI 链接 |
| 格式检查 | `TBD` | CI 链接 |
| `go test -race ./...` | `TBD` | CI 链接 |
| `go vet ./...` | `TBD` | CI 链接 |
| Linux amd64 构建和启动 | `TBD` | 校验值 / 日志 |
| Linux arm64 构建和启动 | `TBD` | 校验值 / 日志 |
| 容器镜像 | 不适用：首版不发布 | 不纳入 v0.1.0 |
| SBOM 和 provenance | 不适用：首版不发布容器 | 不纳入 v0.1.0 |
| Release 归档内容检查 | `TBD` | 文件清单 |

## 3 真实设备和部署验证

| 检查项 | 结果 | 证据 / 备注 |
| --- | --- | --- |
| 真实设备测试报告完成 | `TBD` | 报告编号 |
| 兼容性矩阵只包含实测组合 | `TBD` | 文档链接 |
| 指标语义经海光确认 | `TBD` | 确认记录 |
| 直接运行 / systemd | `TBD` | 测试记录 |
| Docker / Compose | 不适用：首版不支持 | 不纳入 v0.1.0 |
| Kubernetes DaemonSet | 不适用：首版不支持 | 不纳入 v0.1.0 |
| Prometheus 抓取和稳定性 | `TBD` | 测试记录 |

## 4 对外材料

| 检查项 | 结果 | 证据 / 备注 |
| --- | --- | --- |
| README 中英文内容一致 | `TBD` | `TBD` |
| 指标和安装文档准确 | `TBD` | `TBD` |
| Release Notes 无占位符 | `TBD` | `TBD` |
| 发布公告无占位符 | `TBD` | `TBD` |
| 已知限制和支持边界已披露 | `TBD` | `TBD` |
| Prometheus 收录材料已准备 | `TBD` | Release 后提交 |

## 5 验收结论

结论：`通过 / 有条件通过 / 不通过 / TBD`  
遗留问题及截止日期：`TBD`  
鼎茂验收负责人：`TBD`  
海光技术复核人：`TBD`  
确认日期：`TBD`
