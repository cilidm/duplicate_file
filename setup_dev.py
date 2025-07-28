#!/usr/bin/env python3
"""
开发环境设置脚本
"""

import subprocess
import sys
import os


def run_command(cmd, description):
    """运行命令并处理错误"""
    print(f"正在{description}...")
    try:
        result = subprocess.run(cmd, shell=True, check=True, capture_output=True, text=True)
        print(f"✅ {description}成功")
        return True
    except subprocess.CalledProcessError as e:
        print(f"❌ {description}失败: {e}")
        if e.stdout:
            print(f"输出: {e.stdout}")
        if e.stderr:
            print(f"错误: {e.stderr}")
        return False


def install_basic_deps():
    """安装基础依赖"""
    print("=" * 50)
    print("安装基础依赖")
    print("=" * 50)
    
    # 安装生产依赖
    if not run_command("pip install -r requirements.txt", "安装生产依赖"):
        return False
        
    # 安装核心开发依赖
    core_deps = [
        "pytest>=7.0.0",
        "pytest-cov>=4.0.0", 
        "black>=23.0.0",
        "flake8>=6.0.0",
        "isort>=5.0.0"
    ]
    
    for dep in core_deps:
        if not run_command(f"pip install '{dep}'", f"安装 {dep.split('>=')[0]}"):
            print(f"⚠️  跳过 {dep}")
            
    return True


def install_optional_deps():
    """安装可选依赖"""
    print("\n" + "=" * 50)
    print("安装可选依赖")
    print("=" * 50)
    
    optional_deps = [
        ("mypy>=1.0.0", "类型检查"),
        ("pre-commit>=3.0.0", "Git 钩子"),
        ("memory-profiler>=0.60.0", "内存分析"),
        ("sphinx>=7.0.0", "文档生成"),
        ("sphinx-rtd-theme>=1.0.0", "文档主题")
    ]
    
    for dep, desc in optional_deps:
        if not run_command(f"pip install '{dep}'", f"安装 {desc}"):
            print(f"⚠️  跳过可选依赖: {desc}")


def setup_pre_commit():
    """设置 pre-commit"""
    print("\n" + "=" * 50)
    print("设置 pre-commit")
    print("=" * 50)
    
    if os.path.exists(".pre-commit-config.yaml"):
        run_command("pre-commit install", "设置 pre-commit 钩子")
    else:
        print("⚠️  未找到 .pre-commit-config.yaml，跳过 pre-commit 设置")


def create_directories():
    """创建必要的目录"""
    print("\n" + "=" * 50)
    print("创建项目目录")
    print("=" * 50)
    
    dirs = ["reports", "backups", "logs", "data"]
    for dir_name in dirs:
        os.makedirs(dir_name, exist_ok=True)
        print(f"✅ 创建目录: {dir_name}")


def main():
    """主函数"""
    print("DuplicateHunter 开发环境设置")
    print("=" * 60)
    
    # 检查 Python 版本
    if sys.version_info < (3, 8):
        print("❌ 需要 Python 3.8 或更高版本")
        sys.exit(1)
        
    print(f"✅ Python 版本: {sys.version}")
    
    # 安装依赖
    if not install_basic_deps():
        print("❌ 基础依赖安装失败")
        sys.exit(1)
        
    install_optional_deps()
    setup_pre_commit()
    create_directories()
    
    print("\n" + "=" * 60)
    print("🎉 开发环境设置完成！")
    print("=" * 60)
    print("\n可用命令:")
    print("  python app/web/app.py          # 启动 Web 服务")
    print("  python cli.py --scan /path     # 运行命令行工具")
    print("  pytest tests/                  # 运行测试")
    print("  black app/ cli.py              # 代码格式化")
    print("  flake8 app/ cli.py             # 代码检查")
    print("  python examples/basic_usage.py # 运行示例")
    print("\n或使用 Makefile:")
    print("  make help                      # 查看所有命令")
    print("  make run-web                   # 启动 Web 服务")
    print("  make test                      # 运行测试")


if __name__ == "__main__":
    main()