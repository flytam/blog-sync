## 一键生成csdn博客到hexo源文件

> 和[CsdnSyncHexo](https://github.com/flytam/CsdnSyncHexo)一样的功能，只是使用Go编写


### 使用方法

#### 下载可执行程序

- [Mac](bin/blog-sync-mac)
- [Linux](bin/blog-sync-linux)
- [Windows](bin/blog-sync-win)

```bash
# 如通过配置文件
./blog-sync-xx --config=./config.json
```

#### 通过Go环境

待更新..

#### 配置项

详细操作可[参考](https://github.com/flytam/blog/issues/14)

- csdn: csdn用户名，如flytam
- output: 输出hexo markdown文件路径
- cookie: cookie信息


#### Build

[参考](./Makefile)

