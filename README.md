## 一键生成 csdn 博客到 hexo 源文件

> 和[CsdnSyncHexo](https://github.com/flytam/CsdnSyncHexo)一样的功能，只是使用 Go 编写

### 使用方法

#### [下载可执行程序](https://github.com/flytam/blog-sync/releases)

```bash
# 如通过配置文件
./blog-sync-xx --config=./config.json
```


#### 配置项

详细操作可[参考](https://github.com/flytam/blog/issues/14)

> tips: 这边的参数和[CsdnSyncHexo](https://github.com/flytam/CsdnSyncHexo)不太一样

- csdn: csdn 用户名，如 flytam
- output: 输出 hexo markdown 文件路径
- cookie: cookie 信息

#### Build

[参考](./Makefile)
