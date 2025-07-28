#!/usr/bin/env python3
"""
DuplicateHunter 快速启动脚本
"""

import sys
import os
import argparse
from pathlib import Path


def start_web_server():
    """启动 Web 服务器"""
    print("🚀 启动 DuplicateHunter Web 服务...")
    print("访问地址: http://localhost:8080")
    print("按 Ctrl+C 停止服务")
    print("-" * 50)
    
    try:
        from app.web.app import create_app
        app = create_app()
        app.run(host='0.0.0.0', port=8080, debug=False)
    except ImportError as e:
        print(f"❌ 导入错误: {e}")
        print("请先安装依赖: pip install -r requirements.txt")
        sys.exit(1)
    except Exception as e:
        print(f"❌ 启动失败: {e}")
        sys.exit(1)


def run_cli_scan(directory, algorithm='md5', output='text'):
    """运行命令行扫描"""
    if not os.path.exists(directory):
        print(f"❌ 目录不存在: {directory}")
        sys.exit(1)
        
    print(f"🔍 开始扫描目录: {directory}")
    print(f"使用算法: {algorithm}")
    print("-" * 50)
    
    try:
        from app.scanner import FileScanner, ScanConfig
        from app.utils import FileUtils
        
        # 创建配置
        config = ScanConfig(algorithm=algorithm, min_size=1024)
        scanner = FileScanner(config)
        
        # 设置进度回调
        def progress_callback(current, total, message):
            percentage = (current / total * 100) if total > 0 else 0
            print(f"\r进度: {percentage:.1f}% ({current}/{total})", end='', flush=True)
            
        scanner.set_progress_callback(progress_callback)
        
        # 执行扫描
        result = scanner.scan_directory(directory)
        print()  # 换行
        
        # 显示结果
        if not result.duplicate_groups:
            print("✅ 太棒了！没有发现重复文件。")
        else:
            print(f"📊 发现 {len(result.duplicate_groups)} 组重复文件:")
            
            group_id = 1
            for hash_value, files in result.duplicate_groups.items():
                print(f"\n重复组 #{group_id} ({len(files)} 个文件):")
                for file_path in files:
                    file_info = FileUtils.get_file_info(file_path)
                    size_str = file_info.get('size_formatted', 'N/A')
                    print(f"  - {file_path} ({size_str})")
                group_id += 1
                
        # 显示统计
        stats = scanner.get_scan_statistics(result)
        print(f"\n📈 统计信息:")
        for key, value in stats.items():
            print(f"  {key}: {value}")
            
    except ImportError as e:
        print(f"❌ 导入错误: {e}")
        print("请先安装依赖: pip install -r requirements.txt")
        sys.exit(1)
    except Exception as e:
        print(f"❌ 扫描失败: {e}")
        sys.exit(1)


def run_example():
    """运行示例"""
    print("🎯 运行 DuplicateHunter 示例...")
    print("-" * 50)
    
    try:
        import subprocess
        result = subprocess.run([sys.executable, "examples/basic_usage.py"], 
                              capture_output=False, text=True)
        if result.returncode != 0:
            print("❌ 示例运行失败")
            sys.exit(1)
    except Exception as e:
        print(f"❌ 运行示例失败: {e}")
        sys.exit(1)


def main():
    """主函数"""
    parser = argparse.ArgumentParser(
        description='DuplicateHunter 快速启动工具',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
使用示例:
  %(prog)s web                    # 启动 Web 服务
  %(prog)s scan /path/to/dir      # 扫描指定目录
  %(prog)s example                # 运行示例
        """
    )
    
    subparsers = parser.add_subparsers(dest='command', help='可用命令')
    
    # Web 服务命令
    subparsers.add_parser('web', help='启动 Web 服务')
    
    # 扫描命令
    scan_parser = subparsers.add_parser('scan', help='扫描目录')
    scan_parser.add_argument('directory', help='要扫描的目录')
    scan_parser.add_argument('--algorithm', choices=['md5', 'sha1', 'sha256'], 
                           default='md5', help='哈希算法')
    scan_parser.add_argument('--output', choices=['text', 'json'], 
                           default='text', help='输出格式')
    
    # 示例命令
    subparsers.add_parser('example', help='运行示例')
    
    args = parser.parse_args()
    
    if not args.command:
        parser.print_help()
        return
        
    print("DuplicateHunter - 重复文件检测工具")
    print("=" * 50)
    
    try:
        if args.command == 'web':
            start_web_server()
        elif args.command == 'scan':
            run_cli_scan(args.directory, args.algorithm, args.output)
        elif args.command == 'example':
            run_example()
    except KeyboardInterrupt:
        print("\n👋 再见！")
    except Exception as e:
        print(f"❌ 运行失败: {e}")
        sys.exit(1)


if __name__ == '__main__':
    main()