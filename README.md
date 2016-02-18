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
        -P string
            passwd to connect monitor db (default "monitor")
        -U string
            user to connect monitor db (default "monitor")
        -h string
            mysql addr (default "127.0.0.1:3306")
        -p string
            passwd to connect db (default "monitor")
        -t uint
            time interval to monitor s (default 10)
        -u string
            user to connect db (default "monitor")

