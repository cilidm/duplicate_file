#!/usr/bin/env python3
"""
DuplicateHunter 命令行工具
"""

import argparse
import sys
import os
import time
from pathlib import Path
from app.scanner import FileScanner, ScanConfig
from app.utils import FileUtils, ReportGenerator


def main():
    parser = argparse.ArgumentParser(
        description='DuplicateHunter - 重复文件检测和清理工具',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
使用示例:
  %(prog)s --scan /home/user/Documents
  %(prog)s --scan /data --algorithm sha256 --output json
  %(prog)s --scan /data --auto-delete --keep newest
  %(prog)s --scan /data --report --output-file report.html
        """
    )
    
    # 基本参数
    parser.add_argument('--scan', required=True, help='要扫描的目录路径')
    parser.add_argument('--algorithm', choices=['md5', 'sha1', 'sha256'], 
                       default='md5', help='哈希算法 (默认: md5)')
    parser.add_argument('--min-size', type=int, default=1024, 
                       help='最小文件大小（字节，默认: 1024）')
    parser.add_argument('--max-size', type=int, help='最大文件大小（字节）')
    parser.add_argument('--threads', type=int, default=4, 
                       help='扫描线程数 (默认: 4)')
    
    # 过滤参数
    parser.add_argument('--extensions', help='文件扩展名过滤，用逗号分隔 (如: .jpg,.png,.mp4)')
    parser.add_argument('--exclude-dirs', help='排除的目录，用逗号分隔')
    
    # 输出参数
    parser.add_argument('--output', choices=['text', 'json', 'csv'], 
                       default='text', help='输出格式 (默认: text)')
    parser.add_argument('--output-file', help='输出文件路径')
    parser.add_argument('--report', action='store_true', help='生成详细报告')
    
    # 操作参数
    parser.add_argument('--auto-delete', action='store_true', 
                       help='自动删除重复文件')
    parser.add_argument('--keep', choices=['oldest', 'newest', 'first'], 
                       default='first', help='保留哪个文件 (默认: first)')
    parser.add_argument('--backup-dir', help='删除前备份目录')
    parser.add_argument('--dry-run', action='store_true', 
                       help='只显示将要执行的操作，不实际执行')
    
    # 其他参数
    parser.add_argument('--verbose', '-v', action='store_true', help='详细输出')
    parser.add_argument('--quiet', '-q', action='store_true', help='静默模式')
    
    args = parser.parse_args()
    
    # 验证参数
    if not os.path.exists(args.scan):
        print(f"错误: 目录不存在 - {args.scan}", file=sys.stderr)
        sys.exit(1)
        
    if args.auto_delete and not args.dry_run and not args.backup_dir:
        print("警告: 建议在自动删除时指定备份目录 (--backup-dir)")
        
    # 创建扫描配置
    config = ScanConfig(
        algorithm=args.algorithm,
        min_size=args.min_size,
        max_size=args.max_size,
        threads=args.threads
    )
    
    # 处理扩展名过滤
    if args.extensions:
        extensions = [ext.strip() for ext in args.extensions.split(',')]
        config.extensions = set(ext if ext.startswith('.') else f'.{ext}' for ext in extensions)
        
    # 处理排除目录
    if args.exclude_dirs:
        exclude_dirs = [d.strip() for d in args.exclude_dirs.split(',')]
        config.exclude_dirs.update(exclude_dirs)
        
    # 创建扫描器
    scanner = FileScanner(config)
    
    # 设置进度回调
    if not args.quiet:
        def progress_callback(current, total, message):
            if args.verbose:
                print(f"\r进度: {current}/{total} - {message}", end='', flush=True)
            else:
                percentage = (current / total * 100) if total > 0 else 0
                print(f"\r进度: {percentage:.1f}%", end='', flush=True)
                
        scanner.set_progress_callback(progress_callback)
        
    try:
        # 执行扫描
        if not args.quiet:
            print(f"开始扫描目录: {args.scan}")
            print(f"使用算法: {args.algorithm}")
            print(f"线程数: {args.threads}")
            print("-" * 50)
            
        result = scanner.scan_directory(args.scan)
        
        if not args.quiet:
            print()  # 换行
            
        # 显示统计信息
        stats = scanner.get_scan_statistics(result)
        if not args.quiet:
            print("\n扫描统计:")
            for key, value in stats.items():
                print(f"  {key}: {value}")
            print("-" * 50)
            
        # 处理结果
        if not result.duplicate_groups:
            print("太棒了！没有发现重复文件。")
            sys.exit(0)
            
        # 输出结果
        if args.output == 'text':
            display_text_results(result.duplicate_groups, args.verbose)
        elif args.output == 'json':
            output_json_results(result.duplicate_groups, args.output_file)
        elif args.output == 'csv':
            output_csv_results(result.duplicate_groups, args.output_file)
            
        # 生成报告
        if args.report:
            generate_report(result.duplicate_groups, args.output_file)
            
        # 自动删除
        if args.auto_delete:
            auto_delete_duplicates(result.duplicate_groups, args)
            
    except KeyboardInterrupt:
        print("\n扫描被用户中断")
        sys.exit(1)
    except Exception as e:
        print(f"扫描失败: {e}", file=sys.stderr)
        sys.exit(1)


def display_text_results(duplicate_groups, verbose=False):
    """显示文本格式结果"""
    print(f"\n发现 {len(duplicate_groups)} 组重复文件:")
    print("=" * 60)
    
    group_id = 1
    for hash_value, files in duplicate_groups.items():
        print(f"\n重复组 #{group_id} ({len(files)} 个文件):")
        if verbose:
            print(f"  哈希值: {hash_value}")
            
        for i, file_path in enumerate(files):
            file_info = FileUtils.get_file_info(file_path)
            size_str = file_info.get('size_formatted', 'N/A')
            print(f"  [{i+1}] {file_path} ({size_str})")
            
        group_id += 1


def output_json_results(duplicate_groups, output_file=None):
    """输出 JSON 格式结果"""
    if not output_file:
        output_file = f"duplicate_results_{int(time.time())}.json"
        
    success = ReportGenerator.generate_json_report(duplicate_groups, output_file)
    if success:
        print(f"JSON 结果已保存到: {output_file}")
    else:
        print("保存 JSON 结果失败", file=sys.stderr)


def output_csv_results(duplicate_groups, output_file=None):
    """输出 CSV 格式结果"""
    if not output_file:
        output_file = f"duplicate_results_{int(time.time())}.csv"
        
    success = ReportGenerator.generate_csv_report(duplicate_groups, output_file)
    if success:
        print(f"CSV 结果已保存到: {output_file}")
    else:
        print("保存 CSV 结果失败", file=sys.stderr)


def generate_report(duplicate_groups, output_file=None):
    """生成详细报告"""
    if not output_file:
        output_file = f"duplicate_report_{int(time.time())}.html"
    elif not output_file.endswith('.html'):
        output_file += '.html'
        
    success = ReportGenerator.generate_html_report(duplicate_groups, output_file)
    if success:
        print(f"详细报告已生成: {output_file}")
    else:
        print("生成报告失败", file=sys.stderr)


def auto_delete_duplicates(duplicate_groups, args):
    """自动删除重复文件"""
    import time
    
    total_files = sum(len(files) for files in duplicate_groups.values())
    files_to_delete = []
    
    # 确定要删除的文件
    for files in duplicate_groups.values():
        if len(files) <= 1:
            continue
            
        # 根据保留策略选择要保留的文件
        if args.keep == 'oldest':
            # 保留最旧的文件
            files_with_time = [(f, os.path.getmtime(f)) for f in files]
            files_with_time.sort(key=lambda x: x[1])
            keep_file = files_with_time[0][0]
        elif args.keep == 'newest':
            # 保留最新的文件
            files_with_time = [(f, os.path.getmtime(f)) for f in files]
            files_with_time.sort(key=lambda x: x[1], reverse=True)
            keep_file = files_with_time[0][0]
        else:  # first
            # 保留第一个文件
            keep_file = files[0]
            
        # 添加其他文件到删除列表
        for file_path in files:
            if file_path != keep_file:
                files_to_delete.append(file_path)
                
    if not files_to_delete:
        print("没有文件需要删除")
        return
        
    print(f"\n将要删除 {len(files_to_delete)} 个重复文件:")
    
    if args.dry_run:
        print("(模拟运行 - 不会实际删除文件)")
        for file_path in files_to_delete:
            print(f"  [删除] {file_path}")
        return
        
    # 确认删除
    if not args.quiet:
        response = input(f"确定要删除这 {len(files_to_delete)} 个文件吗? (y/N): ")
        if response.lower() not in ['y', 'yes']:
            print("操作已取消")
            return
            
    # 执行删除
    success_count = 0
    for file_path in files_to_delete:
        try:
            success = FileUtils.safe_delete_file(file_path, args.backup_dir)
            if success:
                success_count += 1
                if not args.quiet:
                    print(f"  [已删除] {file_path}")
            else:
                print(f"  [失败] {file_path}", file=sys.stderr)
        except Exception as e:
            print(f"  [错误] {file_path}: {e}", file=sys.stderr)
            
    print(f"\n删除完成: {success_count}/{len(files_to_delete)} 个文件")
    
    if args.backup_dir and success_count > 0:
        print(f"备份目录: {args.backup_dir}")


if __name__ == '__main__':
    main()