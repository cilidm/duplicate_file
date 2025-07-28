# DuplicateHunter 项目结构

```
DuplicateHunter/
├── 📄 README.md                    # 项目说明文档
├── 🚀 start.py                     # 快速启动脚本
├── 🐍 cli.py                       # 命令行工具
├── ⚙️ config.yaml                  # 配置文件
├── 🛠️ setup_dev.py                 # 开发环境设置
├── 🔧 Makefile                     # 自动化工具
├── 📋 requirements.txt             # 生产依赖
├── 📋 requirements-dev.txt         # 开发依赖
├── 🐳 Dockerfile                   # Docker 镜像
├── 🐳 docker-compose.yml           # Docker Compose
├── 🔒 .pre-commit-config.yaml      # 代码质量检查
├── 🚫 .gitignore                   # Git 忽略文件
├── 📝 CHANGELOG.md                 # 更新日志
├── 🤝 CONTRIBUTING.md              # 贡献指南
├── 📄 LICENSE                      # MIT 许可证
├── 📦 app/                         # 核心应用
│   ├── __init__.py                 # 包初始化
│   ├── scanner.py                  # 扫描引擎
│   ├── hasher.py                   # 哈希计算
│   ├── utils.py                    # 工具函数
│   └── web/                        # Web 界面
│       ├── app.py                  # Flask 应用
│       └── templates/
│           └── index.html          # 主页模板
├── 🧪 tests/                       # 测试文件
│   ├── __init__.py
│   ├── test_scanner.py
│   └── test_hasher.py
└── 📚 examples/                    # 使用示例
    ├── basic_usage.py              # 基本用法
    └── test_web.py                 # Web 测试
```

## 核心模块说明

### 📦 app/
- **scanner.py**: 核心扫描引擎，支持多线程文件扫描
- **hasher.py**: 文件哈希计算，支持 MD5/SHA1/SHA256
- **utils.py**: 工具函数集合，包含文件操作、报告生成等
- **web/app.py**: Flask Web 应用，提供 REST API 和 Web 界面

### 🚀 快速启动
- **start.py**: 统一启动脚本，支持 Web 和 CLI 模式
- **cli.py**: 命令行工具，支持批量操作和报告生成

### 🛠️ 开发工具
- **setup_dev.py**: 自动化开发环境设置
- **Makefile**: 常用开发任务自动化
- **.pre-commit-config.yaml**: 代码质量检查配置

### 🐳 部署配置
- **Dockerfile**: Docker 镜像构建配置
- **docker-compose.yml**: 容器编排配置
- **config.yaml**: 应用配置文件

## 使用方法

### Web 模式
```bash
python start.py web
# 访问 http://localhost:8080
```

### 命令行模式
```bash
python start.py scan /path/to/directory
```

### Docker 部署
```bash
docker-compose up -d
```

### 开发环境
```bash
python setup_dev.py
make help