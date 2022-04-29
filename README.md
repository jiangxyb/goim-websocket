基于Gin + Vue + Element UI的前后端分离即时通讯系统，只需要配置redis，即可启动
## 特性
1. 登录：http鉴权与websocket鉴权
2. 向单个用户发送消息
3. 向多个用户发送消息
4. 离线消息
5. 查看已注册用户
6. 查看在线用户
7. 单一用户多客户端登录
8. 注销
## 效果视频展示
待补充
## 快速启动
### 前端启动
```shell
cd web
npm intall
npm run serve
```
### 后端启动
在配置文件`./server/config.yaml`填写redis配置
```shell
cd server
go run main.go
```
## 消息结构
采用json传输，结构如下

```json
{
   "type": "",
   "content": "",
   "from": "",
   "to": ["",""]
}
```
## 前端主要功能实现说明
### 界面
vue结合Element UI实现的简单的单路由界面,本人对前端非专业（勿喷）
### 逻辑
1. 注册，登录
    
    点击注册与登录，向后端发送账号和密码
2. jwt鉴权

   登录逻辑：保存token到localStorage，然后使用原生WebSocket与后端建立websocket连接，成功建立连接后发送带有token的消息,对websocket鉴权,消息为
```json
{
   "type": "auth",
   "content": "w-token",
   "from": "",
   "to": []
}
```
3. 普通通信消息
```json
{
   "type": "normal",
   "content": "用户输入的具体消息内容",
   "from": "",
   "to": []
}
```
 4. 发送普通通信消息

   选定接收对象（一个或多个）， 输入消息内容，发送消息，消息结构为
```json
{
   "type": "normal",
   "content": "用户输入的具体消息内容",
   "from": "",
   "to": ["张三","李四","namei","lufei"]
}
```
## 后端主要功能实现说明
### 注册
收到注册请求后，将账号密码保存到文件，同时缓存中也会通过map更新已注册用户。
### 登录
登录成功后颁发token，然后会收到前端发来的websocket连接请求(type为auth)，对token进行验证，建立websocket连接(websocket仅在建立连接时验证token)
### 在线用户列表
用户建立websocket连接后会将连接保存到数据结构`map[用户名]连接句柄list`的结构，用户断开连接会将连接删除，通过遍历map得到在线用户列表，同时将消息通过websocket发送到前端
```json
{
   "type": "loginlist",
   "content": {"names":["张三","李四","wangwu"]},
   "from": "",
   "to": []
}
```
### 消息发送
根据to字段，去在线map找用户，如果用户在线,找到对应连接发送消息，如果离线，赋值`from`字段，然后存储消息到redis(from为key)，当用户登录时，拿到离线消息发给客户端
```json
{
   "type": "normal",
   "content": "用户输入的具体消息内容",
   "from": "钻石王老五",
   "to": ["张三","李四","namei","lufei"]
}
```
### 单一用户多客户端登录
每次用户建立连接都将websocket连接句柄存到map中list结构，遍历list，将消息发给单个用户的所有登录端

