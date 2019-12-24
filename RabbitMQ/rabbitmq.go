//Email: 120235331@qq.com
//Github: http：//www.github.com/siaoynli
//Date: 2019/12/24 9:42
package RabbitMQ

import (
    "fmt"
    "github.com/streadway/amqp"
    "log"
)
//ampq url 格式
const MQURL="amqp://imoocuser:imoocuser@127.0.0.1:5672/imooc"

type RabbitMQ struct {
    conn  *amqp.Connection
    channel  *amqp.Channel
    //队列
    QueueName  string
    //交换机
    Exchange  string
    //binding key
    Key  string
    //连接信息
    Mqurl  string
}

func NewRabbitMQ(queueName, exchange,key string) *RabbitMQ {
    rabbitmq:= &RabbitMQ{
        QueueName: queueName,
        Exchange:  exchange,
        Key:       key,
        Mqurl:     MQURL,
    }
    var err error
    //创建rabbitmq连接
    rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
    rabbitmq.failOnErr(err,"连接错误，创建失败")

    //创建channel
    rabbitmq.channel, err = rabbitmq.conn.Channel()
    rabbitmq.failOnErr(err,"创建channel失败")

    //返回 rabbitmq实例
    return rabbitmq
}


//断开channel 和conn
func (r *RabbitMQ)Destroy(){
    r.channel.Close()
    r.conn.Close()
}

//自定义错误处理函数
func (r *RabbitMQ)failOnErr(err error,message string){
    if(err !=nil) {
        log.Fatalf("%s:%s",message,err)
        panic(fmt.Sprintf("%s:%s",message,err))
    }
}

//创建订阅模式rabbitmq实例
func NewRabbitMQSPubSub(exchangeName string) *RabbitMQ {
    return  NewRabbitMQ("",exchangeName,"")
}


//创建订阅模式rabbitmq实例
func NewRabbitMQRouting(exchangeName string,routingKey string) *RabbitMQ {
    return  NewRabbitMQ("",exchangeName,routingKey)
}


//创建简单模式下rabbitmq实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
    return  NewRabbitMQ(queueName,"","")
}


//simple 模式下生产
func (r *RabbitMQ) PublishSimple(message string){
    //1 申请队列，队列不存在，则创建，存在，则跳过创建
    _, err := r.channel.QueueDeclare(
          r.QueueName,
          false, //是否持久化
          false,  //是否自动删除
          false,   //是否排他性
          false,  //是否阻塞
          nil,   //额外属性
        )
    if err != nil {
        fmt.Println(err.Error())
    }
    //2 发送消息到队列中
    r.channel.Publish(
        r.Exchange,
        r.QueueName,
        //如果未true，根据exchange和routkey规则，如果无法找到符合条件的队列则会把发送的消息返还给发送者
        false,
        //如果未true，当exchange发送消息到队列后发现队列没有绑定消费者，则会把消息返回给发送者
        false,
        amqp.Publishing{
            ContentType:     "text/plain",
            Body:            []byte(message),
        },
    )
}

func (r *RabbitMQ) ConsumeSimple(){
     //1. 申请队列
    _, err := r.channel.QueueDeclare(
        r.QueueName,
        false, //是否持久化
        false,  //是否自动删除
        false,   //是否排他性
        false,  //是否阻塞
        nil,   //额外属性
    )
    if err != nil {
        fmt.Println(err.Error())
    }
    //接受消息
    messages, err := r.channel.Consume(
            r.QueueName,
            //用来区分多个消费者
            "",
            //是否自动应答
            true,
             //是否排他性
            false,
            //如果设置为true，不能将同一个connection发送的消息传递给这个connection中的消费者
            false,
             //是否阻塞   false 表示阻塞
            false,
            nil,
        )

    if err != nil {
        fmt.Println(err.Error())
    }
    forever := make(chan bool)

    //启用协程处理消息
    go func() {
         for  d:= range messages {
              //实现处理的逻辑
              log.Printf("接收到消息:%s",string(d.Body))
         }
    }()

    log.Printf("【*】消费端等待消息... ctrl+c 退出")
    <-forever
}


//订阅 模式下生产
func (r *RabbitMQ) PublishPub(message string){
    //1 尝试创建交换机
     err := r.channel.ExchangeDeclare(
        r.Exchange,
        "fanout",
        true,
        false,
         false,
         false,
         nil,
    )
     r.failOnErr(err,"尝试创建交换机失败:"+r.Exchange)
    //2 发送消息到队列中
    r.channel.Publish(
        r.Exchange,
        r.QueueName,
        //如果未true，根据exchange和routkey规则，如果无法找到符合条件的队列则会把发送的消息返还给发送者
        false,
        //如果未true，当exchange发送消息到队列后发现队列没有绑定消费者，则会把消息返回给发送者
        false,
        amqp.Publishing{
            ContentType:     "text/plain",
            Body:            []byte(message),
        },
    )
}


func (r *RabbitMQ) ConsumeSub(){
    //1 尝试创建交换机
    err := r.channel.ExchangeDeclare(
        r.Exchange,
        "fanout",
        true,
        false,
        false,
        false,
        nil,
    )
    r.failOnErr(err,"尝试创建交换机失败:"+r.Exchange)
    //2创建队列

    queue, err := r.channel.QueueDeclare(
        r.QueueName,  //这里队列必须为空
        false, //是否持久化
        false,  //是否自动删除
        true,   //是否排他性
        false,  //是否阻塞
        nil,   //额外属性
    )
    r.failOnErr(err,"尝试创建随机队列失败")

    //3交换机绑定队列
    err = r.channel.QueueBind(
        queue.Name,
        "", //订阅模式下为空
        r.Exchange,
        false,
        nil,
    )
    r.failOnErr(err,"交换机绑定队列失败")

    //4 接受消息
    messages, err := r.channel.Consume(
        queue.Name,
        //用来区分多个消费者
        "",
        //是否自动应答
        true,
        //是否排他性
        false,
        //如果设置为true，不能将同一个connection发送的消息传递给这个connection中的消费者
        false,
        //是否阻塞   false 表示阻塞
        false,
        nil,
    )

    r.failOnErr(err,"接收消息失败")

    forever := make(chan bool)

    //启用协程处理消息
    go func() {
        for  d:= range messages {
            //实现处理的逻辑
            log.Printf("接收到消息:%s",string(d.Body))
        }
    }()

    log.Printf("【*】消费端等待消息... ctrl+c 退出")
    <-forever
}


//订阅Routing 模式下生产
func (r *RabbitMQ) PublishRouting(message string){
    //1 尝试创建交换机
    err := r.channel.ExchangeDeclare(
        r.Exchange,
        "direct", //路由模式下必须
        true,
        false,
        false,
        false,
        nil,
    )
    r.failOnErr(err,"尝试创建交换机失败:"+r.Exchange)
    //2 发送消息到队列中
    r.channel.Publish(
        r.Exchange,
        r.Key,
        //如果未true，根据exchange和routkey规则，如果无法找到符合条件的队列则会把发送的消息返还给发送者
        false,
        //如果未true，当exchange发送消息到队列后发现队列没有绑定消费者，则会把消息返回给发送者
        false,
        amqp.Publishing{
            ContentType:     "text/plain",
            Body:            []byte(message),
        },
    )
}


func (r *RabbitMQ) ConsumeRouting(){
    //1 尝试创建交换机

    err := r.channel.ExchangeDeclare(
        r.Exchange,
        "direct",
        true,
        false,
        false,
        false,
        nil,
    )
    r.failOnErr(err,"尝试创建交换机失败:"+r.Exchange)
    //2创建队列

    queue, err := r.channel.QueueDeclare(
        r.QueueName,  //这里队列必须为空
        false, //是否持久化
        false,  //是否自动删除
        true,   //是否排他性
        false,  //是否阻塞
        nil,   //额外属性
    )
    r.failOnErr(err,"尝试创建随机队列失败")

    //3交换机绑定队列
    err = r.channel.QueueBind(
        queue.Name,
        r.Key,
        r.Exchange,
        false,
        nil,
    )
    r.failOnErr(err,"交换机绑定队列失败")

    //4 接受消息
    messages, err := r.channel.Consume(
        queue.Name,
        //用来区分多个消费者
        "",
        //是否自动应答
        true,
        //是否排他性
        false,
        //如果设置为true，不能将同一个connection发送的消息传递给这个connection中的消费者
        false,
        //是否阻塞   false 表示阻塞
        false,
        nil,
    )

    r.failOnErr(err,"接收消息失败")

    forever := make(chan bool)

    //启用协程处理消息
    go func() {
        for  d:= range messages {
            //实现处理的逻辑
            log.Printf("接收到消息:%s",string(d.Body))
        }
    }()

    log.Printf("【*】消费端等待消息... ctrl+c 退出")
    <-forever
}



















