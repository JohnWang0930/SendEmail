package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"strings"
)

type config_map struct {
	Host     string `json:host`
	Port     string `json:port`
	Username string `json:username`
	Password string `json:password`
	From     string `json:from`
	Subject  string `json:subject`
	Nickname string `json:nickname`
	Dbconfig string `json:dbconfig`
}

func initConfig() config_map {
	// 读取文件
	configFile, err := ioutil.ReadFile("./emailConfig.json")

	if err != nil {
		log.Fatal(err)
	}

	// 转换配置
	config := config_map{}
	if err := json.Unmarshal(configFile, &config); err != nil {
		log.Fatal(err)
	}
	return config
}

func main() {
	config := initConfig()

	to := GetEmailListBySql(config.Dbconfig)
	body, err := ioutil.ReadFile("./content.html")
	msg := []byte(
		"To: " + strings.Join(to, ",") + "\r\n" +
			"From: " + config.Nickname + "<" + config.From + ">\r\n" +
			"Subject: " + config.Subject + "\r\n" +
			"Content-Type: text/html; charset=UTF-8" + "\r\n" +
			"\r\n" +
			string(body))

	// 连接smtp服务
	c, err := smtp.Dial(config.Host + ":" + config.Port)
	if err != nil {
		log.Fatal(err)
	}

	// 获取auth
	auth := LoginAuth(config.Username, config.Password)

	// 认证
	err = c.Auth(auth)
	if err != nil {
		log.Fatal(err)
	}

	// 设置发送者和接收者
	if err := c.Mail(config.From); err != nil {
		log.Fatal(err)
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			log.Fatal(err)
		}
	}

	// 发送body
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(wc, "%v", string(msg))
	if err != nil {
		log.Fatal(err)
	}

	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}

}
