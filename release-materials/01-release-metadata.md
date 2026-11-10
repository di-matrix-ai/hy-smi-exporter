# 发布基础信息确认表

## 1 项目和仓库

| 项目 | 最终值 | 责任方 | 状态 |
| --- | --- | --- | --- |
| 项目名 | `hy-smi-exporter` | 双方 | 已确认 |
| 中文名 | 海光 DCU 性能数据 Exporter | 双方 | 待最终确认 |
| GitHub 组织 | `di-matrix-ai` | 鼎茂 | 已确认 |
| 公开仓库 URL | `https://github.com/di-matrix-ai/hy-smi-exporter` | 鼎茂 | 已确认 |
| Go module | `github.com/di-matrix-ai/hy-smi-exporter` | 鼎茂 | 已确认 |
| GHCR 镜像地址 | 首版不发布容器镜像 | 鼎茂 | 已确定 |
| 默认分支 | `main` | 鼎茂 | 已确定 |
| 首版 tag | `v0.1.0` | 鼎茂 | 已确定 |
| 计划发布日期 | 2026-11-13 | 双方 | 待最终确认 |

Go module、源码 import、README、镜像地址、Release URL 和 Prometheus 收录材料必须使用同一 GitHub 组织名。

## 2 法律和开源信息

| 项目 | 最终值 | 责任方 | 状态 |
| --- | --- | --- | --- |
| 鼎茂法定中文全称 | 上海鼎茂信息技术有限公司 | 鼎茂 | 已确认 |
| 鼎茂法定英文全称 | Shanghai Ding Mao Information Technology Co., Ltd. | 鼎茂 | 已确认 |
| 版权年份 | `2026` | 鼎茂 | 待确认 |
| 开源协议 | Apache License 2.0 | 鼎茂 | 已确定 |
| `NOTICE` 最终文本 | `Copyright 2026 Shanghai Ding Mao Information Technology Co., Ltd.` | 鼎茂法务/合规 | 版权主体已确认，待法务最终复核 |
| 第三方依赖审查人 | `TBD` | 鼎茂 | 待填写 |
| 海光公开授权留档编号 | `TBD` | 双方 | 待填写 |

## 3 维护和联系人

| 角色 | 姓名 | GitHub 登录名 | 邮箱/内部联系方式 | 责任方 |
| --- | --- | --- | --- | --- |
| 项目 Owner | 鼎茂发布团队 | `@di-matrix-ai` | `esteric_shen@di-matrix.com` | 鼎茂 |
| 主维护者 | 鼎茂发布团队 | `@di-matrix-ai` | `esteric_shen@di-matrix.com` | 鼎茂 |
| 备份维护者 | 待确认 | 待确认 | 待确认 | 鼎茂 |
| Release 负责人 | Esteric Shen | `@di-matrix-ai`（组织账号） | `esteric_shen@di-matrix.com` | 鼎茂 |
| 安全联系人 | Dingmao security contact | `@di-matrix-ai` | `esteric_shen@di-matrix.com` | 鼎茂 |
| 合规负责人 | `TBD` | `TBD` | `TBD` | 鼎茂 |
| 海光技术协作者 | 王忠勤 | `@zqwangadv` | `zqwangadv@163.com` | 海光 |
| 海光技术协作者 | 代金龙 | `@kinglong08` | `ise_longlong@163.com` | 海光 |
| 海光备份联系人 | 代金龙 | `@kinglong08` | `ise_longlong@163.com` | 海光 |
| 问题升级联系人 | 待确认 | 不适用 | 待确认 | 双方 |

## 4 GitHub 权限建议

| 权限 | 建议配置 |
| --- | --- |
| Dingmao Owner | `di-matrix-ai` 组织；至少 2 个公司管理员账号，避免单账号失联 |
| 主维护者 | Maintain 或 Write |
| 海光协作者 | Triage 或 Write，按实际参与方式确定 |
| 分支保护 | `main` 要求 PR、CI 通过和至少 1 人审阅 |
| Release | 仅鼎茂 Release 负责人和备份负责人可创建正式 Release |
| 安全报告 | 启用 GitHub Private Vulnerability Reporting |

确认人：`TBD`  
确认日期：`TBD`
