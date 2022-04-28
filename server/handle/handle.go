package handle

import (
	"bufio"
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"im/global"
	"im/middlewares"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Reg struct {
	Name     string `json:"name" binding:"required,min=1,max=20"`
	Password string `json:"password" binding:"required,min=3,max=20"`
}
type UList struct {
	Names []string `json:"names"`
}
type LoginStruct struct {
	Name     string `json:"name" binding:"required,min=1,max=20"`
	Password string `json:"password" binding:"required,min=3,max=20"`
}

func Register(c *gin.Context) {
	var reg Reg
	err := c.Bind(&reg)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "用户名或密码格式错误，请重试",
			"code": "4001",
		})
		return
	}
	StoreUserToMap()
	_, ok := global.UserMap.Load(reg.Name)
	if ok {
		fmt.Println("用户已存在")
		c.JSON(http.StatusOK, gin.H{
			"msg":  "用户已存在，请登录或更换用户名注册",
			"code": "4000",
		})
		return
	}
	if err := AddUserToTxt(reg.Name, reg.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "内部错误",
			"code": "5000",
		})
		return
	}
	StoreUserToMap()
	c.JSON(http.StatusOK, gin.H{
		"msg":  "创建用户成功，请登录",
		"code": "2000",
	})
}

func Login(c *gin.Context) {
	var loginData LoginStruct
	err := c.Bind(&loginData)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "用户名或密码格式错误，请重试",
			"code": "4001",
		})
		return
	}
	psw, ok := global.UserMap.Load(loginData.Name)
	if !ok {
		fmt.Println("用户不存在")
		c.JSON(http.StatusOK, gin.H{
			"msg":  "用户不存在，请注册",
			"code": "4003",
		})
		return
	}
	if loginData.Password != psw.(string) {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "密码错误，请重新输入",
			"code": "4005",
		})
		return
	}

	file, err := os.Open(global.SrvConfig.JWTInfo.PrivateKeyPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	pkBytes, err := ioutil.ReadAll(file)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pkBytes))
	tokenGen := middlewares.NewJWTTokenGen("user", privateKey)
	token, err := tokenGen.GenerateToken(loginData.Name, time.Hour*24*20)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"msg":   "登录成功",
		"code":  "2000",
		"name":  loginData.Name,
		"token": token,
	})
}

func LoginList(c *gin.Context) {
	var users UList
	global.LoginMap.Range(func(key, value interface{}) bool {
		users.Names = append(users.Names, key.(string))
		return true
	})
	c.JSON(http.StatusOK, &users)
}
func getLoginList() *UList {
	var users UList
	global.LoginMap.Range(func(key, value interface{}) bool {
		users.Names = append(users.Names, key.(string))
		return true
	})
	return &users
}

func UserInfo(c *gin.Context) {
	name, _ := c.Get("name")
	userName := name.(string)
	c.JSON(http.StatusOK, gin.H{
		"msg":  "成功",
		"code": "2000",
		"name": userName,
	})
}

func UserList(c *gin.Context) {
	var users UList
	global.UserMap.Range(func(key, value interface{}) bool {
		users.Names = append(users.Names, key.(string))
		return true
	})
	c.JSON(http.StatusOK, &users)
}

var i int

type WsInfo struct {
	Type    string   `json:"type"`
	Content string   `json:"content"`
	To      []string `json:"to"`
	From    string   `json:"from"`
}

func WS(ctx *gin.Context) {
	var claim *middlewares.MyClaim
	wsConn, _ := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	for {
		_, data, err := wsConn.ReadMessage()
		if err != nil {
			wsConn.Close()
			if claim != nil {
				RemoveWSConnFromMap(claim.Subject, wsConn)
				r, _ := json.Marshal(gin.H{
					"type":    "loginlist",
					"content": getLoginList(),
					"to":      []string{},
				})
				SendMsgToAllLoginUser(r)
			}
			fmt.Println(claim.Subject, "出错，断开连接：", err)
			fmt.Println("当前在线用户列表：", getLoginList().Names)
			return
		}

		var wsInfo WsInfo
		json.Unmarshal(data, &wsInfo)
		if wsInfo.Type == "auth" {
			claim, err = middlewares.Auth(wsInfo.Content)
			if err != nil {
				// 认证失败
				fmt.Println(err)
				rsp := WsInfo{
					Type:    "no",
					Content: "认证失败，请重新登录",
					To:      []string{},
				}
				r, _ := json.Marshal(rsp)
				wsConn.WriteMessage(websocket.TextMessage, r)
				wsConn.Close()

				continue
			}
			// 认证成功
			// 将连接加入map记录
			AddWSConnToMap(claim.Subject, wsConn)

			fmt.Println(claim.Subject, " 加入连接")
			fmt.Println("当前在线用户列表：", getLoginList().Names)

			rsp := WsInfo{
				Type:    "ok",
				Content: "连接成功，请发送消息",
				To:      []string{},
			}
			r, _ := json.Marshal(rsp)
			// 更新登录列表
			wsConn.WriteMessage(websocket.TextMessage, r)
			r, _ = json.Marshal(gin.H{
				"type":    "loginlist",
				"content": getLoginList(),
				"to":      []string{},
			})
			SendMsgToAllLoginUser(r)
			// 发送离线消息
			cmd := global.Redis.LRange(context.Background(), claim.Subject, 0, -1)
			msgs, err := cmd.Result()
			if err != nil {
				log.Println(err)
				continue
			}
			for _, msg := range msgs {
				wsConn.WriteMessage(websocket.TextMessage, []byte(msg))
			}
			global.Redis.Del(context.Background(), claim.Subject)

		} else {
			rsp, _ := json.Marshal(gin.H{
				"type":    "normal",
				"content": wsInfo.Content,
				"to":      []string{},
				"from":    claim.Subject,
			})
			SendMsgToOtherUser(rsp, claim.Subject, wsInfo.To...)
		}
	}
	wsConn.Close()
}

var (
	Upgrader = websocket.Upgrader{
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func AddWSConnToMap(userName string, wsConn *websocket.Conn) {
	// 同一用户可以有多个ws连接（登录多次）
	loginListInter, ok := global.LoginMap.Load(userName)
	if !ok {
		// 之前没登录
		loginList := list.New()
		loginList.PushBack(wsConn)
		global.LoginMap.Store(userName, loginList)
	} else {
		// 多次登录
		loginList := loginListInter.(*list.List)
		loginList.PushBack(wsConn)
		global.LoginMap.Store(userName, loginList)
	}
}

func RemoveWSConnFromMap(userName string, wsConn *websocket.Conn) {
	loginListInter, ok := global.LoginMap.Load(userName)
	if !ok {
		fmt.Println("没有连接可以关闭")
	} else {
		// 有连接
		loginList := loginListInter.(*list.List)
		if loginList.Len() <= 1 {
			global.LoginMap.Delete(userName)
		} else {
			for e := loginList.Front(); e != nil; e = e.Next() {
				if e.Value.(*websocket.Conn) == wsConn {
					loginList.Remove(e)
					break

				}
			}
			global.LoginMap.Store(userName, loginList)
		}
	}
}

func SendMsgToOtherUser(data []byte, myName string, otherUserName ...string) {
	for _, otherName := range otherUserName {
		if otherName != myName {
			v, ok := global.LoginMap.Load(otherName)
			if ok {
				// 在线，发送给目标用户的所有客户端
				l := v.(*list.List)
				for e := l.Front(); e != nil; e = e.Next() {
					conn := e.Value.(*websocket.Conn)
					conn.WriteMessage(websocket.TextMessage, data)
				}
			} else {
				_, ok := global.UserMap.Load(otherName)
				if ok {
					//离线消息缓存到redis
					global.Redis.LPush(context.Background(), otherName, data)
				}
			}
		}

		//global.LoginMap.Range(func(key, value interface{}) bool {
		//	if key.(string) != myName && key.(string) == otherName {
		//		l := value.(*list.List)
		//		for e := l.Front(); e != nil; e = e.Next() {
		//			conn := e.Value.(*websocket.Conn)
		//			conn.WriteMessage(websocket.TextMessage, data)
		//		}
		//	}
		//	return true
		//})
	}
}
func SendMsgToAllLoginUser(data []byte) {
	global.LoginMap.Range(func(key, value interface{}) bool {
		l := value.(*list.List)
		for e := l.Front(); e != nil; e = e.Next() {
			conn := e.Value.(*websocket.Conn)
			conn.WriteMessage(websocket.TextMessage, data)
		}
		return true
	})
}

func AddUserToTxt(userName, psw string) error {
	f, err := os.OpenFile("./user.txt", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = f.WriteString(fmt.Sprintf("\nu:%v\np:%v", userName, psw))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func StoreUserToMap() {
	f, err := os.Open("./user.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	r := bufio.NewReader(f)
	for {
		str, err := r.ReadString('\n')
		if err == io.EOF {
			fmt.Println(err)
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		flagAndData := strings.Split(str, ":")
		if len(flagAndData) != 2 {
			continue
		}
		userName := strings.TrimRight(flagAndData[1], "\n")
		userName = strings.TrimRight(userName, "\r")
		if flagAndData[0] == "u" {
			str, err := r.ReadString('\n')

			flagAndData := strings.Split(str, ":")
			if len(flagAndData) != 2 {
				continue
			}
			var psw = strings.TrimRight(flagAndData[1], "\n")
			psw = strings.TrimRight(psw, "\r")
			if flagAndData[0] == "p" {
				global.UserMap.Store(userName, psw)
				if err == io.EOF {
					fmt.Println(err)
					break
				}
				if err != nil {
					fmt.Println(err)
					break
				}
			} else {
				continue
			}
		} else {
			continue
		}

	}
}
