package backend

import (
	"os"
	"strconv"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"net/http"

	"github.com/liuhan907/waka/waka-cow/database"
	"github.com/liuhan907/waka/waka-cow/modules/hall/hall_message"
	"github.com/liuhan907/waka/waka-cow/proto"
)

var (
	log = logrus.WithFields(logrus.Fields{
		"pid":    os.Getpid(),
		"module": "cow.backend",
	})
)

// 消息转发目标创建器
type TargetCreator func() *actor.PID

// 配置
type Option struct {
	// 创建器
	TargetCreator TargetCreator

	// 监听地址
	Address string
}

func Start(option Option) {
	target := option.TargetCreator()

	router := gin.Default()
	router.GET("/player/changed/:id", func(c *gin.Context) {
		param := c.Param("id")
		id, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			c.Status(400)
			return
		}
		database.RefreshCache(database.Player(id))
		log.WithFields(logrus.Fields{
			"id": id,
		}).Debug("player changed")

		c.Status(200)
	})
	router.GET("/configuration/changed/", func(c *gin.Context) {
		database.RefreshConfiguration()

		log.Debug("configuration changed")

		c.Status(200)
	})
	router.GET("/room/flowing/query", func(c *gin.Context) {
		ch := make(chan interface{})
		defer close(ch)
		target.Tell(&hall_message.GetFlowingRoom{
			Respond: func(response []*cow_proto.NiuniuRoomData, e error) {
				if e != nil {
					ch <- e
				} else {
					ch <- response
				}
			},
		})
		response := <-ch
		switch evd := response.(type) {
		case []*cow_proto.NiuniuRoomData:
			c.JSON(http.StatusOK, gin.H{"rooms": evd})
			c.Status(200)
		default:
			c.BindJSON(
				struct {
					Err interface{}
				}{
					Err: evd,
				})
		}
	})
	router.GET("/room/player/query", func(c *gin.Context) {
		ch := make(chan interface{})
		defer close(ch)

		target.Tell(&hall_message.GetPlayerRoom{
			Respond: func(response []*cow_proto.NiuniuRoomData, e error) {
				if e != nil {
					ch <- e
				} else {
					ch <- response
				}
			},
		})
		response := <-ch
		switch evd := response.(type) {
		case []*cow_proto.NiuniuRoomData:
			log.WithFields(logrus.Fields{
				"response": response,
			}).Warnln("room/player/query")
			c.JSON(http.StatusOK, gin.H{"rooms": evd})
			c.Status(200)
		default:
			c.BindJSON(
				struct {
					Err interface{}
				}{
					Err: evd,
				})
		}
	})
	router.GET("/player/online", func(c *gin.Context) {
		ch := make(chan interface{})
		defer close(ch)
		target.Tell(&hall_message.GetOnlinePlayer{
			Respond: func(response []int32, e error) {
				if e != nil {
					ch <- e
				} else {
					ch <- response
				}
			},
		})
		response := <-ch
		switch evd := response.(type) {
		case []int32:
			c.BindJSON(evd)
			c.JSON(http.StatusOK, gin.H{"players": evd})
			c.Status(200)
		default:
			c.BindJSON(
				struct {
					Err interface{}
				}{
					Err: evd,
				})
		}
	})

	go func() {
		err := router.Run(option.Address)
		if err != nil {
			log.WithFields(logrus.Fields{
				"address": option.Address,
				"err":     err,
			}).Fatalln("listen failed")
		}
	}()

	log.WithFields(logrus.Fields{
		"address": option.Address,
	}).Infoln("listen started")
}
