//Email: 120235331@qq.com
//Github: http：//www.github.com/siaoynli
//Date: 2019/12/24 10:36

package main

import (
	"fmt"
	"rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main(){
	rabbitmq := RabbitMQ.NewRabbitMQRouting("imoocExchangeRoute","imooc")
	rabbitmq2 := RabbitMQ.NewRabbitMQRouting("imoocExchangeRoute","imooc2")

	for num:=0;num<=100 ; num++  {
		rabbitmq.PublishRouting("hello imooc:"+strconv.Itoa(num))
		rabbitmq2.PublishRouting("hello imooc2:"+strconv.Itoa(num))
		time.Sleep(time.Second)
		fmt.Println("发送成功",num)
	}
}


