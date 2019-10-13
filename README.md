# flynet 
## 前言
> 前段时间做某个项目，由于涉及到tcp/udp方面的知识比较多，于是就索性趁热打铁，写个工具来强化相关知识。另外由于并非十分擅长Golang，所以也顺便再了解下Golang吧。

## 简介
[flynet](https://github.com/asche910/flynet) 是一款Golang语言编写的命令行工具，目前支持的功能包括：

* Http代理
* 本地Socks5代理
* C/S模式的Socks5代理，支持TCP/UDP方式
* 内网穿透
* ...
项目目前分为clien端和sever端，除http、本地socks5代理两端都支持外，其余功能需要两端配合使用。

## 使用方式
###  安装
Windows、linux用户可以直接在[Releases页面](https://github.com/asche910/flynet/releases)下载对应的版本即可，其他平台可自行下载源码编译。

Windows中命令行进入到相应目录，```.\win-client.exe ...```或  ```.\win-server.exe ...```

Linux中同样的， ```./linux-server ...```或```./linux-client ...```

在下文中皆以```server ...```或```client ...```表示。

尝试运行后，如果输出如下信息表示成功：
```
Usage: flynet [options]
  -M, --mode        choose which mode to run. the mode must be one of['http', 'socks5',
                    'socks5-tcp', 'socks5-udp', 'forward']
  -L, --listen      choose which port(s) to listen or forward
  -S, --server      the server address client connect to
  -V, --verbose     output detail info
  -l, --log         output detail info to log file
  -H, --help        show detail usage

Mail bug reports and suggestions to <asche910@gmail.com>
or github: https://github.com/asche910/flynet
```

### Http代理
http代理直接在本机上开启Http代理，client和server都支持，命令如下：
```
server -M http -L 8848 
```
或
```
client -M http -L 8848
```
表示在本机8848端口上开启了Http代理服务，如果没有任何信息输出则表示启动成功，毕竟linux的一大哲学就是：
> 没有消息就是好消息

当然如果还是想看到消息的话，可以在后面加上 ```-V```或```--verbose```参数，这样的话就会输出很多消息了。或者也可以加上```-l```或```--log```参数来启动日志文件，会在运行目录下生成一个 ```flynet.log```文件。

### 本地Socks5代理
本机上开启socks5代理的话，也是非常简单的，client和server都支持，命令如下：
```
server -M socks5  -L 8848
```
或
```
client -M socks5 -L 8848
```
这就表示在本机8848端口上开启了socks5代理，然后Chrome配合SwitchyOmega就可以很好的上网了。

### C/S模式的Socks5代理-TCP
前面的那个是在本地上的socks5代理，这个则是client和server相互配合的socks5代理，并且中间是以tcp协议传输。用途的话，自由发挥吧。使用方法如下：

**服务端**
```
server -M socks5-tcp -L 8888
```
**客户端**
```
client -M socks5-tcp -L 8848 -S asche.top:8888
```
这里的例子是假设我服务器域名为 asche.top，然后客户端在8848端口开启了socks5代理，然后流量是以TCP的方式转发到了服务器的8888端口上，交由服务器去请求相应的目标网站，再把请求结果返回给客户端。如果可以，中间流量再进行加密，保证了传输的安全性。


### C/S模式的Socks5代理-UDP
这个和上面tcp那个非常相似，不同的是这个使用UDP报文进行传输。毕竟UDP在某些方面有它自身的优势，而且某些重要的协议主要使用udp传输，比如DNS协议。下面来介绍具体用法：

**服务端**
```
server -M socks5-udp -L 53
```
**客户端**
```
client -M socks5-udp -L 8848 -S asche.top:53
```
这里同样以域名asche.top、端口53为例，客户端在8848端口开启了socks5代理，然后所有流量通过udp方式传输到服务端的53端口上，服务端收到后解析请求，然后将所有请求发至目标网站，再将结果以udp方式返回到客户端。同样的是中间传输也进行了加密。


### 内网穿透
> 内网穿透，即NAT穿透，网络连接时术语，计算机是局域网内时，外网与内网的计算机节点需要连接通信，有时就会出现不支持内网穿透。就是说映射端口，能让外网的电脑找到处于内网的电脑，提高下载速度

简单点说就是让外网能够访问到内网中的机器。这里该工具所做的就是将内网的某个端口映射到服务器的某个端口中去，这样通过访问服务器的某个端口就可以间接的访问到内网中的端口了。方法如下：

**服务端**
```
server -M forward -L 8888 8080
```
**客户端**
```
server -M forward -L 80 -S asche.top:8888
```
 同样假设服务器域名为asche.top, 这样所完成的就是将客户端的80端口映射到了服务端的8080端口上，中间的数据传输是通过服务端监听8888来完成的。然后我们访问asche.top:8080看到的内容应该就是客户端80端口上的内容了。

## 结语
项目目前功能也比较局限，日后应该会加上更多功能。另外地址位于 [flynet](https://github.com/asche910/flynet), 还望大家多多支持！
