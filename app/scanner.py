"""
文件扫描器 - 核心扫描引擎
"""

import os
import time
import threading
from pathlib import Path
from typing import Dict, List, Set, Optional, Callable
from concurrent.futures import ThreadPoolExecutor, as_completed
from dataclasses import dataclass
from .hasher import FileHasher
from .utils import FileUtils


@dataclass
class ScanResult:
    """扫描结果数据类"""
    total_files: int = 0
    total_size: int = 0
    duplicate_groups: Dict[str, List[str]] = None
    scan_time: float = 0.0
    errors: List[str] = None
    
    def __post_init__(self):
        if self.duplicate_groups is None:
            self.duplicate_groups = {}
        if self.errors is None:
            self.errors = []


@dataclass
class ScanConfig:
    """扫描配置"""
    algorithm: str = "md5"  # md5, sha1, sha256
    min_size: int = 1024  # 最小文件大小（字节）
    max_size: Optional[int] = None  # 最大文件大小（字节）
    extensions: Optional[Set[str]] = None  # 允许的文件扩展名
    exclude_dirs: Set[str] = None  # 排除的目录
    threads: int = 4  # 扫描线程数
    chunk_size: int = 8192  # 文件读取块大小
    
    def __post_init__(self):
        if self.exclude_dirs is None:
            self.exclude_dirs = {'.git', 'node_modules', '__pycache__', '.vscode', '.idea'}


class FileScanner:
    """文件扫描器"""
    
    def __init__(self, config: ScanConfig = None):
        self.config = config or ScanConfig()
        self.hasher = FileHasher(self.config.algorithm, self.config.chunk_size)
        self._stop_event = threading.Event()
        self._progress_callback: Optional[Callable] = None
        
    def set_progress_callback(self, callback: Callable[[int, int, str], None]):
        """设置进度回调函数
        
        Args:
            callback: 回调函数，参数为 (当前进度, 总数, 当前文件)
        """
        self._progress_callback = callback
        
    def stop_scan(self):
        """停止扫描"""
        self._stop_event.set()
        
    def scan_directory(self, directory: str) -> ScanResult:
        """扫描目录中的重复文件
        
        Args:
            directory: 要扫描的目录路径
            
        Returns:
            ScanResult: 扫描结果
        """
        start_time = time.time()
        result = ScanResult()
        
        try:
            # 收集所有文件
            files = self._collect_files(directory)
            if not files:
                return result
                
            result.total_files = len(files)
            result.total_size = sum(os.path.getsize(f) for f in files)
            
            # 按文件大小分组
            size_groups = self._group_by_size(files)
            
            # 计算哈希值并找出重复文件
            result.duplicate_groups = self._find_duplicates(size_groups)
            
        except Exception as e:
            result.errors.append(f"扫描错误: {str(e)}")
            
        result.scan_time = time.time() - start_time
        return result
        
    def _collect_files(self, directory: str) -> List[str]:
        """收集目录中的所有文件"""
        files = []
        directory_path = Path(directory)
        
        if not directory_path.exists():
            raise FileNotFoundError(f"目录不存在: {directory}")
            
        for root, dirs, filenames in os.walk(directory):
            # 排除指定目录
            dirs[:] = [d for d in dirs if d not in self.config.exclude_dirs]
            
            for filename in filenames:
                if self._stop_event.is_set():
                    break
                    
                file_path = os.path.join(root, filename)
                
                # 检查文件是否符合条件
                if self._should_include_file(file_path):
                    files.append(file_path)
                    
        return files
        
    def _should_include_file(self, file_path: str) -> bool:
        """检查文件是否应该包含在扫描中"""
        try:
            # 检查文件是否存在且可读
            if not os.path.isfile(file_path) or not os.access(file_path, os.R_OK):
                return False
                
            # 检查文件大小
            file_size = os.path.getsize(file_path)
            if file_size < self.config.min_size:
                return False
                
            if self.config.max_size and file_size > self.config.max_size:
                return False
                
            # 检查文件扩展名
            if self.config.extensions:
                file_ext = Path(file_path).suffix.lower()
                if file_ext not in self.config.extensions:
                    return False
                    
            return True
            
        except (OSError, IOError):
            return False
            
    def _group_by_size(self, files: List[str]) -> Dict[int, List[str]]:
        """按文件大小分组"""
        size_groups = {}
        
        for file_path in files:
            if self._stop_event.is_set():
                break
                
            try:
                file_size = os.path.getsize(file_path)
                if file_size not in size_groups:
                    size_groups[file_size] = []
                size_groups[file_size].append(file_path)
            except (OSError, IOError):
                continue
                
        # 只返回有多个文件的组
        return {size: files for size, files in size_groups.items() if len(files) > 1}
        
    def _find_duplicates(self, size_groups: Dict[int, List[str]]) -> Dict[str, List[str]]:
        """在相同大小的文件中找出重复文件"""
        duplicate_groups = {}
        total_files = sum(len(files) for files in size_groups.values())
        processed_files = 0
        
        with ThreadPoolExecutor(max_workers=self.config.threads) as executor:
            # 为每个大小组提交哈希计算任务
            future_to_files = {}
            
            for size, files in size_groups.items():
                if self._stop_event.is_set():
                    break
                    
                if len(files) > 1:  # 只处理有多个文件的组
                    future = executor.submit(self._process_size_group, files)
                    future_to_files[future] = files
                    
            # 收集结果
            for future in as_completed(future_to_files):
                if self._stop_event.is_set():
                    break
                    
                try:
                    group_duplicates = future.result()
                    duplicate_groups.update(group_duplicates)
                    
                    processed_files += len(future_to_files[future])
                    
                    # 更新进度
                    if self._progress_callback:
                        self._progress_callback(processed_files, total_files, "计算文件哈希值...")
                        
                except Exception as e:
                    print(f"处理文件组时出错: {e}")
                    
        return duplicate_groups
        
    def _process_size_group(self, files: List[str]) -> Dict[str, List[str]]:
        """处理相同大小的文件组"""
        hash_groups = {}
        
        for file_path in files:
            if self._stop_event.is_set():
                break
                
            try:
                file_hash = self.hasher.calculate_hash(file_path)
                if file_hash:
                    if file_hash not in hash_groups:
                        hash_groups[file_hash] = []
                    hash_groups[file_hash].append(file_path)
                    
            except Exception as e:
                print(f"计算文件哈希值失败 {file_path}: {e}")
                
        # 只返回有重复的组
        return {hash_val: files for hash_val, files in hash_groups.items() if len(files) > 1}
        
    def get_scan_statistics(self, result: ScanResult) -> Dict:
        """获取扫描统计信息"""
        duplicate_files = []
        duplicate_size = 0
        
        for files in result.duplicate_groups.values():
            duplicate_files.extend(files)
            # 计算重复文件占用的空间（除了保留一个文件）
            if files:
                file_size = os.path.getsize(files[0])
                duplicate_size += file_size * (len(files) - 1)
                
        return {
            "总文件数": result.total_files,
            "总大小": FileUtils.format_size(result.total_size),
            "重复文件组数": len(result.duplicate_groups),
            "重复文件数": len(duplicate_files),
            "可释放空间": FileUtils.format_size(duplicate_size),
            "扫描耗时": f"{result.scan_time:.2f}秒",
            "错误数": len(result.errors)
        }