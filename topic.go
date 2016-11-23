package ali_mns

import (
	"fmt"
	"os"
	"strings"
	"net/http"
	"encoding/base64"
	"io/ioutil"
	"encoding/pem"
	"crypto/x509"
	"sort"
	"bytes"
	"github.com/gogap/errors"
)

var (
	//DefaultTopicQPSLimit      int32 = 2000
	certUrlCached	[]byte
	certCached	*x509.Certificate
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

// Verify Notification Signature
func VerifyNotificationSignature(req *http.Request) error {

	// 整理Header数据
	var authorization, url, contentMd5, contentType, date string
	var mnsSplit = make([]string, 0, 4)

	var count int = 0
	for k, v := range req.Header {
		switch k1 := strings.ToLower(k); k1 {
		case "authorization":
			count++
			authorization = v[0]
		case "content-md5":
			count++
			contentMd5 = v[0]
		case "content-type":
			count++
			contentType = strings.ToLower(v[0])
		case "date":
			count++
			date = v[0]
		case "x-mns-request-id":
			count++
			mnsSplit = append(mnsSplit, k1 + ":" + v[0])
		case "x-mns-version":
			count++
			mnsSplit = append(mnsSplit, k1 + ":" + v[0])
		case "x-mns-signing-cert-url":
			count++
			mnsSplit = append(mnsSplit, k1 + ":" + v[0])
			url = v[0]
		}
	}
	if count < 7 {
		return ERR_MNS_INVALID_NOTIFICATION_HEADER.New(nil)
	}

	// 生成待签名字符串
	sort.Strings(mnsSplit)
	var str2Sign string = req.Method + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n"
	for _, str := range mnsSplit {
		str2Sign += str + "\n"
	}
	str2Sign += req.RequestURI

	// 判断是否需要重新获取X509证书
	certUrl, err := base64.StdEncoding.DecodeString(url)
	if err != nil {
		return ERR_DECODE_URL_FAILED.New(errors.Params{"err": err, "url": url})

	}

	// 获取证书并缓存
	if bytes.Compare(certUrlCached, certUrl) != 0 {
		refreshCert(certUrl)
		certUrlCached = certUrl
	}

	// Authorization解密
	sig2Check, err := base64.StdEncoding.DecodeString(authorization)
	if err != nil {
		return ERR_MNS_SIGNATURE_DOES_NOT_MATCH.New(errors.Params{"err": err})
	}

	// 校验签名
	err = certCached.CheckSignature(x509.SHA1WithRSA, []byte(str2Sign), sig2Check)
	if err != nil {
		return ERR_MNS_SIGNATURE_DOES_NOT_MATCH.New(errors.Params{"err": err})
	}

	return nil
}

// Decode incoming Notification from Topic mode
func ParseNotification(req *http.Request, msg *TopicNotification) (statusCode int, err error) {

	// 校验签名
	err = VerifyNotificationSignature(req)
	if err != nil {
		statusCode = 403
		return
	}

	statusCode = 204

	// 解析消息
	decoder := NewAliMNSDecoder()
	if e := decoder.Decode(req.Body, msg); e != nil {
		err = ERR_UNMARSHAL_NOTIFICATION_FAILED.New(errors.Params{"err": e})
		return
	}

	// TODO 检查MD5

	return
}

// Update certificate
func refreshCert(certUrl []byte) {
	resp, err := http.Get(string(certUrl))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	block, _ := ioutil.ReadAll(resp.Body)

	p, _ := pem.Decode(block)
	certCached, err = x509.ParseCertificate(p.Bytes)
	if err != nil {
		return
	}
}