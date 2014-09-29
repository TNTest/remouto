remouto  --imouto hosts 配置工具。
=======

### 文件说明
- 与hosts 同目录的 hosts.bk是修改前的备份文件； 
- hosts.imouto是真正起作用的hosts内容；

### 注意事项
- 运行时windows下防火墙会报修改hosts文件的警告， 请信任操作。
- 运行时需要管理员权限。 

### 已知问题：
  - hosts文件在某些情况下即使管理员权限也无法修改， 程序运行不成功。 解决办法共两种：
    - 修复文件权限为正常状态； 
    - 手工追加hosts.imouto中的内容到hosts文件中
