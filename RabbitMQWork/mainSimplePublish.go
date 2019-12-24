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
	rabbitmq := RabbitMQ.NewRabbitMQSimple("imoocSimple")

	for num:=0;num<=100 ; num++  {
		rabbitmq.PublishSimple("hello imooc:"+strconv.Itoa(num))
		time.Sleep(time.Second)
		fmt.Println("发送成功",num)
	}
}


