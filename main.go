package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const (
	contentLen = 1024
)

var (
	secret             string
	lastSuccessRunTime time.Time
	currentPort        string
)

type ChangePortRequest struct {
	Secret string `form:"secret"`
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().Unix())

	secret = os.Getenv("GIN_SECRET")

	if secret == "" {
		// 随机生成字符串作为secret
		secret = randString(32)
	}

	log.Printf("secret: %v", secret)

	lastSuccessRunTime = time.Now().Add(-time.Hour)

	r := gin.Default()

	r.GET("/current-port", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, fmt.Sprintf("current port is %v", currentPort))
	})

	r.GET("/change-port", func(c *gin.Context) {
		// 每分钟只能运行一次
		if time.Now().Before(lastSuccessRunTime.Add(time.Minute)) {
			c.JSON(http.StatusInternalServerError, "fuck you!!! you can only change the port once every minute!!")
			return
		}

		var req ChangePortRequest
		err := c.ShouldBindWith(&req, binding.Query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		// 判断密码是否正确
		if req.Secret != secret {
			c.JSON(http.StatusInternalServerError, "fuck you!!! the secret is wrong!!")
			return
		}

		// 读取模板文件
		content, err := os.ReadFile("launch.sh.template")
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		contentStr := string(content)

		// 随机生成port
		newPort := strconv.Itoa(10000 + rand.Intn(2000))

		// 写入新的port到模板里
		contentStr = strings.ReplaceAll(contentStr, "{{port}}", newPort)

		// 写入到launch.sh
		file, err := os.OpenFile("launch.sh", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0744)
		_, err = file.WriteString(contentStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		file.Close()

		cmd := exec.Command("docker", "ps", "-q")
		containerIDBytes, err := cmd.CombinedOutput()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		containerID := strings.TrimSpace(string(containerIDBytes))
		log.Printf("old container id: %s", containerID)

		if containerID != "" {
			// kill掉当前docker
			cmd = exec.Command("docker", "kill", containerID)
			output, err := cmd.CombinedOutput()
			log.Printf("docker kill output: %s", output)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		// 重新启动docker
		cmd = exec.Command("bash", "./launch.sh")
		err = cmd.Run()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		log.Printf("change to new port %v, old port: %v", newPort, currentPort)

		lastSuccessRunTime = time.Now()
		currentPort = newPort

		c.JSON(http.StatusOK, gin.H{
			"new_port": newPort,
		})
	})
	r.Run("0.0.0.0:19090")

}
