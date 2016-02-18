dbmonitor

--------------
## 功能介绍
- 轻量级MySQL监控系统
- golang1.5

## 基本原理
  集中收集信息式的监控，定时的执行 show global status where variable_name in (xxx),获取数据库状态信息，并存储进MySQL
  与阈值比较，触发对应报警动作

## 使用说明
- git clone https://github.com/tangbinbin/dbmonitor.git
- make

### 查看使用帮助
    ./bin/dbmonitor -h
    Usage of ./bin/dbmonitor:
        -p string
            连接存储信息的MySQL的密码 (default "monitor")
        -u string
            连接存储信息的MySQL的用户名 (default "monitor")
        -h string
            存储信息的MySQL的地址 (default "127.0.0.1:3306")
        -P string
            连接监控MySQL的密码 (default "monitor")
        -t uint
            监控时间间隔 (default 10)
        -U string
            连接监控MySQL的用户 (default "monitor")

