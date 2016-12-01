package ali_mns

import (
	"encoding/xml"
)

type TopicMessageResponse struct {
	XMLName   xml.Name `xml:"Message" json:"-"`
	Code      string   `xml:"Code,omitempty" json:"code,omitempty"`
	Message   string   `xml:"Message,omitempty" json:"message,omitempty"`
	RequestId string   `xml:"RequestId,omitempty" json:"request_id,omitempty"`
	HostId    string   `xml:"HostId,omitempty" json:"host_id,omitempty"`
}

type TopicMessageSendRequest struct {
	XMLName      xml.Name    `xml:"Message" json:"message"`
	MessageBody  Base64Bytes `xml:"MessageBody" json:"message_body"`
	MessageTag   string      `xml:"MessageTag,omitempty" json:"message_tag,omitempty"`
}

type TopicMessageSendResponse struct {
	TopicMessageResponse
	MessageId      string `xml:"MessageId" json:"message_id"`
	MessageBodyMD5 string `xml:"MessageBodyMD5" json:"message_body_md5"`
}

type TopicSubscribeRequest struct {
	XMLName		xml.Name	`xml:"Subscription" json:"-"`
	Endpoint	string		`xml:"Endpoint" json:"endpoint"`
	FilterTag	string		`xml:"FilterTag,omitempty" json:"filter_tag,omitempty"`
	NotifyStrategy	string		`xml:"NotifyStrategy,omitempty" json:"notify_strategy,omitempty"`
	NotifyContentFormat	string	`xml:"NotifyContentFormat,omitempty" json:"notify_content_format,omitempty"`
}

type TopicSubscription struct {
	XMLName		xml.Name	`xml:"Subscription" json:"-"`
	Subscriber	string		`xml:"Subscriber" json:"subscriber"`
	TopicOwner	string		`xml:"TopicOwner" json:"topic_owner"`
	TopicName	string		`xml:"TopicName" json:"topic_name"`
	Endpoint	string		`xml:"Endpoint" json:"endpoint"`
	NotifyStrategy	string		`xml:"NotifyStrategy,omitempty" json:"notify_strategy,omitempty"`
	NotifyContentFormat	string	`xml:"NotifyContentFormat,omitempty" json:"notify_content_format,omitempty"`
	FilterTag	string		`xml:"FilterTag" json:"filter_tag"`
	CreateTime	int64		`xml:"CreateTime" json:"create_time"`
	LastModifyTime	int64		`xml:"LastModifyTime" json:"last_modify_time"`
}

type CreateTopicRequest struct {
	XMLName                xml.Name `xml:"Topic" json:"-"`
	MaxMessageSize         int32    `xml:"MaximumMessageSize,omitempty" json:"maximum_message_size,omitempty"`
}

type TopicAttribute struct {
	XMLName                xml.Name `xml:"Topic" json:"-"`
	TopicName              string   `xml:"TopicName,omitempty" json:"topic_name,omitempty"`
	CreateTime             int64    `xml:"CreateTime,omitempty" json:"create_time,omitempty"`
	LastModifyTime         int64    `xml:"LastModifyTime,omitempty" json:"last_modify_time,omitempty"`
	MaxMessageSize         int32    `xml:"MaximumMessageSize,omitempty" json:"maximum_message_size,omitempty"`
	MessageRetentionPeriod int32    `xml:"MessageRetentionPeriod,omitempty" json:"message_retention_period,omitempty"`
	MessageCount		int64	`xml:"MessageCount,omitempty" json:"message_count,omitempty"`
}

type Topic struct {
	TopicURL string `xml:"TopicURL" json:"url"`
}

type Topics struct {
	XMLName    xml.Name `xml:"Topics" json:"-"`
	Topics     []Topic  `xml:"Topic" json:"topics"`
	NextMarker string   `xml:"NextMarker" json:"next_marker"`
}

type TopicNotification struct {
	XMLName		xml.Name	`xml:"Notification" json:"-"`
	TopicOwner	string		`xml:"TopicOwner" json:"topic_owner"`
	TopicName	string		`xml:"TopicName" json:"topic_name"`
	Subscriber	string		`xml:"Subscriber" json:"subscriber"`
	SubscriptionName	string	`xml:"SubscriptionName" json:"subscription_name"`
	MessageId	string		`xml:"MessageId" json:"message_id"`
	Message		Base64Bytes	`xml:"Message" json:"message"`
	MessageMD5	string		`xml:"MessageMD5" json:"message_md5"`
	MessageTag	string		`xml:"MessageTag" json:"message_tag"`
	PublishTime	int64		`xml:"PublishTime" json:"publish_time"`
}