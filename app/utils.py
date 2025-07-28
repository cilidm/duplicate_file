"""
工具函数模块
"""

import os
import shutil
import json
import csv
from datetime import datetime
from pathlib import Path
from typing import List, Dict, Any, Optional


class FileUtils:
    """文件操作工具类"""
    
    @staticmethod
    def format_size(size_bytes: int) -> str:
        """格式化文件大小
        
        Args:
            size_bytes: 字节数
            
        Returns:
            格式化后的大小字符串
        """
        if size_bytes == 0:
            return "0B"
            
        size_names = ["B", "KB", "MB", "GB", "TB"]
        i = 0
        
        while size_bytes >= 1024 and i < len(size_names) - 1:
            size_bytes /= 1024.0
            i += 1
            
        return f"{size_bytes:.2f}{size_names[i]}"
        
    @staticmethod
    def get_file_info(file_path: str) -> Dict[str, Any]:
        """获取文件信息
        
        Args:
            file_path: 文件路径
            
        Returns:
            文件信息字典
        """
        try:
            stat = os.stat(file_path)
            return {
                'path': file_path,
                'name': os.path.basename(file_path),
                'size': stat.st_size,
                'size_formatted': FileUtils.format_size(stat.st_size),
                'modified_time': datetime.fromtimestamp(stat.st_mtime).isoformat(),
                'created_time': datetime.fromtimestamp(stat.st_ctime).isoformat(),
                'extension': Path(file_path).suffix.lower(),
                'directory': os.path.dirname(file_path)
            }
        except (OSError, IOError):
            return {}
            
    @staticmethod
    def safe_delete_file(file_path: str, backup_dir: Optional[str] = None) -> bool:
        """安全删除文件（可选备份）
        
        Args:
            file_path: 要删除的文件路径
            backup_dir: 备份目录，如果提供则先备份再删除
            
        Returns:
            是否删除成功
        """
        try:
            if backup_dir and os.path.exists(file_path):
                # 创建备份
                backup_path = FileUtils.create_backup(file_path, backup_dir)
                if not backup_path:
                    return False
                    
            # 删除文件
            os.remove(file_path)
            return True
            
        except (OSError, IOError) as e:
            print(f"删除文件失败 {file_path}: {e}")
            return False
            
    @staticmethod
    def create_backup(file_path: str, backup_dir: str) -> Optional[str]:
        """创建文件备份
        
        Args:
            file_path: 源文件路径
            backup_dir: 备份目录
            
        Returns:
            备份文件路径，失败返回 None
        """
        try:
            if not os.path.exists(backup_dir):
                os.makedirs(backup_dir)
                
            # 生成备份文件名
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            file_name = os.path.basename(file_path)
            backup_name = f"{timestamp}_{file_name}"
            backup_path = os.path.join(backup_dir, backup_name)
            
            # 复制文件
            shutil.copy2(file_path, backup_path)
            return backup_path
            
        except (OSError, IOError) as e:
            print(f"创建备份失败 {file_path}: {e}")
            return None
            
    @staticmethod
    def move_file(src_path: str, dst_dir: str) -> bool:
        """移动文件到指定目录
        
        Args:
            src_path: 源文件路径
            dst_dir: 目标目录
            
        Returns:
            是否移动成功
        """
        try:
            if not os.path.exists(dst_dir):
                os.makedirs(dst_dir)
                
            file_name = os.path.basename(src_path)
            dst_path = os.path.join(dst_dir, file_name)
            
            # 如果目标文件已存在，添加序号
            counter = 1
            while os.path.exists(dst_path):
                name, ext = os.path.splitext(file_name)
                dst_path = os.path.join(dst_dir, f"{name}_{counter}{ext}")
                counter += 1
                
            shutil.move(src_path, dst_path)
            return True
            
        except (OSError, IOError) as e:
            print(f"移动文件失败 {src_path}: {e}")
            return False


class ReportGenerator:
    """报告生成器"""
    
    @staticmethod
    def generate_json_report(duplicate_groups: Dict[str, List[str]], output_file: str) -> bool:
        """生成 JSON 格式报告
        
        Args:
            duplicate_groups: 重复文件组
            output_file: 输出文件路径
            
        Returns:
            是否生成成功
        """
        try:
            report_data = {
                'generated_at': datetime.now().isoformat(),
                'total_groups': len(duplicate_groups),
                'duplicate_groups': []
            }
            
            for hash_value, files in duplicate_groups.items():
                group_data = {
                    'hash': hash_value,
                    'file_count': len(files),
                    'files': [FileUtils.get_file_info(f) for f in files]
                }
                report_data['duplicate_groups'].append(group_data)
                
            with open(output_file, 'w', encoding='utf-8') as f:
                json.dump(report_data, f, indent=2, ensure_ascii=False)
                
            return True
            
        except Exception as e:
            print(f"生成 JSON 报告失败: {e}")
            return False
            
    @staticmethod
    def generate_csv_report(duplicate_groups: Dict[str, List[str]], output_file: str) -> bool:
        """生成 CSV 格式报告
        
        Args:
            duplicate_groups: 重复文件组
            output_file: 输出文件路径
            
        Returns:
            是否生成成功
        """
        try:
            with open(output_file, 'w', newline='', encoding='utf-8') as f:
                writer = csv.writer(f)
                writer.writerow(['组ID', '哈希值', '文件路径', '文件大小', '修改时间'])
                
                group_id = 1
                for hash_value, files in duplicate_groups.items():
                    for file_path in files:
                        file_info = FileUtils.get_file_info(file_path)
                        writer.writerow([
                            group_id,
                            hash_value,
                            file_path,
                            file_info.get('size_formatted', ''),
                            file_info.get('modified_time', '')
                        ])
                    group_id += 1
                    
            return True
            
        except Exception as e:
            print(f"生成 CSV 报告失败: {e}")
            return False
            
    @staticmethod
    def generate_html_report(duplicate_groups: Dict[str, List[str]], output_file: str) -> bool:
        """生成 HTML 格式报告
        
        Args:
            duplicate_groups: 重复文件组
            output_file: 输出文件路径
            
        Returns:
            是否生成成功
        """
        try:
            html_content = f"""
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>重复文件扫描报告</title>
    <style>
        body {{ font-family: Arial, sans-serif; margin: 20px; }}
        .header {{ background-color: #f8f9fa; padding: 20px; border-radius: 5px; margin-bottom: 20px; }}
        .group {{ border: 1px solid #ddd; margin-bottom: 20px; border-radius: 5px; }}
        .group-header {{ background-color: #e9ecef; padding: 10px; font-weight: bold; }}
        .file-list {{ padding: 10px; }}
        .file-item {{ padding: 5px 0; border-bottom: 1px solid #eee; }}
        .file-path {{ font-family: monospace; }}
        .file-size {{ color: #666; }}
        .stats {{ display: flex; gap: 20px; }}
        .stat-item {{ text-align: center; }}
    </style>
</head>
<body>
    <div class="header">
        <h1>重复文件扫描报告</h1>
        <div class="stats">
            <div class="stat-item">
                <div><strong>{len(duplicate_groups)}</strong></div>
                <div>重复文件组</div>
            </div>
            <div class="stat-item">
                <div><strong>{sum(len(files) for files in duplicate_groups.values())}</strong></div>
                <div>重复文件数</div>
            </div>
        </div>
        <p>生成时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}</p>
    </div>
"""
            
            group_id = 1
            for hash_value, files in duplicate_groups.items():
                html_content += f"""
    <div class="group">
        <div class="group-header">
            重复组 #{group_id} - {len(files)} 个文件 (哈希: {hash_value[:16]}...)
        </div>
        <div class="file-list">
"""
                
                for file_path in files:
                    file_info = FileUtils.get_file_info(file_path)
                    html_content += f"""
            <div class="file-item">
                <div class="file-path">{file_path}</div>
                <div class="file-size">大小: {file_info.get('size_formatted', 'N/A')} | 修改时间: {file_info.get('modified_time', 'N/A')}</div>
            </div>
"""
                
                html_content += """
        </div>
    </div>
"""
                group_id += 1
                
            html_content += """
</body>
</html>
"""
            
            with open(output_file, 'w', encoding='utf-8') as f:
                f.write(html_content)
                
            return True
            
        except Exception as e:
            print(f"生成 HTML 报告失败: {e}")
            return False


class ConfigManager:
    """配置管理器"""
    
    @staticmethod
    def load_config(config_file: str) -> Dict[str, Any]:
        """加载配置文件
        
        Args:
            config_file: 配置文件路径
            
        Returns:
            配置字典
        """
        try:
            with open(config_file, 'r', encoding='utf-8') as f:
                if config_file.endswith('.json'):
                    return json.load(f)
                elif config_file.endswith('.yaml') or config_file.endswith('.yml'):
                    import yaml
                    return yaml.safe_load(f)
                else:
                    raise ValueError("不支持的配置文件格式")
                    
        except Exception as e:
            print(f"加载配置文件失败: {e}")
            return {}
            
    @staticmethod
    def save_config(config: Dict[str, Any], config_file: str) -> bool:
        """保存配置文件
        
        Args:
            config: 配置字典
            config_file: 配置文件路径
            
        Returns:
            是否保存成功
        """
        try:
            with open(config_file, 'w', encoding='utf-8') as f:
                if config_file.endswith('.json'):
                    json.dump(config, f, indent=2, ensure_ascii=False)
                elif config_file.endswith('.yaml') or config_file.endswith('.yml'):
                    import yaml
                    yaml.dump(config, f, default_flow_style=False, allow_unicode=True)
                else:
                    raise ValueError("不支持的配置文件格式")
                    
            return True
            
        except Exception as e:
            print(f"保存配置文件失败: {e}")
            return False