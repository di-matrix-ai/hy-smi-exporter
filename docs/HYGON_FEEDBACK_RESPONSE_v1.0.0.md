# hy-smi-exporter v1.0.0 海光测试反馈整改说明

文档日期：2026 年 11 月 12 日
适用版本：`hy-smi-exporter v1.0.0` 候选版本
反馈环境：K100_AI，驱动 `6.3.31-V1.5.1`
整改方：上海鼎茂信息技术有限公司
代码分支：`fix-hygon-v1.0.0-feedback`
整改 PR：<https://github.com/di-matrix-ai/hy-smi-exporter/pull/2>
当前状态：代码、文档和 CI 验证已完成；海光协作者已于 2026 年 11 月 12 日批准 PR；K100_AI 定向复测作为发布后跟踪项

## 1. 整改结论

海光提出的三项建议均已完成代码和文档整改：

1. 项目涉及的设备管理命令已统一为 `hy-smi`；
2. 增加 K100_AI 设备唯一标识和健康状态的补充采集及字段兼容；
3. 增加 1 秒统计指标开关，默认关闭耗时较长的统计采集。

截至 2026 年 11 月 12 日，PR 的依赖校验、格式检查、竞态测试、单元测试、`go vet`、常规构建及 Linux `amd64`、`arm64` 交叉构建均已通过，海光协作者账号 `zqwangadv` 已在 PR #2 完成技术审核并批准变更。由于当前不具备 K100_AI 真机环境，唯一标识、健康状态以及 8 卡节点采集耗时仍需后续定向复测；该事项作为已公开的待验证项管理，不阻断 `v1.0.0` 发布。

## 2. 反馈与整改结果

| 序号 | 海光反馈 | 本次整改 | 当前验证状态 |
| --- | --- | --- | --- |
| 1 | DCU / HCU 没有 `hcu-smi` 命令，应统一为 `hy-smi` | 删除当前公开代码、配置、部署说明和发布文档中的 `hcu-smi` 描述；Exporter 默认且仅调用 `hy-smi` | 文档扫描及 CI 已通过 |
| 2 | K100_AI 上 `hy-smi --showuniqueid` 可显示 UUID，但 `hy_smi_info` 中 `unique_id`、`health` 为空 | 主查询缺少字段时，分别执行 `hy-smi --showuniqueid` 和 `hy-smi --healthcheck`；兼容 `Unique ID`、`UniqueID`、`UUID`、`Health`、`Health Status`、`Health State`、`Health Check`、`Device Health` 等字段名；按设备序号合并补采结果 | 单元测试、CI 和协作者审核已通过；K100_AI 真机复测发布后补充 |
| 3 | 部分指标统计过去 1 秒数据，8 卡机型执行耗时较长 | 新增 `--hy-smi.collect-1s-stats` 参数和 `HY_SMI_EXPORTER_COLLECT_1S_STATS` 环境变量；默认值为 `false`；默认不执行 `hy-smi -a`，不输出利用率及 PCIe、XHCL、DF、UMC 带宽统计 | 配置、采集和指标输出测试已通过；8 卡环境耗时对比发布后补充 |

## 3. 具体改动

### 3.1 命令名称统一

- 默认命令固定为 `hy-smi`；
- 命令路径仍可通过 `--hy-smi.path` 或 `HY_SMI_EXPORTER_BINARY_PATH` 指定；
- 当前发布范围内的 README、安装说明、兼容性说明、指标文档及发布材料均不再包含 `hcu-smi`；
- Exporter 只调用查询类参数，不执行设备配置、驱动加载、设备重置或 RAS 注入操作。

### 3.2 Unique ID 和健康状态补采

默认主查询为：

```bash
hy-smi --json
```

当主查询结果缺少设备唯一标识或健康状态时，Exporter 按需执行：

```bash
hy-smi --showuniqueid --json
hy-smi --healthcheck --json
```

如果对应参数不支持 JSON 输出，则自动回退到文本输出：

```bash
hy-smi --showuniqueid
hy-smi --healthcheck
```

补采结果按设备序号合并到主查询结果，并通过以下指标输出：

- `hy_smi_info{unique_id="...",health="..."}`；
- `hy_smi_health_info{health="..."}`。

唯一标识或健康状态的可选补充查询失败时，会记录调试日志，但不会将已成功的主查询误判为整体采集失败。

### 3.3 1 秒统计开关

默认配置：

```bash
--hy-smi.collect-1s-stats=false
```

默认关闭时：

- 主查询使用 `hy-smi --json`；
- 不执行 `hy-smi -a`；
- 不输出利用率及 PCIe、XHCL、DF、UMC 带宽等 1 秒统计指标；
- 减少 8 卡等多卡节点上的命令执行时间和 Prometheus 抓取压力。

显式开启时：

```bash
--hy-smi.collect-1s-stats=true
```

主查询改为 `hy-smi -a --json`，JSON 不可用时回退到 `hy-smi -a`，并输出对应的利用率和带宽指标。生产环境启用后，应根据实测最大耗时调整 `--hy-smi.timeout` 和 Prometheus `scrape_timeout`。

## 4. 自动化验证结果

PR 检查地址：<https://github.com/di-matrix-ai/hy-smi-exporter/pull/2/checks>

| 检查项 | 结果 |
| --- | --- |
| `go mod verify` 依赖完整性检查 | 通过 |
| `gofmt` 格式检查 | 通过 |
| `go test -race ./...` 单元测试及竞态检查 | 通过 |
| `go vet ./...` 静态检查 | 通过 |
| `go build` 常规构建 | 通过 |
| Linux `amd64` 交叉构建 | 通过 |
| Linux `arm64` 交叉构建 | 通过 |
| 发布材料完整性及占位符检查 | 通过 |

本次新增测试覆盖：

- UUID、Unique ID 和健康状态字段别名解析；
- 多卡补采结果按设备序号合并；
- 1 秒统计默认关闭；
- 启用开关后使用 `hy-smi -a --json`；
- 默认关闭时相关利用率和带宽指标不输出。

## 5. 发布后海光定向复测清单

### 5.1 基础信息

请记录并脱敏后返回：

- 设备型号及卡数量；
- 驱动版本、Linux 发行版和内核版本；
- 候选二进制版本及 SHA-256；
- 测试日期和测试结论。

### 5.2 Unique ID 和健康状态复测

使用候选二进制在 K100_AI / `6.3.31-V1.5.1` 节点启动：

```bash
./hy-smi-exporter \
  --hy-smi.path=/opt/hyhal/bin/hy-smi \
  --hy-smi.timeout=15s \
  --hy-smi.collect-1s-stats=false
```

执行：

```bash
hy-smi --showuniqueid
hy-smi --healthcheck
curl -s http://127.0.0.1:9584/metrics | grep '^hy_smi_info'
curl -s http://127.0.0.1:9584/metrics | grep '^hy_smi_health_info'
```

通过标准：

- `hy_smi_up=1`，设备数量与物理卡数量一致；
- `hy_smi_info` 的 `unique_id` 非空，并与 `hy-smi --showuniqueid` 对应设备一致；
- `hy_smi_info` 的 `health` 非空，并与 `hy-smi --healthcheck` 一致；
- `hy_smi_health_info` 正常输出；
- 多卡设备的设备序号、UUID 和健康状态不存在错位。

### 5.3 1 秒统计及耗时复测

1. 使用默认关闭模式，连续请求 `/metrics` 10 次，记录最小、最大和平均耗时；
2. 使用 `--hy-smi.collect-1s-stats=true` 重启 Exporter；
3. 确认利用率及 PCIe、XHCL、DF、UMC 带宽指标恢复输出；
4. 再连续请求 `/metrics` 10 次，记录最小、最大和平均耗时；
5. 根据开启后的最大耗时评估 Prometheus `scrape_timeout`。

## 6. 发布安排与后续跟踪

海光协作者账号 `zqwangadv` 已于 2026 年 11 月 12 日批准 PR #2，协作者技术审核要求已经满足。`v1.0.0` 按以下条件发布：

1. 鼎茂更新兼容性矩阵、RC 验收记录和发布放行记录，如实披露 K100_AI 待复测状态；
2. 合并后的 `main` 分支 CI 再次通过；
3. 正式 Tag 指向最终合并提交，Release 产物及 SHA-256 校验完整；
4. K100_AI / `6.3.31-V1.5.1` 定向复测在发布后继续跟踪，复测通过后补充测试记录和兼容矩阵；如发现代码问题，通过后续补丁版本修复。

PR 审核仅表示相关技术协作者对本次变更的代码审阅意见，不表述为海光官方认证、产品背书或联合发布。复测材料对外传递前，应对序列号、UUID、主机名、IP 地址和 PCIe 地址进行脱敏。
