# 海光驱动兼容参考及安装要点

本文件整理自海光提供的《驱动兼容列表及安装方法.md》，用于发布前技术评审和首版兼容性矩阵的输入。该表是厂商驱动参考，不等同于 `hy-smi-exporter` 已完成的实测兼容承诺。公开发布前请由海光确认版本、型号和安装说明是否可以公开。

## 1 驱动兼容参考

| 驱动版本 | 支持硬件 | 推荐 DTK 版本 |
| --- | --- | --- |
| `rock-4.0-23-5.6.25` | Z100 / Z100L | 21.04 / 21.10 / 22.04 |
| `rock-4.5.2-5.11.40` | Z100 / Z100L | 22.04 / 22.10 / 23.04 |
| `rock-5.2.0-5.16.18` | Z100 / Z100L | 22.10 / 23.04 |
| `rock-5.2.0-5.16.29` | Z100 / Z100L / K100 | 23.04 / 23.10 |
| `rock-5.7.1-6.2.26+` | Z100 / Z100L / K100 / K100-AI | 24.04 / 25.04 |
| `rock-6.3.8+` | BW1000 / K100-AI / K100 / Z100L / Z100 | 25.04* |
| `rock-6.3.30+` | BW1100 / BW1000 / BW150 / BW100 / BW10 / K100-AI / K100 / Z100L / Z100 | 26.04 |

`25.04*` 的 DTK 对应关系需要海光在公开发布前确认具体含义和适用条件。

## 2 安装前置依赖

CentOS/RHEL 类系统：

```bash
yum install -y rpm-build gcc-c++ cmake automake \
  elfutils-libelf-devel libdrm libdrm-devel pciutils
yum install -y kernel-devel-$(uname -r) kernel-modules-extra-$(uname -r)
```

Ubuntu 类系统：

```bash
apt install -y gcc g++ cmake automake libelf-dev \
  libdrm_amdgpu1 libtinfo5 pciutils libdrm-dev
apt install -y linux-headers-$(uname -r) linux-modules-extra-$(uname -r)
```

包名可能因发行版而变化，以目标系统软件源和海光安装包要求为准。安装前确认：

- `kernel` 与 `kernel-devel` / `linux-headers` 版本一致；
- 新版本驱动要求 CMake 高于 3.2；
- 新版本驱动要求 GCC 高于 7.3；
- 安装和重启驱动服务具备相应权限。

## 3 驱动安装流程

驱动安装包应从海光官方或双方确认的授权渠道获取，不放入 exporter 仓库或公开发布包。

```bash
chmod +x rock*.run
./rock*.run
systemctl restart hymgr
```

非虚拟化场景是否增加 `-A` 参数，按海光建议执行。安装程序参数可用 `./rock*.run -h` 查看。安装新驱动通常不需要手工卸载旧驱动；安装过程中会处理旧版本。

安装过程中：

- `Do you want to remove legacy hymgr config file? [Y/n]:` 选择 `yes` 会重新生成 `/etc/hymgr.cfg`，选择 `no` 会保留现有配置；
- 非 Z100L 环境可能提示 `Update vbios? [y/N]:`。首次安装是否更新 VBIOS 和是否重启物理机，必须结合海光建议和维护窗口决定；
- 云计算环境如不允许重启，应提前确认不更新 VBIOS 的风险。

## 4 驱动模块检查

6.2.x 及以下版本通常显示 `hydcu` 模块，6.3.x 及以上版本通常显示 `hycu` 模块：

```bash
lsmod | grep -E 'hydcu|hycu'
```

模块名称变化不能单独作为 exporter 兼容结论。测试时还要记录 `hy-smi` / `hcu-smi` 版本、设备型号、固件、DTK、Linux 发行版和内核版本。

## 5 与 `hy-smi-exporter` 的关系

- 该参考表只说明海光驱动资料列出的硬件和 DTK 关系；
- exporter 首版只支持宿主机直接运行或 systemd；
- Docker、Docker Compose、Kubernetes 和容器镜像不属于 `v0.1.0` 兼容范围；
- 只有在真实设备上完成 JSON 主查询、文本回退、指标对照和稳定性测试的组合，才能写入 `docs/compatibility.md` 的“已验证”部分；
- 驱动、固件、DTK 或工具升级后必须重新执行回归。

海光公开确认人：`TBD`  
确认日期：`TBD`  
公开范围：`TBD`
