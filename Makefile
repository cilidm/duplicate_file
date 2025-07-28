# DuplicateHunter Makefile

.PHONY: help install dev test lint format clean build docker run-web run-cli

# 默认目标
help:
	@echo "DuplicateHunter 开发工具"
	@echo ""
	@echo "可用命令:"
	@echo "  install     安装依赖"
	@echo "  dev         安装开发依赖"
	@echo "  test        运行测试"
	@echo "  lint        代码检查"
	@echo "  format      代码格式化"
	@echo "  clean       清理临时文件"
	@echo "  build       构建 Docker 镜像"
	@echo "  docker      运行 Docker 容器"
	@echo "  run-web     启动 Web 服务"
	@echo "  run-cli     运行命令行工具示例"

# 安装生产依赖
install:
	pip install -r requirements.txt

# 安装开发依赖
dev: install
	pip install -r requirements-dev.txt
	pre-commit install

# 运行测试
test:
	pytest tests/ -v --cov=app --cov-report=html --cov-report=term

# 代码检查
lint:
	flake8 app/ cli.py tests/
	mypy app/ cli.py

# 代码格式化
format:
	black app/ cli.py tests/ examples/
	isort app/ cli.py tests/ examples/

# 清理临时文件
clean:
	find . -type f -name "*.pyc" -delete
	find . -type d -name "__pycache__" -delete
	find . -type d -name "*.egg-info" -exec rm -rf {} +
	rm -rf build/ dist/ .coverage htmlcov/ .pytest_cache/
	rm -rf reports/ backups/ logs/

# 构建 Docker 镜像
build:
	docker build -t duplicatehunter:latest .

# 运行 Docker 容器
docker: build
	docker-compose up -d

# 停止 Docker 容器
docker-stop:
	docker-compose down

# 启动 Web 服务
run-web:
	python app/web/app.py

# 运行命令行工具示例
run-cli:
	python examples/basic_usage.py

# 运行命令行工具（需要指定目录）
cli:
	@echo "使用方法: make cli DIR=/path/to/directory"
	@if [ -z "$(DIR)" ]; then \
		echo "错误: 请指定扫描目录，例如: make cli DIR=/home/user/Documents"; \
		exit 1; \
	fi
	python cli.py --scan $(DIR) --verbose

# 生成文档
docs:
	@echo "生成文档..."
	@mkdir -p docs/
	@echo "文档生成功能待实现"

# 发布准备
release-check: test lint
	@echo "发布前检查完成"

# 创建发布
release: release-check
	@echo "创建发布..."
	@echo "请手动创建 Git 标签和 GitHub Release"

# 开发环境设置
setup-dev: dev
	@echo "开发环境设置完成"
	@echo "可以使用以下命令:"
	@echo "  make run-web    # 启动 Web 服务"
	@echo "  make test       # 运行测试"
	@echo "  make lint       # 代码检查"
	@echo "  make format     # 代码格式化"

# 性能测试
benchmark:
	@echo "运行性能测试..."
	python -m pytest tests/ -k "benchmark" -v

# 安全检查
security:
	@echo "运行安全检查..."
	@echo "安装 bandit 进行安全检查: pip install bandit"
	@echo "运行: bandit -r app/"

# 依赖更新
update-deps:
	@echo "依赖更新功能需要安装 pip-tools"
	@echo "运行: pip install pip-tools"
	@echo "然后: pip-compile requirements.in"

# 项目统计
stats:
	@echo "项目统计信息:"
	@echo "代码行数:"
	@find app/ -name "*.py" | xargs wc -l | tail -1
	@echo "测试行数:"
	@find tests/ -name "*.py" | xargs wc -l | tail -1
	@echo "文件数量:"
	@find . -name "*.py" | wc -l