package ali_mns

import (
	"fmt"
	"os"
	"strings"
)

var (
	//DefaultTopicQPSLimit      int32 = 2000
)

type AliMNSTopic interface {
	Name() string
	SendMessage(message TopicMessageSendRequest) (resp TopicMessageSendResponse, err error)
}

type MNSTopic struct {
	name       string
	client     MNSClient
	stopChan   chan bool
	//qpsLimit   int32
	//qpsMonitor *QPSMonitor
	decoder    MNSDecoder
}

// I prefer subscribe to topic with http endpoint
// So this library is really for subscribing to topic and publishing message
// 'Cause I don't see the need to limit QPS so it's gone for now...
func NewMNSTopic(name string, client MNSClient, /*qps ...int32*/) AliMNSTopic {
	if name == "" {
		panic("ali_mns: topic name could not be empty")
	}

	topic := new(MNSTopic)
	topic.client = client
	topic.name = name
	topic.stopChan = make(chan bool)
	//topic.qpsLimit = DefaultTopicQPSLimit
	topic.decoder = NewAliMNSDecoder()

	//if qps != nil && len(qps) == 1 && qps[0] > 0 {
	//	topic.qpsLimit = qps[0]
	//}

	proxyURL := ""
	topicProxyEnvKey := PROXY_PREFIX + strings.Replace(strings.ToUpper(name), "-", "_", -1)
	if url := os.Getenv(topicProxyEnvKey); url != "" {
		proxyURL = url
	}

	client.SetProxy(proxyURL)

	//topic.qpsMonitor = NewQPSMonitor(5)

	return topic
}

func (p *MNSTopic) Name() string {
	return p.name
}

func (p *MNSTopic) SendMessage(message TopicMessageSendRequest) (resp TopicMessageSendResponse, err error) {
	//p.checkQPS()
	_, err = send(p.client, p.decoder, POST, nil, message, fmt.Sprintf("topics/%s/%s", p.name, "messages"), &resp)
	return
}

