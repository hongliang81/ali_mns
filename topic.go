package ali_mns

import (
	"fmt"
	"os"
	"strings"
	"net/http"
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

// Decode incoming Notification from Topic mode
func ParseNotification(req *http.Request, msg *TopicNotification) (statusCode int, err error) {

	// 整理Header数据
	var url, contentMd5, contentType, date string
	var mnsSplit = make([]string, 0, 4)

	for k, v := range req.Header {
		switch k1 := strings.ToLower(k); k1 {
		case "content-md5":
			contentMd5 = v[0]
		case "content-type":
			contentType = strings.ToLower(v[0])
		case "date":
			contentType = v[0]
		case "x-mns-request-id":
			mnsSplit = append(mnsSplit, v[0])
		case "x-mns-version":
			mnsSplit = append(mnsSplit, v[0])
		case "x-mns-signing-cert-url":
			mnsSplit = append(mnsSplit, v[0])
			url = v[0]
		}
	}

	fmt.Printf("url[%s]\nmd5[%s]\ntype[%s]\ndate[%s]\n", url, contentMd5, contentType, date)

	//// 获取X509证书
	//certUrl, err := base64.StdEncoding.DecodeString(url)
	//if err != nil {
	//	// TODO
	//	return
	//}
	//resp, err := http.Get(string(certUrl))
	//if err != nil {
	//	// TODO
	//	return
	//}
	//defer resp.Body.Close()
	//
	//block, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(block))	// TODO

	// 计算待签名字符串

	// Authorization解密

	// 认证

	// 返回

	return
}