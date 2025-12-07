# SwiftChatGo-朗讯IM即时通讯系统
## 项目简介：
该项目是基于 Golang 生态构建的高可靠分布式实时聊天系统，支持单聊 / 群聊场景下多类型消息的实时传输、离线拉取与持久化存储，通过集群化部署与多级优化，保障高并发下的消息可靠性与系统稳定性。
## 项目功能：

## 项目配置：
docker run -d --name swiftchat -p 8081:8081 --link mysql:mysql --link redis:redis -v /data/swiftchat/upload:/app/asset/upload --restart always swiftchat:latest