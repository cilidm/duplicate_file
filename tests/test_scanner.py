"""
文件扫描器测试
"""

import os
import tempfile
import pytest
from pathlib import Path
from app.scanner import FileScanner, ScanConfig


class TestFileScanner:
    """文件扫描器测试类"""
    
    def setup_method(self):
        """测试前准备"""
        self.temp_dir = tempfile.mkdtemp()
        self.config = ScanConfig(algorithm='md5', min_size=1)
        self.scanner = FileScanner(self.config)
        
    def teardown_method(self):
        """测试后清理"""
        import shutil
        shutil.rmtree(self.temp_dir, ignore_errors=True)
        
    def create_test_file(self, filename, content):
        """创建测试文件"""
        file_path = os.path.join(self.temp_dir, filename)
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content)
        return file_path
        
    def test_scan_empty_directory(self):
        """测试扫描空目录"""
        result = self.scanner.scan_directory(self.temp_dir)
        assert result.total_files == 0
        assert len(result.duplicate_groups) == 0
        
    def test_scan_no_duplicates(self):
        """测试扫描无重复文件"""
        self.create_test_file('file1.txt', 'content1')
        self.create_test_file('file2.txt', 'content2')
        
        result = self.scanner.scan_directory(self.temp_dir)
        assert result.total_files == 2
        assert len(result.duplicate_groups) == 0
        
    def test_scan_with_duplicates(self):
        """测试扫描有重复文件"""
        self.create_test_file('file1.txt', 'same content')
        self.create_test_file('file2.txt', 'same content')
        self.create_test_file('file3.txt', 'different content')
        
        result = self.scanner.scan_directory(self.temp_dir)
        assert result.total_files == 3
        assert len(result.duplicate_groups) == 1
        
        # 检查重复组
        duplicate_group = list(result.duplicate_groups.values())[0]
        assert len(duplicate_group) == 2
        
    def test_min_size_filter(self):
        """测试最小文件大小过滤"""
        self.config.min_size = 10
        scanner = FileScanner(self.config)
        
        self.create_test_file('small.txt', 'small')  # 5 字节
        self.create_test_file('large.txt', 'large content here')  # > 10 字节
        
        result = scanner.scan_directory(self.temp_dir)
        assert result.total_files == 1  # 只有大文件被扫描
        
    def test_extension_filter(self):
        """测试文件扩展名过滤"""
        self.config.extensions = {'.txt'}
        scanner = FileScanner(self.config)
        
        self.create_test_file('file1.txt', 'content')
        self.create_test_file('file2.py', 'content')
        
        result = scanner.scan_directory(self.temp_dir)
        assert result.total_files == 1  # 只有 .txt 文件被扫描
        
    def test_exclude_dirs(self):
        """测试排除目录"""
        # 创建子目录
        sub_dir = os.path.join(self.temp_dir, '.git')
        os.makedirs(sub_dir)
        
        self.create_test_file('file1.txt', 'content')
        self.create_test_file('.git/config', 'git config')
        
        result = self.scanner.scan_directory(self.temp_dir)
        assert result.total_files == 1  # .git 目录被排除