# 驱动与 Exporter 安装指南

本文档说明海光设备节点上安装驱动、检查运行环境和部署 `hy-smi-exporter` 的基本流程。

## 1 安装前确认

安装前记录以下信息：

- 设备型号和数量；
- Linux 发行版和版本；
- 当前内核版本；
- CPU 架构；
- 计划安装的驱动、固件和 DTK 版本；
- 计划使用的 `hy-smi` / `hcu-smi` 版本；
- 首版部署方式：直接运行或 systemd。Docker、Docker Compose 和 Kubernetes
  不属于 `v0.1.0` 兼容范围。

驱动、DTK 和硬件的参考对应关系见 [compatibility.md](compatibility.md)。

## 2 安装系统依赖

### 2.1 CentOS、Rocky 或 RHEL 类系统

```bash
yum install -y rpm-build gcc-c++ cmake automake \
  elfutils-libelf-devel libdrm libdrm-devel pciutils

yum install -y kernel-devel-$(uname -r) kernel-modules-extra-$(uname -r)
```

### 2.2 Ubuntu 类系统

```bash
apt install -y gcc g++ cmake automake libelf-dev \
  libdrm_amdgpu1 libtinfo5 pciutils libdrm-dev

apt install -y linux-headers-$(uname -r) linux-modules-extra-$(uname -r)
```

包名可能因发行版版本不同而变化，正式安装前应以目标系统的软件源为准。

安装前置条件：

- `kernel` 与 `kernel-devel` / `linux-headers` 版本匹配；
- 新版本驱动的 CMake 版本高于 3.2；
- 新版本驱动的 GCC 版本大于 7.3；
- 具备安装内核模块和重启驱动服务的权限。

## 3 安装海光驱动

驱动安装包应从海光官方或双方确认的授权渠道获取。不要将驱动安装包、固件或厂商工具二进制提交到 exporter 仓库。

```bash
# 给安装包增加执行权限
chmod +x rock*.run

# 安装驱动；非虚拟化场景是否追加 -A 按海光建议执行
./rock*.run

# 重启驱动管理服务
systemctl restart hymgr
```

可通过以下命令查看安装参数：

```bash
./rock*.run -h
```

### 3.1 安装注意事项

1. 通常不需要先卸载旧驱动，安装新驱动时安装程序会处理旧版本。
2. 如果出现 `Do you want to remove legacy hymgr config file? [Y/n]:`，选择 `yes` 会重新生成 `/etc/hymgr.cfg`，选择 `no` 会保留原配置。
3. 非 Z100L 环境首次安装可能出现 `Update vbios? [y/N]:`。是否更新 VBIOS 以及是否重启物理机，应根据海光建议和现场运维窗口决定。
4. 云计算场景如果不允许重启物理节点，应提前确认不更新 VBIOS 的风险和替代方案。

## 4 检查驱动是否加载

海光资料指出，不同驱动版本的内核模块名称可能不同。可以先执行：

```bash
lsmod | grep -E 'hydcu|hycu'
```

同时确认设备节点和工具：

```bash
ls -l /dev/kfd /dev/dri 2>/dev/null
which hy-smi || which hcu-smi
hy-smi --version 2>/dev/null || hcu-smi --version
```

在正式测试前执行只读查询：

```bash
hy-smi -a --json
hy-smi -a
```

如果实际工具名称为 `hcu-smi`，将命令替换为对应工具；也可以通过 exporter 参数指定工具路径。

## 5 安装和运行 Exporter

### 5.1 直接运行

```bash
./bin/hy-smi-exporter \
  --web.listen-address=:9584 \
  --hy-smi.path=/usr/bin/hy-smi \
  --hy-smi.timeout=5s \
  --log.level=info
```

检查服务：

```bash
curl http://127.0.0.1:9584/healthz
curl http://127.0.0.1:9584/metrics | grep hy_smi_
```

### 5.2 systemd

仓库提供 `examples/hy-smi-exporter.service` 和
`examples/hy-smi-exporter.env`。确认宿主机已安装驱动及
`hy-smi` / `hcu-smi` 后，将二进制安装到目标路径并按发行版注册服务。

### 5.3 首版部署边界

Docker、Docker Compose、Kubernetes 和容器镜像不属于 `v0.1.0` 的兼容、
测试或生产部署范围，仓库不提供相关部署文件。

## 6 驱动升级后的处理

升级驱动、固件、DTK 或 `hy-smi` / `hcu-smi` 后，不应直接认定 exporter 继续兼容。应按以下顺序处理：

1. 记录升级前后的版本信息；
2. 执行 `hy-smi -a --json` 和 `hy-smi -a` 并保存脱敏样例；
3. 对照 [compatibility.md](compatibility.md) 检查字段、单位和枚举；
4. 执行厂商测试包中的冒烟测试和完整测试清单；
5. 必要时补充解析 fixture、单元测试和兼容性矩阵；
6. 只有测试通过后，才将该组合标记为“已验证”。

## 7 常见问题

| 现象 | 排查方法 |
| --- | --- |
| `hy-smi` 找不到 | 检查 PATH，或使用 `--hy-smi.path` 指定绝对路径 |
| `hy_smi_up` 为 0 | 检查日志、驱动、工具路径、权限和命令超时 |
| `hy_smi_up` 为 1 但设备数为 0 | 检查驱动模块、设备节点和厂商工具输出 |
| JSON 查询失败 | 检查工具版本，确认文本回退是否可用，并保留脱敏输出 |
| 指标缺少 | 对照工具原始输出和兼容矩阵，确认字段是否由当前驱动提供 |
| 驱动安装失败 | 检查内核头文件、`kernel-devel`、GCC、CMake 和安装日志 |
| 需要进程指标 | 当前 exporter 尚未提供进程指标，不能仅通过挂载 debugfs 获得 |

## 8 安全要求

- exporter 只调用查询类命令，不执行设备配置修改、驱动加载、设备重置或 RAS 注入；
- 不要把驱动、固件或厂商工具二进制放入公开仓库和测试包；
- 生产环境限制 `/metrics` 访问范围；
- 对外提交日志和样例前，脱敏序列号、Unique ID、主机名、IP 地址、PCIe 地址和客户信息；
- 驱动安装、VBIOS 更新和节点重启应在经过批准的维护窗口内执行。
