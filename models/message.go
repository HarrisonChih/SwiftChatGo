package models

import (
	"context"
	"fmt"
	"ginchat/utils"
	"github.com/fatih/set"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// 聊天消息模型
type Message struct {
	gorm.Model
	UserId     int64  //发送者
	TargetId   int64  //接受者
	Type       int    //发送类型  1私聊  2群聊  3心跳
	Media      int    //消息类型  1文字 2表情包 3语音 4图片 /表情包
	Content    string //消息内容
	CreateTime uint64 //创建时间
	ReadTime   uint64 //读取时间
	Pic        string
	Url        string
	Desc       string
	Amount     int //其他数字统计
}

func (table *Message) TableName() string {
	return "message"

}

// WebSocket 连接的核心封装模型 每个在线用户对应一个node
type Node struct {
	Conn          *websocket.Conn //连接
	Addr          string          //客户端地址
	FirstTime     uint64          //首次连接时间
	HeartbeatTime uint64          //心跳时间
	LoginTime     uint64          //登录时间
	DataQueue     chan []byte     //消息
	GroupSets     set.Interface   //好友 / 群
}

// 映射关系 维护在线用户
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwlocker sync.RWMutex

func Chat(response http.ResponseWriter, request *http.Request) {
	//1.获取参数，并检验token等合法性
	query := request.URL.Query()
	userId, _ := strconv.ParseInt(query.Get("userId"), 10, 64)
	//msgType := query.Get("type")
	//targetId := query.Get("targetId")
	//context := query.Get("context")
	isvalida := true //todo checkToken()
	//升级HTTP连接为WebSocket（允许合法请求跨域
	conn, err := (&websocket.Upgrader{
		//token 校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(response, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//2.获取conn
	currentTime := uint64(time.Now().Unix())
	node := &Node{
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(), //客户端地址
		HeartbeatTime: currentTime,                //心跳时间
		LoginTime:     currentTime,                //登录时间
		DataQueue:     make(chan []byte, 50),
		GroupSets:     set.New(set.ThreadSafe),
	}

	//3.用户关系

	//4.userid 与 node 绑定并加锁
	rwlocker.Lock()
	clientMap[userId] = node
	rwlocker.Unlock()

	//5.启动消息发送协程（异步消费消息队列）
	go sendProc(node)

	//6. 启动消息接收协程（监听客户端发送的消息）
	go recvProc(node)
	//sendMsg(int64(userId), []byte("欢迎进入聊天系统"))
}

// sendProc：异步发送消息（消费DataQueue中的消息）
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue: //// 从消息队列取消息
			fmt.Println("[ws]sendProc >>>> msg :", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// recvProc：监听客户端消息（生产消息到UDP通道）
func recvProc(node *Node) {
	for {
		_, message, err := node.Conn.ReadMessage() // 读取客户端消息
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := Message{}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Println(err)
		}
		broadMsg(message) // 将消息写入UDP发送通道，转发给其他服务节点
		if msg.Type != 3 {
			dispatch(message)
			broadMsg(message) // 将消息写入UDP发送通道，转发给其他服务节点
			fmt.Println("[ws] <<<< ", string(message))
		} else {
			currentTime := uint64(time.Now().Unix())
			node.Heartbeat(currentTime)
		}

	}
}

var udpSendChan = make(chan []byte) // UDP消息发送通道

func broadMsg(message []byte) {
	udpSendChan <- message
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

// udpSendProc：UDP发送协程（向指定UDP地址广播消息）
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(58, 198, 176, 23),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		select {
		case data := <-udpSendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// udpRecvProc：UDP接收协程（监听UDP端口，接收其他服务节点的消息）
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: viper.GetInt("port.udp"),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	for {
		var buf [1024]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

// dispatch：根据消息类型分发消息
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1:
		sendMsg(msg.TargetId, data)
		fmt.Println("dispatch  data :", string(data))
	case 2:
		sendGroupMsg(msg.TargetId, data) //发送的群ID ，消息内容
	//case 3:
	//	sendAllMsg()
	default:
		return
	}

}

func sendMsg(userId int64, msg []byte) {

	rwlocker.RLock()
	node, ok := clientMap[userId]
	rwlocker.RUnlock()
	jsonMsg := Message{}
	err := json.Unmarshal(msg, &jsonMsg)
	ctx := context.Background()
	targetIdStr := strconv.Itoa(int(userId))
	userIdStr := strconv.Itoa(int(jsonMsg.UserId))
	jsonMsg.CreateTime = uint64(time.Now().Unix())
	//r, err := utils.Red.Get(ctx, "online_"+userIdStr).Result()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//if r != "" {
	//	if ok {
	//		fmt.Println("sendMsg >>> userID: ", userId, "  msg:", string(msg))
	//		node.DataQueue <- msg
	//	}
	//}
	if ok {
		fmt.Println("sendMsg >>> userID: ", userId, "  msg:", string(msg))
		node.DataQueue <- msg
	}
	var key string
	if userId > jsonMsg.UserId {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}
	res, err := utils.Red.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
	}
	score := float64(len(res)) + 1
	ress, e := utils.Red.ZAdd(ctx, key, redis.Z{score, msg}).Result() //jsonMsg
	//res, e := utils.Red.Do(ctx, "zadd", key, 1, jsonMsg).Result() //备用 后续拓展 记录完整msg
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(ress)
}

func sendGroupMsg(targetId int64, msg []byte) {
	fmt.Println("开始群发消息")
	userIds := SearchUserByGroupId(uint(targetId))
	for i := 0; i < len(userIds); i++ {
		//排除给自己的
		if targetId != int64(userIds[i]) {
			sendMsg(int64(userIds[i]), msg)
		}

	}
}

// 需要重写此方法才能完整的msg转byte[]
func (msg Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(msg)
}

// 获取缓存里面的消息
func RedisMsg(userIdA int64, userIdB int64, start int64, end int64, isRev bool) []string {
	//rwlocker.RLock()
	//node, ok := clientMap[userIdA]
	//rwlocker.RUnlock()
	//jsonMsg := Message{}
	//json.Unmarshal(msg, &jsonMsg)
	ctx := context.Background()
	userIdStr := strconv.Itoa(int(userIdA))
	targetIdStr := strconv.Itoa(int(userIdB))
	var key string
	if userIdA > userIdB {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}

	//rels, err := utils.Red.ZRevRange(ctx, key, 0, 10).Result() //根据score倒叙

	var rels []string
	var err error
	if isRev {
		rels, err = utils.Red.ZRange(ctx, key, start, end).Result()
	} else {
		rels, err = utils.Red.ZRevRange(ctx, key, start, end).Result()
	}
	if err != nil {
		fmt.Println(err) //没有找到
	}
	// 发送推送消息
	/**
	// 后台通过websoket 推送消息
	for _, val := range rels {
		fmt.Println("sendMsg >>> userID: ", userIdA, "  msg:", val)
		node.DataQueue <- []byte(val)
	}**/
	return rels
}

// 更新用户心跳
func (node *Node) Heartbeat(currentTime uint64) {
	node.HeartbeatTime = currentTime
	return
}

// 清理超时连接
func CleanConnection(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("cleanConnection err", r)
		}
	}()
	//fmt.Println("定时任务,清理超时连接 ", param)
	//node.IsHeartbeatTimeOut()
	currentTime := uint64(time.Now().Unix())
	for i := range clientMap {
		node := clientMap[i]
		if node.IsHeartbeatTimeOut(currentTime) {
			fmt.Println("心跳超时..... 关闭连接：", node)
			node.Conn.Close()
		}
	}
	return result
}

// 用户心跳是否超时
func (node *Node) IsHeartbeatTimeOut(currentTime uint64) (timeout bool) {
	if node.HeartbeatTime+viper.GetUint64("timeout.HeartbeatMaxTime") <= currentTime {
		fmt.Println("心跳超时。。。自动下线", node)
		timeout = true
	}
	return
}
