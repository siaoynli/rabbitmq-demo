//Email: 120235331@qq.com
//Github: http：//www.github.com/siaoynli
//Date: 2019/12/24 10:36

package main

import (
	"fmt"
	"rabbitmq/RabbitMQ"
)

func main(){
	rabbitmq := RabbitMQ.NewRabbitMQSimple("imoocSimple")
	rabbitmq.PublishSimple("hello imooc")
	fmt.Println("发送成功")
}


