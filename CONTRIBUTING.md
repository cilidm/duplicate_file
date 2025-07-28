# 贡献指南

感谢您对 DuplicateHunter 项目的关注！我们欢迎所有形式的贡献。

## 如何贡献

### 报告问题

如果您发现了 bug 或有功能建议，请：

1. 检查 [Issues](https://github.com/cilidm/DuplicateHunter/issues) 确认问题未被报告
2. 创建新的 Issue，详细描述问题或建议
3. 提供复现步骤、环境信息和相关日志

### 提交代码

1. **Fork 项目**
   ```bash
   git clone https://github.com/your-username/DuplicateHunter.git
   cd DuplicateHunter
   ```

2. **创建开发分支**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **设置开发环境**
   ```bash
   python -m venv venv
   source venv/bin/activate  # Linux/Mac
   # 或 venv\Scripts\activate  # Windows
   
   pip install -r requirements.txt
   pip install -r requirements-dev.txt
   ```

4. **进行开发**
   - 遵循现有的代码风格
   - 添加必要的测试
   - 更新相关文档

5. **运行测试**
   ```bash
   pytest tests/
   flake8 app/
   black app/ --check
   ```

6. **提交更改**
   ```bash
   git add .
   git commit -m "feat: 添加新功能描述"
   git push origin feature/your-feature-name
   ```

7. **创建 Pull Request**
   - 详细描述更改内容
   - 关联相关的 Issues
   - 确保所有检查通过

## 开发规范

### 代码风格

- 使用 [Black](https://black.readthedocs.io/) 进行代码格式化
- 使用 [Flake8](https://flake8.pycqa.org/) 进行代码检查
- 遵循 [PEP 8](https://www.python.org/dev/peps/pep-0008/) 规范

### 提交信息

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

- `feat:` 新功能
- `fix:` 修复 bug
- `docs:` 文档更新
- `style:` 代码格式调整
- `refactor:` 代码重构
- `test:` 测试相关
- `chore:` 构建或辅助工具的变动

### 测试

- 为新功能添加单元测试
- 确保测试覆盖率不低于 80%
- 运行完整测试套件确保无回归

## 开发环境

### 本地开发

```bash
# 启动 Web 服务
python app/web/app.py

# 运行命令行工具
python cli.py --scan /path/to/directory

# 运行测试
pytest tests/ -v

# 代码格式化
black app/ cli.py
flake8 app/ cli.py
```

### Docker 开发

```bash
# 构建镜像
docker build -t duplicatehunter:dev .

# 运行容器
docker-compose up -d

# 查看日志
docker-compose logs -f
```

## 项目结构

```
DuplicateHunter/
├── app/                    # 核心应用代码
│   ├── __init__.py
│   ├── scanner.py         # 扫描引擎
│   ├── hasher.py         # 哈希计算
│   ├── utils.py          # 工具函数
│   └── web/              # Web 界面
│       ├── app.py
│       └── templates/
├── tests/                # 测试文件
├── docs/                 # 文档
├── examples/             # 使用示例
├── cli.py               # 命令行工具
├── config.yaml          # 配置文件
├── requirements.txt     # 依赖列表
└── README.md           # 项目说明
```

## 发布流程

1. 更新版本号（`app/__init__.py`）
2. 更新 CHANGELOG.md
3. 创建 Git 标签
4. 构建和发布 Docker 镜像
5. 发布 GitHub Release

## 社区

- [GitHub Discussions](https://github.com/cilidm/DuplicateHunter/discussions) - 讨论和问答
- [Issues](https://github.com/cilidm/DuplicateHunter/issues) - Bug 报告和功能请求

## 许可证

通过贡献代码，您同意您的贡献将在 MIT 许可证下授权。

---

再次感谢您的贡献！🎉