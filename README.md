# DuplicateHunter 🔍

一个高效的重复文件检测和清理工具，帮助您快速找到并管理系统中的重复文件，释放宝贵的存储空间。

![Python](https://img.shields.io/badge/Python-3.8+-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)

## ✨ 特性

- 🚀 **高性能扫描**：基于 MD5/SHA256 算法的快速文件比较
- 🌐 **Web 界面**：直观的 Web 管理界面，支持可视化操作
- 📁 **多格式支持**：支持所有文件类型的重复检测
- 🔄 **批量操作**：支持批量删除、移动和备份重复文件
- 💾 **内存优化**：智能内存管理，支持大文件和大量文件扫描
- 🐳 **容器化部署**：提供 Docker 镜像，一键部署
- 📊 **详细报告**：生成详细的扫描报告和统计信息
- 🔒 **安全可靠**：文件操作前自动备份，支持撤销操作
- 🎯 **智能差异显示**：重复文件名差异高亮，快速识别细微差别

## 🚀 快速开始

### 方式一：Docker 部署（推荐）

```bash
# 克隆项目
git clone https://github.com/cilidm/duplicate_file.git
cd duplicate_file

# 构建镜像
docker build -t duplicatehunter:latest .

# 运行容器
docker run -d \
  --name duplicate-hunter \
  -p 8080:8080 \
  -v /path/to/scan:/app/data \
  duplicatehunter:latest

# 访问 Web 界面
open http://localhost:8080
```

### 方式二：本地安装

```bash
# 克隆项目
git clone https://github.com/cilidm/duplicate_file.git
cd duplicate_file

# 安装依赖
pip install -r requirements.txt

# 启动 Web 服务
python start.py web

# 或使用命令行工具
python start.py scan /path/to/directory
```

## 📖 使用指南

### Web 界面使用

1. **选择扫描目录**：在 Web 界面中选择要扫描的目录
2. **配置扫描选项**：设置文件大小限制、文件类型过滤等
3. **开始扫描**：点击开始按钮，实时查看扫描进度
4. **查看结果**：浏览重复文件列表，智能差异高亮显示
5. **执行操作**：选择删除、移动或备份重复文件

### 命令行使用

```bash
# 基本扫描
python cli.py --scan /home/user/Documents

# 指定算法和输出格式
python cli.py --scan /data --algorithm sha256 --output json

# 批量删除重复文件（保留最新的）
python cli.py --scan /data --auto-delete --keep newest

# 生成详细报告
python cli.py --scan /data --report --output-file report.html
```

## 🔧 配置选项

### 扫描配置

```yaml
# config.yaml
scan:
  algorithm: "md5"  # md5, sha1, sha256
  min_size: 1024    # 最小文件大小（字节）
  max_size: null    # 最大文件大小（字节）
  extensions:       # 文件扩展名过滤
    - ".jpg"
    - ".png"
    - ".mp4"
  exclude_dirs:     # 排除目录
    - ".git"
    - "node_modules"
    - "__pycache__"

performance:
  threads: 4        # 扫描线程数
  chunk_size: 8192  # 文件读取块大小
  memory_limit: "1GB"  # 内存使用限制

web:
  host: "0.0.0.0"
  port: 8080
  debug: false
```

## 📊 性能基准

| 文件数量 | 总大小 | 扫描时间 | 内存使用 |
|---------|--------|----------|----------|
| 10,000  | 1GB    | 30s      | 50MB     |
| 100,000 | 10GB   | 5min     | 200MB    |
| 1,000,000 | 100GB | 45min   | 800MB    |

*测试环境：Intel i7-8700K, 16GB RAM, SSD*

## 🏗️ 架构设计

```
DuplicateHunter/
├── app/
│   ├── __init__.py
│   ├── scanner.py          # 核心扫描引擎
│   ├── hasher.py          # 文件哈希计算
│   ├── web/               # Web 界面
│   │   ├── app.py
│   │   └── templates/
│   └── utils.py           # 工具函数
├── tests/                 # 测试文件
├── examples/              # 使用示例
├── start.py               # 快速启动脚本
├── cli.py                 # 命令行工具
└── config.yaml            # 配置文件
```

## 🔍 算法说明

### 文件比较策略

1. **快速预筛选**：首先比较文件大小
2. **哈希计算**：对相同大小的文件计算哈希值
3. **分块读取**：大文件采用分块读取，避免内存溢出
4. **智能缓存**：缓存已计算的哈希值，提高重复扫描效率

### 支持的哈希算法

- **MD5**：速度快，适合快速扫描
- **SHA1**：平衡速度和安全性
- **SHA256**：最高安全性，适合重要文件

## 🐳 Docker 部署

### 使用 Docker Compose

```bash
# 克隆项目
git clone https://github.com/cilidm/duplicate_file.git
cd duplicate_file

# 使用 Docker Compose 启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## 📱 Web 界面特性

- **直观的扫描配置**：支持目录选择、算法选择、文件过滤等
- **实时进度显示**：扫描过程中显示实时进度和统计信息
- **智能结果展示**：重复文件差异高亮显示，便于快速识别
- **批量操作支持**：支持批量删除和文件管理
- **多格式报告**：支持 HTML、JSON、CSV 格式报告导出

## 🔒 安全特性

- **操作确认**：所有删除操作需要用户确认
- **自动备份**：删除前自动创建备份
- **操作日志**：记录所有文件操作历史
- **权限检查**：验证文件操作权限
- **撤销功能**：支持撤销最近的操作

## 🤝 贡献指南

我们欢迎所有形式的贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详细信息。

### 开发环境搭建

```bash
# 克隆项目
git clone https://github.com/cilidm/duplicate_file.git
cd duplicate_file

# 方式一：使用设置脚本（推荐）
python setup_dev.py

# 方式二：手动安装
pip install -r requirements.txt
pip install pytest black flake8 isort

# 运行测试
pytest tests/

# 启动开发服务器
python start.py web
```

## 📋 待办事项

- [ ] 支持网络驱动器扫描
- [ ] 添加文件预览功能
- [ ] 实现增量扫描
- [ ] 支持云存储集成
- [ ] 添加 API 接口
- [ ] 移动端适配

## 🆚 与其他工具对比

| 特性 | DuplicateHunter | dupeGuru | Duplicate Cleaner |
|------|----------------|----------|-------------------|
| Web 界面 | ✅ | ❌ | ❌ |
| 命令行 | ✅ | ✅ | ❌ |
| Docker 支持 | ✅ | ❌ | ❌ |
| 批量操作 | ✅ | ✅ | ✅ |
| 多算法支持 | ✅ | ✅ | ❌ |
| 智能差异显示 | ✅ | ❌ | ❌ |
| 开源免费 | ✅ | ✅ | ❌ |

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

感谢所有贡献者和用户的支持！

## 📞 联系我们

- 🐛 Issues: [GitHub Issues](https://github.com/cilidm/duplicate_file/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/cilidm/duplicate_file/discussions)
- 📧 Email: 通过 GitHub Issues 联系

---

⭐ 如果这个项目对您有帮助，请给我们一个 Star！