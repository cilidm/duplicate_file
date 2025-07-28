"""
DuplicateHunter - 重复文件检测和清理工具

一个高效的重复文件检测工具，支持 Web 界面和命令行操作。
"""

__version__ = "1.0.0"
__author__ = "cilidm"
__description__ = "高效的重复文件检测和清理工具"

from .scanner import FileScanner
from .hasher import FileHasher
from .utils import FileUtils

__all__ = ['FileScanner', 'FileHasher', 'FileUtils']