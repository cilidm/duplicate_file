> 开始： 
> 初次执行 go run main.go -n -p /home/src
>
> 再次执行 go run main.go

> -n 执行新任务 
> -p 检测路径

> 初次执行会先扫描文件并检查重复文件，这个步骤耗时比较长，执行完后 可以在浏览器 http://localhost:8000 查看

> 没有安装golang环境的，可以在命令行执行二进制文件，首次执行要追加-n -p指令，会在目录下生成views网页文件夹以及static js文件夹。

> runtime文件夹保存日志及levelDB数据文件

> 测试环境 darwin10.15.5  linux  win10

> todo 现在保存的key->val 只需要带索引的，其他的删除以节省空间
> 新任务生成时删除全部旧数据