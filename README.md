# LanBridge
将多个局域网桥接起来，彼此就像在一个局域网一样互相访问。

Connect multiple LANs to access each other as if they were on one LAN. 


![image](https://github.com/cuirongjie/LanBridge/blob/main/raw_folder/img1.png?raw=true)


# 特点 Features
1. 将多个局域网连接在一起，不同局域网之间可以任意通讯
2. 支持 RDP、SSH、http、https、websocket、ftp、tcp 等任意基于tcp的协议
3. 支持 内网穿透
4. 支持客户端与服务端连接认证，支持客户端与客户端连接认证
5. 安全，可以不暴露服务端口到公网
6. 不依赖任何第三方包
7. 任何系统不会检测出木马病毒

# 用法 Instructions

## 准备
1. 1台有公网IP的主机（任意云主机即可）；
2. 至少两个在不同局域网的主机。

## 用法一、连接多个局域网

#### 实现：在 192.168.1.12 上访问 http://192.168.1.11:9080/ 相当于访问了 http://172.16.16.22:8080/
![image](https://github.com/cuirongjie/LanBridge/blob/main/raw_folder/img1.png?raw=true)

#### 准备：server程序拷贝到服务器(111.204.166.168)上，client程序拷贝到192.168.1.11、172.16.16.21上。

#### 服务端
##### 111.204.166.168中，server.config：
```
{
  "Port": 8000
}
```

#### 客户端
##### 192.168.1.11中，client.config：
```
{
  "ServerAddr": "111.204.166.168 :28010",
  "MyCode": "Client111",
  "Mappings": [
    {
      "LocalPort": 9080,
      "RemoteCode": "Client1621",
      "DistAddr": "172.16.16.22:8080"
    }
  ]
}
```
##### 172.16.16.21中，client.config：
```
{
  "ServerAddr": "111.204.166.168:8000",
  "MyCode": "Client1621"
}
```

运行server和client，在局域网1中任意机器(如192.168.1.12)上访问：http://192.168.1.11:9080/ ，就相当于访问了172.16.16.22上的8080端口服务。
#### 说明：
1. 服务器暴露了连接端口(8000)，没暴露任何其他服务端口，所以服务是安全的；
2. 此方法亦适用于其他协议（RDP、SSH、http、https、websocket、ftp、tcp 等任意基于tcp的协议）。


## 用法二、内网穿透
#### 实现：互联网用户访问 http://111.204.166.168:9080/ 相当于访问了 http://172.16.16.22:8080/
![image](https://github.com/cuirongjie/LanBridge/blob/main/raw_folder/img1.png?raw=true)

#### 准备：server程序拷贝到服务器(111.204.166.168)上，client程序拷贝到172.16.16.21上。

#### 服务端
111.204.166.168中，server.config：
```
{
  "Port": 8000,
  "Mappings": [
    {
      "ServerPort": 8880,
      "RemoteCode": "Client1621",
      "DistAddr": "8080"
    }
}
```

#### 客户端
172.16.16.21中，client.config：
```
{
  "ServerAddr": "111.204.166.168:8000",
  "MyCode": "Client1621"
}
```

运行server和client，任意互联网用户访问：http://111.204.166.168:9080/ ，就相当于访问了172.16.16.22上的8080端口服务。
#### 说明：
1. 服务器暴露了连接端口(8000)和服务端口(9080)；
2. 此方法亦适用于其他协议（RDP、SSH、http、https、websocket、ftp、tcp 等任意基于tcp的协议）。

## 用法三、增加安全性
#### 实现：在“用法一”的基础上，为客户机172.16.16.21增加连接密码

#### 服务端
111.204.166.168中，server.config：
```
{
  "Port": 8000
}
```

#### 客户端
192.168.1.11中，client.config：
```
{
  "ServerAddr": "111.204.166.168 :28010",
  "MyCode": "Client111",
  "Mappings": [
    {
      "LocalPort": 9080,
      "RemoteCode": "Client1621",
      "RemotePassword": "123456",
      "DistAddr": "172.16.16.22:8080"
    }
  ]
}
```

172.16.16.21中，client.config：
```
{
  "ServerAddr": "111.204.166.168:8000",
  "MyCode": "Client1621",
  "MyPassword": "123456"
}
```

#### 说明：
1. 此时172.16.16.21增加了安全性，添加了访问密码123456；192.168.1.11也需要对应配置上访问密码123456；
2. 如果需要为192.168.1.11增加安全性，添加密码，可以采用一样的方式。


#### 实现：在“用法二”的基础上，为服务器111.204.166.168添加连接密码、白名单

#### 服务端
111.204.166.168中，server.config：
```
{
  "Port": 8000,
  "ServerPassword": "qwe123",
  "Whitelist": ["Client1", "Client2", "Client1621"],
  "Mappings": [
    {
      "ServerPort": 8880,
      "RemoteCode": "Client1621",
      "DistAddr": "8080"
    }
}
```

#### 客户端
172.16.16.21中，client.config：
```
{
  "ServerAddr": "111.204.166.168:8000",
  "ServerPassword": "qwe123",
  "MyCode": "Client1621"
}
```

#### 说明：
1. 此时为服务器111.204.166.168添加了连接密码qwe123，每一个客户机也需要对应配置上访问密码qwe123；
2. 此时为服务器111.204.166.168添加白名单，只有识别码(MyCode)在白名单("Client1", "Client2", "Client1621")中的三台客户端可以连接服务器；
3. 如果在“用法一”的基础上为服务器添加连接密码和白名单，可以采用一样的方式。

## 完整配置

#### 服务端
```
{
  "Port": 8000,
  "ServerPassword": "qwe123",
  "Whitelist": ["Client1", "Client2", "Client3"],
  "Mappings": [
    {
      "ServerPort": 8880,
      "RemoteCode": "Client3",
      "DistAddr": "192.168.1.12:8080"
    },
    {
      "ServerPort": 8881,
      "RemoteCode": "Client4",
      "DistAddr": "172.16.16.22:3389"
    }
  ]
}
```

#### 客户端
```
{
  "ServerAddr": "111.204.166.168:8000",
  "ServerPassword": "qwe123",
  "MyCode": "Client3",
  "MyPassword": "123456",
  "Mappings": [
    {
      "LocalPort": 18000,
      "RemoteCode": "Client4",
      "RemotePassword": "111111",
      "DistAddr": "192.168.1.15:3389"
    },
    {
      "LocalPort": 18001,
      "RemoteCode": "Client5",
      "RemotePassword": "222222",
      "DistAddr": "172.16.16.22:8080"
    }
  ]
}
```
