"""
Web 应用主程序
"""

from flask import Flask, render_template, request, jsonify, send_file
import os
import threading
import time
from datetime import datetime
from ..scanner import FileScanner, ScanConfig, ScanResult
from ..utils import FileUtils, ReportGenerator


app = Flask(__name__)
app.config['SECRET_KEY'] = 'duplicate-hunter-secret-key'

# 全局变量
current_scan = None
scan_thread = None
scan_results = {}


@app.route('/')
def index():
    """主页"""
    return render_template('index.html')


@app.route('/api/scan', methods=['POST'])
def start_scan():
    """开始扫描"""
    global current_scan, scan_thread
    
    data = request.get_json()
    directory = data.get('directory')
    
    if not directory or not os.path.exists(directory):
        return jsonify({'error': '目录不存在'}), 400
        
    # 检查是否已有扫描在进行
    if current_scan and scan_thread and scan_thread.is_alive():
        return jsonify({'error': '已有扫描在进行中'}), 400
        
    # 创建扫描配置
    config = ScanConfig(
        algorithm=data.get('algorithm', 'md5'),
        min_size=data.get('min_size', 1024),
        max_size=data.get('max_size'),
        threads=data.get('threads', 4)
    )
    
    # 处理文件扩展名过滤
    extensions = data.get('extensions', [])
    if extensions:
        config.extensions = set(ext.lower() for ext in extensions)
        
    # 处理排除目录
    exclude_dirs = data.get('exclude_dirs', [])
    if exclude_dirs:
        config.exclude_dirs.update(exclude_dirs)
        
    # 创建扫描器
    current_scan = FileScanner(config)
    scan_id = str(int(time.time()))
    
    # 启动扫描线程
    scan_thread = threading.Thread(
        target=run_scan,
        args=(current_scan, directory, scan_id)
    )
    scan_thread.start()
    
    return jsonify({
        'scan_id': scan_id,
        'message': '扫描已开始'
    })


@app.route('/api/scan/<scan_id>/status')
def get_scan_status(scan_id):
    """获取扫描状态"""
    if scan_id not in scan_results:
        return jsonify({'status': 'not_found'}), 404
        
    result = scan_results[scan_id]
    
    if result['status'] == 'running':
        return jsonify({
            'status': 'running',
            'progress': result.get('progress', {}),
            'message': result.get('message', '扫描中...')
        })
    elif result['status'] == 'completed':
        return jsonify({
            'status': 'completed',
            'result': result['data'],
            'statistics': result['statistics']
        })
    elif result['status'] == 'error':
        return jsonify({
            'status': 'error',
            'error': result['error']
        })
        
    return jsonify({'status': 'unknown'})


@app.route('/api/scan/<scan_id>/stop', methods=['POST'])
def stop_scan(scan_id):
    """停止扫描"""
    global current_scan
    
    if current_scan:
        current_scan.stop_scan()
        return jsonify({'message': '扫描已停止'})
        
    return jsonify({'error': '没有正在进行的扫描'}), 400


@app.route('/api/files/delete', methods=['POST'])
def delete_files():
    """删除文件"""
    data = request.get_json()
    files = data.get('files', [])
    backup_dir = data.get('backup_dir')
    
    results = []
    for file_path in files:
        success = FileUtils.safe_delete_file(file_path, backup_dir)
        results.append({
            'file': file_path,
            'success': success
        })
        
    return jsonify({'results': results})


@app.route('/api/files/move', methods=['POST'])
def move_files():
    """移动文件"""
    data = request.get_json()
    files = data.get('files', [])
    target_dir = data.get('target_dir')
    
    if not target_dir:
        return jsonify({'error': '目标目录不能为空'}), 400
        
    results = []
    for file_path in files:
        success = FileUtils.move_file(file_path, target_dir)
        results.append({
            'file': file_path,
            'success': success
        })
        
    return jsonify({'results': results})


@app.route('/api/report/generate', methods=['POST'])
def generate_report():
    """生成报告"""
    try:
        data = request.get_json()
        scan_id = data.get('scan_id')
        format_type = data.get('format', 'html')
        
        if scan_id not in scan_results or scan_results[scan_id]['status'] != 'completed':
            return jsonify({'error': '扫描结果不存在'}), 400
            
        duplicate_groups = scan_results[scan_id]['data']['duplicate_groups']
        timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
        
        # 确保报告目录存在
        reports_dir = os.path.abspath('reports')
        os.makedirs(reports_dir, exist_ok=True)
        
        filename = f'duplicate_report_{timestamp}.{format_type}'
        filepath = os.path.join(reports_dir, filename)
        
        success = False
        if format_type == 'json':
            success = ReportGenerator.generate_json_report(duplicate_groups, filepath)
        elif format_type == 'csv':
            success = ReportGenerator.generate_csv_report(duplicate_groups, filepath)
        elif format_type == 'html':
            success = ReportGenerator.generate_html_report(duplicate_groups, filepath)
        else:
            return jsonify({'error': '不支持的报告格式'}), 400
            
        if success:
            return jsonify({
                'filename': filename,
                'filepath': filepath,
                'download_url': f'/api/report/download/{filename}'
            })
        else:
            return jsonify({'error': '生成报告失败'}), 500
            
    except Exception as e:
        return jsonify({'error': f'生成报告时发生错误: {str(e)}'}), 500


@app.route('/api/report/download/<filename>')
def download_report(filename):
    """下载报告"""
    try:
        # 使用绝对路径
        reports_dir = os.path.abspath('reports')
        filepath = os.path.join(reports_dir, filename)
        
        if not os.path.exists(filepath):
            return jsonify({'error': f'文件不存在: {filepath}'}), 404
            
        return send_file(filepath, as_attachment=True)
        
    except Exception as e:
        return jsonify({'error': f'下载文件时发生错误: {str(e)}'}), 500


def run_scan(scanner: FileScanner, directory: str, scan_id: str):
    """运行扫描（在后台线程中）"""
    global scan_results
    
    # 初始化扫描结果
    scan_results[scan_id] = {
        'status': 'running',
        'progress': {'current': 0, 'total': 0, 'message': '准备扫描...'},
        'start_time': time.time()
    }
    
    # 设置进度回调
    def progress_callback(current, total, message):
        scan_results[scan_id]['progress'] = {
            'current': current,
            'total': total,
            'message': message
        }
    
    scanner.set_progress_callback(progress_callback)
    
    try:
        # 执行扫描
        result = scanner.scan_directory(directory)
        
        # 获取统计信息
        statistics = scanner.get_scan_statistics(result)
        
        # 保存结果
        scan_results[scan_id] = {
            'status': 'completed',
            'data': {
                'total_files': result.total_files,
                'total_size': result.total_size,
                'duplicate_groups': result.duplicate_groups,
                'scan_time': result.scan_time,
                'errors': result.errors
            },
            'statistics': statistics,
            'end_time': time.time()
        }
        
    except Exception as e:
        scan_results[scan_id] = {
            'status': 'error',
            'error': str(e),
            'end_time': time.time()
        }


def create_app():
    """创建 Flask 应用"""
    # 确保必要的目录存在
    os.makedirs('reports', exist_ok=True)
    os.makedirs('backups', exist_ok=True)
    
    return app


if __name__ == '__main__':
    app = create_app()
    app.run(host='0.0.0.0', port=8080, debug=True)