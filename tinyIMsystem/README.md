# Go-IM-System

## English

An IM-System program based on Aceld's Golang lesson. Using message 
read-write separation.

The system includes structs Server, User and Client. Server and 
User are both server-side components. Server encapsulates the server
source related to connections, while User represents the users 
connected to the server. They work together to handle the communication
between the client and the server. The last one is the client-side 
component.

Server mainly offers methods for initialization, start and handling 
of service. Run method init and start service to listen to the server
port and deal with requests using the Handle method after accepting 
them. Create a goroutine to broadcast messages asynchronously.

## 简体中文

基于刘丹冰Go课程的即时通信系统，使用了读写分离模型。

系统主要结构体包括：Server、User和Client。 前二者同属于服务端，User是对一个连接及
其相关Server资源的包装；后者是客户端。

Server主要方法包括`Run()`启动服务监听端口，接收请求并交由`Handler()`处理，同时开
启goroutine异步广播信息。其他的功能主要由User实现。

User中封装了用户上线、下线、超时、监听方法。同时根据Client发送的信息，选择执行公聊、
私聊、 查询在线用户和重命名的功能。
