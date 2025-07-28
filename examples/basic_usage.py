#!/usr/bin/env python3
"""
DuplicateHunter 基本使用示例
"""

import os
import sys
import tempfile
from pathlib import Path

# 添加项目根目录到 Python 路径
sys.path.insert(0, str(Path(__file__).parent.parent))

from app.scanner import FileScanner, ScanConfig
from app.utils import FileUtils, ReportGenerator


def create_sample_files():
    """创建示例文件用于测试"""
    temp_dir = tempfile.mkdtemp(prefix='duplicate_test_')
    print(f"创建测试目录: {temp_dir}")
    
    # 创建一些重复文件
    files_data = [
        ('file1.txt', '这是重复内容'),
        ('file2.txt', '这是重复内容'),  # 与 file1.txt 重复
        ('file3.txt', '这是不同的内容'),
        ('subdir/file4.txt', '这是重复内容'),  # 与 file1.txt 重复
        ('subdir/file5.txt', '另一个不同的内容'),
        ('image1.jpg', '假装这是图片内容'),
        ('image2.jpg', '假装这是图片内容'),  # 与 image1.jpg 重复
    ]
    
    for file_path, content in files_data:
        full_path = os.path.join(temp_dir, file_path)
        os.makedirs(os.path.dirname(full_path), exist_ok=True)
        
        with open(full_path, 'w', encoding='utf-8') as f:
            f.write(content)
            
    return temp_dir


def basic_scan_example():
    """基本扫描示例"""
    print("=" * 60)
    print("基本扫描示例")
    print("=" * 60)
    
    # 创建测试文件
    test_dir = create_sample_files()
    
    try:
        # 创建扫描配置
        config = ScanConfig(
            algorithm='md5',
            min_size=1,  # 最小文件大小 1 字节
            threads=2    # 使用 2 个线程
        )
        
        # 创建扫描器
        scanner = FileScanner(config)
        
        # 设置进度回调
        def progress_callback(current, total, message):
            print(f"进度: {current}/{total} - {message}")
            
        scanner.set_progress_callback(progress_callback)
        
        # 执行扫描
        print(f"开始扫描目录: {test_dir}")
        result = scanner.scan_directory(test_dir)
        
        # 显示结果
        print(f"\n扫描完成！")
        print(f"总文件数: {result.total_files}")
        print(f"总大小: {FileUtils.format_size(result.total_size)}")
        print(f"扫描耗时: {result.scan_time:.2f} 秒")
        print(f"发现 {len(result.duplicate_groups)} 组重复文件")
        
        # 显示重复文件详情
        if result.duplicate_groups:
            print("\n重复文件详情:")
            group_id = 1
            for hash_value, files in result.duplicate_groups.items():
                print(f"\n重复组 #{group_id}:")
                print(f"  哈希值: {hash_value}")
                print(f"  文件数: {len(files)}")
                for file_path in files:
                    file_info = FileUtils.get_file_info(file_path)
                    print(f"    - {file_path} ({file_info.get('size_formatted', 'N/A')})")
                group_id += 1
                
        # 获取统计信息
        stats = scanner.get_scan_statistics(result)
        print(f"\n统计信息:")
        for key, value in stats.items():
            print(f"  {key}: {value}")
            
    finally:
        # 清理测试文件
        import shutil
        shutil.rmtree(test_dir, ignore_errors=True)
        print(f"\n已清理测试目录: {test_dir}")


def advanced_scan_example():
    """高级扫描示例"""
    print("\n" + "=" * 60)
    print("高级扫描示例")
    print("=" * 60)
    
    # 创建测试文件
    test_dir = create_sample_files()
    
    try:
        # 创建高级扫描配置
        config = ScanConfig(
            algorithm='sha256',  # 使用更安全的 SHA256
            min_size=10,         # 最小文件大小 10 字节
            extensions={'.txt'}, # 只扫描 .txt 文件
            threads=4            # 使用 4 个线程
        )
        
        # 创建扫描器
        scanner = FileScanner(config)
        
        # 执行扫描
        print(f"开始高级扫描目录: {test_dir}")
        print("配置: SHA256 算法，只扫描 .txt 文件，最小 10 字节")
        
        result = scanner.scan_directory(test_dir)
        
        # 显示结果
        print(f"\n扫描完成！")
        print(f"符合条件的文件数: {result.total_files}")
        print(f"发现 {len(result.duplicate_groups)} 组重复文件")
        
        # 生成报告
        if result.duplicate_groups:
            print("\n生成报告...")
            
            # 生成 HTML 报告
            html_file = os.path.join(test_dir, 'report.html')
            success = ReportGenerator.generate_html_report(result.duplicate_groups, html_file)
            if success:
                print(f"HTML 报告已生成: {html_file}")
                
            # 生成 JSON 报告
            json_file = os.path.join(test_dir, 'report.json')
            success = ReportGenerator.generate_json_report(result.duplicate_groups, json_file)
            if success:
                print(f"JSON 报告已生成: {json_file}")
                
    finally:
        # 清理测试文件
        import shutil
        shutil.rmtree(test_dir, ignore_errors=True)
        print(f"\n已清理测试目录: {test_dir}")


def file_operations_example():
    """文件操作示例"""
    print("\n" + "=" * 60)
    print("文件操作示例")
    print("=" * 60)
    
    # 创建测试文件
    test_dir = create_sample_files()
    backup_dir = os.path.join(test_dir, 'backups')
    
    try:
        # 创建一个测试文件
        test_file = os.path.join(test_dir, 'test_delete.txt')
        with open(test_file, 'w') as f:
            f.write('这个文件将被删除')
            
        print(f"创建测试文件: {test_file}")
        
        # 安全删除文件（带备份）
        print(f"安全删除文件到备份目录: {backup_dir}")
        success = FileUtils.safe_delete_file(test_file, backup_dir)
        
        if success:
            print("文件删除成功！")
            print(f"原文件已删除: {not os.path.exists(test_file)}")
            print(f"备份目录存在: {os.path.exists(backup_dir)}")
            
            # 列出备份文件
            if os.path.exists(backup_dir):
                backup_files = os.listdir(backup_dir)
                print(f"备份文件: {backup_files}")
        else:
            print("文件删除失败！")
            
        # 文件信息获取示例
        remaining_files = []
        for root, dirs, files in os.walk(test_dir):
            for file in files:
                if not file.startswith('.'):
                    file_path = os.path.join(root, file)
                    remaining_files.append(file_path)
                    
        if remaining_files:
            print(f"\n剩余文件信息:")
            for file_path in remaining_files[:3]:  # 只显示前3个
                info = FileUtils.get_file_info(file_path)
                print(f"  文件: {info.get('name', 'N/A')}")
                print(f"    大小: {info.get('size_formatted', 'N/A')}")
                print(f"    修改时间: {info.get('modified_time', 'N/A')}")
                print()
                
    finally:
        # 清理测试文件
        import shutil
        shutil.rmtree(test_dir, ignore_errors=True)
        print(f"已清理测试目录: {test_dir}")


def main():
    """主函数"""
    print("DuplicateHunter 使用示例")
    print("=" * 60)
    
    try:
        # 运行示例
        basic_scan_example()
        advanced_scan_example()
        file_operations_example()
        
        print("\n" + "=" * 60)
        print("所有示例运行完成！")
        print("=" * 60)
        
    except KeyboardInterrupt:
        print("\n示例被用户中断")
    except Exception as e:
        print(f"\n运行示例时出错: {e}")
        import traceback
        traceback.print_exc()


if __name__ == '__main__':
    main()