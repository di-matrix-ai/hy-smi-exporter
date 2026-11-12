# 海光测试反馈整改说明与复测清单

适用版本：`hy-smi-exporter v1.0.0` 候选版本  
反馈环境：K100_AI，驱动 `6.3.31-V1.5.1`  
整改方：上海鼎茂信息技术有限公司  
状态：代码和文档已整改，待海光候选版本复测

## 1 反馈与整改结果

| 序号 | 海光反馈 | 整改结果 | 复测状态 |
| --- | --- | --- | --- |
| 1 | DCU / HCU 设备管理命令统一为 `hy-smi` | 已删除代码注释、参数说明和发布材料中的错误命令名称，程序只调用 `hy-smi` | 文档检查待 CI 确认 |
| 2 | K100_AI 上 `hy-smi --showuniqueid` 有 UUID，但 `hy_smi_info` 的 `unique_id`、`health` 为空 | 主查询缺少字段时，新增 `--showuniqueid`、`--healthcheck` 补充查询；兼容 `Unique ID`、`UniqueID`、`UUID` 和多种健康状态字段名；补充结果按设备序号合并 | 待 K100_AI 复测 |
| 3 | 部分指标统计过去 1 秒数据，8 卡机型执行耗时较长 | 新增 `--hy-smi.collect-1s-stats` 和 `HY_SMI_EXPORTER_COLLECT_1S_STATS`；默认值为 `false`，默认不输出利用率及 PCIe、XHCL、DF、UMC 带宽统计 | 待 8 卡环境对比耗时 |

必需主查询成功时，`hy_smi_up=1`。唯一标识、健康状态等可选补充查询失败会记录调试日志，但不会把整个设备采集误报为不可用。

## 2 K100_AI 定向复测

使用候选二进制在 K100_AI / `6.3.31-V1.5.1` 节点启动：

```bash
./hy-smi-exporter \
  --hy-smi.path=/opt/hyhal/bin/hy-smi \
  --hy-smi.timeout=15s \
  --hy-smi.collect-1s-stats=false
```

核对原始命令和 Exporter 输出：

```bash
hy-smi --showuniqueid
hy-smi --healthcheck
curl -s http://127.0.0.1:9584/metrics | grep '^hy_smi_info'
curl -s http://127.0.0.1:9584/metrics | grep '^hy_smi_health_info'
```

通过标准：

- `hy_smi_info` 的 `unique_id` 与 `hy-smi --showuniqueid` 对应设备一致且非空；
- `hy_smi_info` 的 `health` 非空，并与 `hy-smi --healthcheck` 一致；
- `hy_smi_health_info` 正常输出；
- `hy_smi_up=1`，设备数量与物理卡数量一致；
- 多卡设备的设备序号与 UUID、健康状态不存在错位。

## 3 1 秒统计开关复测

默认关闭模式下，确认利用率及 PCIe、XHCL、DF、UMC 带宽指标不输出，并记录连续 10 次 `/metrics` 请求耗时。

然后使用以下参数重启：

```bash
./hy-smi-exporter \
  --hy-smi.path=/opt/hyhal/bin/hy-smi \
  --hy-smi.timeout=30s \
  --hy-smi.collect-1s-stats=true
```

确认上述可选指标恢复输出，再记录连续 10 次请求的最小、最大和平均耗时。8 卡节点的 Prometheus `scrape_timeout` 应大于开启统计后的最大实际耗时，并保留安全余量。

## 4 反馈材料

请返回以下脱敏材料：

- 候选二进制版本和 SHA-256；
- 设备型号、驱动版本、Linux 和内核版本；
- `--showuniqueid`、`--healthcheck` 脱敏输出；
- 默认关闭和显式开启两种模式的 `/metrics` 关键行；
- 两种模式各 10 次请求的耗时统计；
- 复测结论、测试人和日期。

序列号、UUID、主机名、IP 地址和 PCIe 地址对外传递前应脱敏。复测通过后，再更新兼容性矩阵、RC 验收记录和正式发布放行确认单。
