# flynet 
[![release][1]][4]  [![size][2]][5] [![License][3]][6]
 
[1]: https://img.shields.io/github/v/release/asche910/flynet
[2]: https://img.shields.io/github/repo-size/asche910/flynet
[3]: https://img.shields.io/github/license/asche910/flynet
[4]: https://github.com/asche910/flynet/releases
[5]: https://github.com/asche910/flynet
[6]: https://github.com/asche910/flynet/blob/master/LICENSE



## Features
[flynet](https://github.com/asche910/flynet) Is a command-line tool written in Golang language, currently supported features include：

* [Http proxy](#Http-Proxy)
* [Local Socks5 proxy](#Local-Socks5-Proxy)
* [C/S mode of Socks5 proxy by TCP](#CS-mode-of-Socks5-proxy-by-TCP)
* [C/S mode of Socks5 proxy by UDP](#CS-mode-of-Socks5-proxy-by-UDP)
* [NAT traversal](#NAT-traversal)
* ...

The project is currently divided into the client and the sever. The http and local socks5 proxy, both sides of flynet support. The other functions need to be used at both sides.

## Getting Started
###  Installing
There are two ways to install

##### Download from [Releases Page](https://github.com/asche910/flynet/releases)
Go to the [Releases Page](https://github.com/asche910/flynet/releases), download the corresponding version


##### Build from source

```shell 
go get -u -v github.com/asche910/flynet
```
In windows system，```.\client.exe ...```or  ```.\server.exe ...```

and linux system， ```./server ...``` or ```./client ...```

will show like this ```server ...``` or ```client ...``` in the following article

after run this command, it should output some messages like this:
```
Usage: flynet [options]
  -M, --mode        choose which mode to run. the mode must be one of['http', 'socks5',
                    'socks5-tcp', 'socks5-udp', 'forward']
  -L, --listen      choose which port(s) to listen or forward
  -S, --server      the server address client connect to
  -m, --method      choose a encrypt method, which must be one of ['aes-128-cfb','aes-192-cfb',
                    'aes-256-cfb', 'aes-128-ctr', 'aes-192-ctr', 'aes-256-ctr', 'rc4-md5',
                    'rc4-md5-6', 'chacha20', 'chacha20-ietf'], default is 'aes-256-cfb'
  -P, --password    password of server
  -p, --pac         having this flag, pac-mode will open
  -C, --config      read from config file
  -V, --verbose     output detail info
  -l, --log         output detail info to log file
  -H, --help        show detail usage

Mail bug reports and suggestions to <asche910@gmail.com>
or github: https://github.com/asche910/flynet
```

### Http Proxy

The http proxy directly opens the Http proxy on this machine. Both the client and the server support it. The commands are as follows:

```
server -M http -L 8848 
```
or
```
client -M http -L 8848
```
It means that the Http proxy service is enabled on the port of 8848. If there is no output, it means the startup is successful. After all, one of philosophies of linux is:

> No news is good news 

Of course, if you still want to see some messages, you could run this command with options, like ```-V``` or ```--verbose```, which will output detail message to the terminal.
and the option ```-l``` or ```--log```will write message to the file named 'flynet.log'. of course you can also indicate the file name after the option

### Local Socks5 Proxy

It is very simple to start the socks5 proxy on this machine. Both the client and the server support it. The commands are as follows:

```
server -M socks5  -L 1080
```
or
```
client -M socks5 -L 1080
```
This means that the socks5 proxy is enabled on the port 1080 of the machine, and then Chrome can work well with **Switchy Omega**.

### C/S mode of Socks5 proxy by TCP

The previous one is the socks5 agent on the local, this one is the socks5 proxy that the client and the server cooperate with each other, and the middle is transmitted by the **TCP** protocol. 
By using this mode, you can bypass the **GFW** easily. The usage is as follows:

**Server**
```
server -M socks5-tcp -L 8888
```
**Client**
```
client -M socks5-tcp -L 1080 -S example.com:8888
```


The example here is to assume that my server domain name is example.com ( you can also use ip directly:`-S xxx.xxx.xxx.xxx:8888`), then the client starts the socks5 proxy on port 1080, 
and then the traffic is forwarded to the server's 8888 port by TCP, and the server requests the corresponding target website. Return the result of the request to the client.
The intermediate traffic is encrypted (default method is "aes-256-cfb") to ensure the security of the transmission.
If you want, you can add more options, like
 * ```-P password``` add a password as the key for encryption or decryption
 * ```-m method```   use different encrypt method to encrypt your traffic
 * ```-p```  adding this on client side will start **PAC** mode. but you should setting the url as System auto proxy url manually

If you feel that it is too much trouble to enter a bunch of parameters every time，you can just use ```-C flynet.cnnf``` to load a config file in your current directory

### C/S mode of Socks5 proxy by UDP

This is very similar to the above tcp. The difference is that this is sending and receiving all packets through kcp (UDP) not TCP.
After all, UDP has its own advantages in some aspects, and some important protocols mainly use udp transmission, such as the DNS protocol.
One of uses is that you can have free internet access **without authentication** when using **campus-wifi** or **public-wifi**.
Here's a look at the specific usage:

**Server**
```
server -M socks5-udp -L 53
```
**Client**
```
client -M socks5-udp -L 8848 -S example.com:53
```

Here also take the domain name example.com and port 53 as examples.The client opens the socks5 proxy on udp port 53, and then all traffic
 is transmitted to the server port 53 through the udp mode. After receiving the request, the server then requests all requests.
Send to the target website and return the results to the client in udp mode. The same is that all traffic is encrypted.

### NAT traversal

> Network address translator traversal is a computer networking technique of establishing and maintaining Internet protocol connections across gateways that implement network address translation (NAT).
  
To put it simply, the external network can access the machines in the internal network. What the tool does here is to map a port
 on the intranet to a port on the server, so that by accessing a port on the server, you can indirectly access the port in the intranet.
 Methods as below:

**Server**
```
server -M forward -L 8888 8080
```
**Client**
```
server -M forward -L 80 -S example.com:8888
```

Also assume that the server domain name is example.com. The goal is mapping the port 80 of the client to the server port 8080. the middle of the data transmission is done through the server listening 8888. 
Then we visit example.com:8080 could see the content on the client port 80.

## Conclusion

The current function of the project is relatively limited, and more functions should be added in the future.
If the project is useful to you, a **star** is best favour for me!
