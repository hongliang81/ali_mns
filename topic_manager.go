package ali_mns

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gogap/errors"
	"strconv"
)

type AliTopicManager interface {
	CreateTopic(location MNSLocation, topicName string, maxMessageSize int32) (err error)
	SetTopicAttributes(location MNSLocation, topicName string, maxMessageSize int32) (err error)
	GetTopicAttributes(location MNSLocation, topicName string) (attr TopicAttribute, err error)
	DeleteTopic(location MNSLocation, topicName string) (err error)
	ListTopic(location MNSLocation, nextMarker string, retNumber int32, prefix string) (topics Topics, err error)
}

type MNSTopicManager struct {
	ownerId         string
	credential      Credential
	accessKeyId     string
	accessKeySecret string

	decoder MNSDecoder
}

func checkTopicName(topicName string) (err error) {
	if len(topicName) > 256 {
		err = ERR_MNS_TOPIC_NAME_IS_TOO_LONG.New()
		return
	}
	return
}

func checkTopicMaxMessageSize(maxSize int32) (err error) {
	if maxSize < 1024 || maxSize > 65536 {
		err = ERR_MNS_MAX_MESSAGE_SIZE_RANGE_ERROR.New()
		return
	}
	return
}

func NewMNSTopicManager(ownerId, accessKeyId, accessKeySecret string) AliTopicManager {
	return &MNSTopicManager{
		ownerId:         ownerId,
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
		decoder:         new(AliMNSDecoder),
	}
}

func checkTopicAttributes(maxMessageSize int32) (err error) {
	if err = checkTopicMaxMessageSize(maxMessageSize); err != nil {
		return
	}
	return
}

func (p *MNSTopicManager) CreateTopic(location MNSLocation, topicName string, maxMessageSize int32) (err error) {
	topicName = strings.TrimSpace(topicName)

	if err = checkTopicName(topicName); err != nil {
		return
	}

	if err = checkTopicAttributes(maxMessageSize); err != nil {
		return
	}

	message := CreateTopicRequest{
		MaxMessageSize:         maxMessageSize,
	}

	url := fmt.Sprintf("http://%s.mns.%s.aliyuncs.com", p.ownerId, string(location))

	cli := NewAliMNSClient(url, p.accessKeyId, p.accessKeySecret)

	var code int
	if code, err = send(cli, p.decoder, PUT, nil, &message, "topics/"+ topicName, nil); err != nil {
		return
	}

	switch code {
	case http.StatusOK:
		return
	case http.StatusNoContent:
		{
			err = ERR_MNS_TOPIC_ALREADY_EXIST_AND_HAVE_SAME_ATTR.New(errors.Params{"name": topicName})
			return
		}
	case http.StatusConflict:
		{
			err = ERR_MNS_TOPIC_ALREADY_EXIST.New(errors.Params{"name": topicName})
			return
		}
	}

	return
}

func (p *MNSTopicManager) SetTopicAttributes(location MNSLocation, topicName string, maxMessageSize int32) (err error) {
	topicName = strings.TrimSpace(topicName)

	if err = checkQueueName(topicName); err != nil {
		return
	}

	if err = checkTopicAttributes(maxMessageSize); err != nil {
		return
	}

	message := CreateTopicRequest{
		MaxMessageSize:         maxMessageSize,
	}

	url := fmt.Sprintf("http://%s.mns.%s.aliyuncs.com", p.ownerId, string(location))

	cli := NewAliMNSClient(url, p.accessKeyId, p.accessKeySecret)

	_, err = send(cli, p.decoder, PUT, nil, &message, fmt.Sprintf("topics/%s?metaoverride=true", topicName), nil)
	return
}

func (p *MNSTopicManager) GetTopicAttributes(location MNSLocation, topicName string) (attr TopicAttribute, err error) {
	topicName = strings.TrimSpace(topicName)

	if err = checkTopicName(topicName); err != nil {
		return
	}

	url := fmt.Sprintf("http://%s.mns.%s.aliyuncs.com", p.ownerId, string(location))

	cli := NewAliMNSClient(url, p.accessKeyId, p.accessKeySecret)

	_, err = send(cli, p.decoder, GET, nil, nil, "topics/"+ topicName, &attr)

	return
}

func (p *MNSTopicManager) DeleteTopic(location MNSLocation, topicName string) (err error) {
	topicName = strings.TrimSpace(topicName)

	if err = checkTopicName(topicName); err != nil {
		return
	}

	url := fmt.Sprintf("http://%s.mns.%s.aliyuncs.com", p.ownerId, string(location))

	cli := NewAliMNSClient(url, p.accessKeyId, p.accessKeySecret)

	_, err = send(cli, p.decoder, DELETE, nil, nil, "topics/"+ topicName, nil)

	return
}

func (p *MNSTopicManager) ListTopic(location MNSLocation, nextMarker string, retNumber int32, prefix string) (topics Topics, err error) {

	url := fmt.Sprintf("http://%s.mns.%s.aliyuncs.com", p.ownerId, string(location))

	cli := NewAliMNSClient(url, p.accessKeyId, p.accessKeySecret)

	header := map[string]string{}

	marker := strings.TrimSpace(nextMarker)
	if len(marker) > 0 {
		if marker != "" {
			header["x-mns-marker"] = marker
		}
	}

	if retNumber > 0 {
		if retNumber >= 1 && retNumber <= 1000 {
			header["x-mns-ret-number"] = strconv.Itoa(int(retNumber))
		} else {
			err = REE_MNS_GET_TOPIC_RET_NUMBER_RANGE_ERROR.New()
			return
		}
	}

	prefix = strings.TrimSpace(prefix)
	if prefix != "" {
		header["x-mns-prefix"] = prefix
	}

	_, err = send(cli, p.decoder, GET, header, nil, "topics", &topics)

	return
}
