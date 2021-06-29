# 0xDnsRebind

# Usage

简单配置下 [config.yml](./config.yml)
```
domain:
  main: dns-rebind.io # 域名
  ns: [ns1.app.io, ns2.app.io] # 域名NS记录
  ip: 1.1.1.8 # 域名A记录IP
  rebind: 127.0.0.1 # 重新绑定的IP 
```

如何确保顺序访问的？

对每个访问的ip做了缓存，第一次返回 domain.ip 第二次就返回 domain.rebind 

## Demo:

运行 `./0xdnsrebind`

dig 本地测试：

![image](https://user-images.githubusercontent.com/26270009/123769634-145e1380-d8fc-11eb-8ce3-20d694999679.png)
