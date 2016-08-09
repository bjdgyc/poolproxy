#poolproxy

## Introduction
Pool Proxy
连接池代理是一个使用golang编写的简单连接池代理工具

该工具提供了连接池功能，并可设置最大连接数、连接最大空闲时间、定时检测并断开空闲连接

该工具提供了一个透明的代理接口，可以为下游程序提供带连接池的代理功能

## TODO
目前仅实现了redis的代理功能

本地使用时，建议监听`Unix domain socket`
可以有效减少TCP握手消耗，提高系统性能


## Installation

`go get `

## Toml config

``` toml

[options.redis]
    #日志文件
    logfile = ""
    # 代理监听设置
    # 可以设置为Unix socket
    # 如: /var/run/poolproxy.socket
    # 也可以设置为TCP端口
    addr = ":8080"
    read_timeout = 0
    write_timeout = 0
    pool_timeout = 0

    #远程连接设置
    raddr = "192.168.56.102:6379"
    ruser = ""
    rpass = ""
    rpool_size = 0
    ridle_timeout = 0
    ridle_check_frequency = 3
    
```

## Start

`go build && ./poolproxy -c ./config.toml`
