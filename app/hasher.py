"""
文件哈希计算器
"""

import hashlib
import os
from typing import Optional


class FileHasher:
    """文件哈希计算器"""
    
    SUPPORTED_ALGORITHMS = {
        'md5': hashlib.md5,
        'sha1': hashlib.sha1,
        'sha256': hashlib.sha256
    }
    
    def __init__(self, algorithm: str = 'md5', chunk_size: int = 8192):
        """初始化哈希计算器
        
        Args:
            algorithm: 哈希算法 (md5, sha1, sha256)
            chunk_size: 文件读取块大小
        """
        if algorithm not in self.SUPPORTED_ALGORITHMS:
            raise ValueError(f"不支持的哈希算法: {algorithm}")
            
        self.algorithm = algorithm
        self.chunk_size = chunk_size
        self._hash_cache = {}  # 哈希值缓存
        
    def calculate_hash(self, file_path: str) -> Optional[str]:
        """计算文件的哈希值
        
        Args:
            file_path: 文件路径
            
        Returns:
            文件的哈希值，如果计算失败返回 None
        """
        try:
            # 检查缓存
            file_stat = os.stat(file_path)
            cache_key = f"{file_path}:{file_stat.st_mtime}:{file_stat.st_size}"
            
            if cache_key in self._hash_cache:
                return self._hash_cache[cache_key]
                
            # 计算哈希值
            hash_obj = self.SUPPORTED_ALGORITHMS[self.algorithm]()
            
            with open(file_path, 'rb') as f:
                while chunk := f.read(self.chunk_size):
                    hash_obj.update(chunk)
                    
            file_hash = hash_obj.hexdigest()
            
            # 缓存结果
            self._hash_cache[cache_key] = file_hash
            
            return file_hash
            
        except (OSError, IOError) as e:
            print(f"无法读取文件 {file_path}: {e}")
            return None
        except Exception as e:
            print(f"计算哈希值时出错 {file_path}: {e}")
            return None
            
    def calculate_partial_hash(self, file_path: str, max_bytes: int = 1024 * 1024) -> Optional[str]:
        """计算文件的部分哈希值（用于快速比较大文件）
        
        Args:
            file_path: 文件路径
            max_bytes: 最大读取字节数
            
        Returns:
            文件的部分哈希值
        """
        try:
            hash_obj = self.SUPPORTED_ALGORITHMS[self.algorithm]()
            bytes_read = 0
            
            with open(file_path, 'rb') as f:
                while bytes_read < max_bytes:
                    chunk_size = min(self.chunk_size, max_bytes - bytes_read)
                    chunk = f.read(chunk_size)
                    
                    if not chunk:
                        break
                        
                    hash_obj.update(chunk)
                    bytes_read += len(chunk)
                    
            return hash_obj.hexdigest()
            
        except (OSError, IOError) as e:
            print(f"无法读取文件 {file_path}: {e}")
            return None
        except Exception as e:
            print(f"计算部分哈希值时出错 {file_path}: {e}")
            return None
            
    def clear_cache(self):
        """清空哈希值缓存"""
        self._hash_cache.clear()
        
    def get_cache_size(self) -> int:
        """获取缓存大小"""
        return len(self._hash_cache)
        
    def get_algorithm_info(self) -> dict:
        """获取算法信息"""
        return {
            'algorithm': self.algorithm,
            'chunk_size': self.chunk_size,
            'cache_size': self.get_cache_size(),
            'supported_algorithms': list(self.SUPPORTED_ALGORITHMS.keys())
        }