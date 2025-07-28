"""
文件哈希计算器测试
"""

import os
import tempfile
import pytest
from app.hasher import FileHasher


class TestFileHasher:
    """文件哈希计算器测试类"""
    
    def setup_method(self):
        """测试前准备"""
        self.temp_dir = tempfile.mkdtemp()
        self.hasher = FileHasher('md5')
        
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
        
    def test_md5_hash(self):
        """测试 MD5 哈希计算"""
        file_path = self.create_test_file('test.txt', 'hello world')
        hash_value = self.hasher.calculate_hash(file_path)
        
        # "hello world" 的 MD5 值
        expected = '5d41402abc4b2a76b9719d911017c592'
        assert hash_value == expected
        
    def test_sha256_hash(self):
        """测试 SHA256 哈希计算"""
        hasher = FileHasher('sha256')
        file_path = self.create_test_file('test.txt', 'hello world')
        hash_value = hasher.calculate_hash(file_path)
        
        # "hello world" 的 SHA256 值
        expected = 'b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9'
        assert hash_value == expected
        
    def test_same_content_same_hash(self):
        """测试相同内容产生相同哈希值"""
        file1 = self.create_test_file('file1.txt', 'same content')
        file2 = self.create_test_file('file2.txt', 'same content')
        
        hash1 = self.hasher.calculate_hash(file1)
        hash2 = self.hasher.calculate_hash(file2)
        
        assert hash1 == hash2
        
    def test_different_content_different_hash(self):
        """测试不同内容产生不同哈希值"""
        file1 = self.create_test_file('file1.txt', 'content1')
        file2 = self.create_test_file('file2.txt', 'content2')
        
        hash1 = self.hasher.calculate_hash(file1)
        hash2 = self.hasher.calculate_hash(file2)
        
        assert hash1 != hash2
        
    def test_hash_cache(self):
        """测试哈希值缓存"""
        file_path = self.create_test_file('test.txt', 'test content')
        
        # 第一次计算
        hash1 = self.hasher.calculate_hash(file_path)
        cache_size1 = self.hasher.get_cache_size()
        
        # 第二次计算（应该使用缓存）
        hash2 = self.hasher.calculate_hash(file_path)
        cache_size2 = self.hasher.get_cache_size()
        
        assert hash1 == hash2
        assert cache_size1 == cache_size2 == 1
        
    def test_nonexistent_file(self):
        """测试不存在的文件"""
        hash_value = self.hasher.calculate_hash('/nonexistent/file.txt')
        assert hash_value is None
        
    def test_unsupported_algorithm(self):
        """测试不支持的算法"""
        with pytest.raises(ValueError):
            FileHasher('unsupported_algorithm')