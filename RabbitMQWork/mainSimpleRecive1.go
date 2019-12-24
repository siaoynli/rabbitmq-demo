//Email: 120235331@qq.com
//Github: http：//www.github.com/siaoynli
//Date: 2019/12/24 10:38

package main

import "rabbitmq/RabbitMQ"

func main(){
	rabbitmq :=RabbitMQ.NewRabbitMQSimple("imoocSimple")
	rabbitmq.ConsumeSimple()
}


