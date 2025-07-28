#!/usr/bin/env python3
"""
Web 应用测试脚本
"""

import os
import sys
import tempfile
import shutil
from pathlib import Path

# 添加项目根目录到 Python 路径
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from app.web.app import create_app


def create_test_files():
    """创建测试文件"""
    test_dir = tempfile.mkdtemp(prefix='duplicate_test_')
    print(f"创建测试目录: {test_dir}")
    
    # 创建一些重复文件
    content1 = "这是测试文件内容1"
    content2 = "这是测试文件内容2"
    
    # 创建重复文件组1
    with open(os.path.join(test_dir, 'file1.txt'), 'w', encoding='utf-8') as f:
        f.write(content1)
    with open(os.path.join(test_dir, 'file1_copy.txt'), 'w', encoding='utf-8') as f:
        f.write(content1)
        
    # 创建重复文件组2
    subdir = os.path.join(test_dir, 'subdir')
    os.makedirs(subdir)
    with open(os.path.join(subdir, 'file2.txt'), 'w', encoding='utf-8') as f:
        f.write(content2)
    with open(os.path.join(subdir, 'file2_copy.txt'), 'w', encoding='utf-8') as f:
        f.write(content2)
        
    # 创建唯一文件
    with open(os.path.join(test_dir, 'unique.txt'), 'w', encoding='utf-8') as f:
        f.write("这是唯一的文件内容")
        
    return test_dir


def cleanup_test_files(test_dir):
    """清理测试文件"""
    if os.path.exists(test_dir):
        shutil.rmtree(test_dir)
        print(f"清理测试目录: {test_dir}")


def main():
    """主函数"""
    print("DuplicateHunter Web 应用测试")
    print("=" * 50)
    
    # 创建测试文件
    test_dir = create_test_files()
    
    try:
        # 创建 Flask 应用
        app = create_app()
        
        print(f"\n🚀 启动 Web 服务...")
        print(f"访问地址: http://localhost:8080")
        print(f"测试目录: {test_dir}")
        print(f"按 Ctrl+C 停止服务")
        print("-" * 50)
        
        # 启动应用
        app.run(host='0.0.0.0', port=8080, debug=True)
        
    except KeyboardInterrupt:
        print("\n👋 服务已停止")
    except Exception as e:
        print(f"❌ 启动失败: {e}")
    finally:
        # 清理测试文件
        cleanup_test_files(test_dir)


if __name__ == '__main__':
    main()