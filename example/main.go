package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/hongliang81/ali_mns"
	"fmt"
)

type appConf struct {
	Url             string `json:"url"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

func main() {
	conf := appConf{}

	if bFile, e := ioutil.ReadFile("app.conf"); e != nil {
		panic(e)
	} else {
		if e := json.Unmarshal(bFile, &conf); e != nil {
			panic(e)
		}
	}

	topic := ali_mns.NewMNSTopicManager("1340859151301362", conf.AccessKeyId, conf.AccessKeySecret)
	err := topic.CreateTopic(ali_mns.Beijing, "test", 65536)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = topic.SetTopicAttributes(ali_mns.Beijing, "test", 65534)
	if err != nil {
		fmt.Println(err.Error())
	}

	attr, err := topic.GetTopicAttributes(ali_mns.Beijing, "test")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%+v\n", attr)
	}

	//client := ali_mns.NewAliMNSClient(conf.Url,
	//	conf.AccessKeyId,
	//	conf.AccessKeySecret)

	//msg := ali_mns.MessageSendRequest{
	//	MessageBody:  []byte("hello gogap/ali_mns"),
	//	DelaySeconds: 0,
	//	Priority:     8}
	//
	//queue := ali_mns.NewMNSQueue("test", client)
	//ret, err := queue.SendMessage(msg)
	//
	//if err != nil {
	//	logs.Error(err)
	//} else {
	//	logs.Pretty("response:", ret)
	//}
	//
	//respChan := make(chan ali_mns.MessageReceiveResponse)
	//errChan := make(chan error)
	//go func() {
	//	for {
	//		select {
	//		case resp := <-respChan:
	//			{
	//				logs.Pretty("response:", resp)
	//				logs.Debug("change the visibility: ", resp.ReceiptHandle)
	//				if ret, e := queue.ChangeMessageVisibility(resp.ReceiptHandle, 5); e != nil {
	//					logs.Error(e)
	//				} else {
	//					logs.Pretty("visibility changed", ret)
	//				}
	//
	//				logs.Debug("delete it now: ", resp.ReceiptHandle)
	//				if e := queue.DeleteMessage(resp.ReceiptHandle); e != nil {
	//					logs.Error(e)
	//				}
	//			}
	//		case err := <-errChan:
	//			{
	//				logs.Error(err)
	//			}
	//		}
	//	}
	//
	//}()
	//
	//queue.ReceiveMessage(respChan, errChan)
	//for {
	//	time.Sleep(time.Second * 1)
	//}

}
//
//func main() {
//	conf := appConf{}
//
//	if bFile, e := ioutil.ReadFile("app.conf"); e != nil {
//		panic(e)
//	} else {
//		if e := json.Unmarshal(bFile, &conf); e != nil {
//			panic(e)
//		}
//	}
//
//	client := ali_mns.NewAliMNSClient(conf.Url,
//		conf.AccessKeyId,
//		conf.AccessKeySecret)
//
//	msg := ali_mns.MessageSendRequest{
//		MessageBody:  []byte("hello gogap/ali_mns"),
//		DelaySeconds: 0,
//		Priority:     8}
//
//	queue := ali_mns.NewMNSQueue("test", client)
//	ret, err := queue.SendMessage(msg)
//
//	if err != nil {
//		logs.Error(err)
//	} else {
//		logs.Pretty("response:", ret)
//	}
//
//	respChan := make(chan ali_mns.MessageReceiveResponse)
//	errChan := make(chan error)
//	go func() {
//		for {
//			select {
//			case resp := <-respChan:
//				{
//					logs.Pretty("response:", resp)
//					logs.Debug("change the visibility: ", resp.ReceiptHandle)
//					if ret, e := queue.ChangeMessageVisibility(resp.ReceiptHandle, 5); e != nil {
//						logs.Error(e)
//					} else {
//						logs.Pretty("visibility changed", ret)
//					}
//
//					logs.Debug("delete it now: ", resp.ReceiptHandle)
//					if e := queue.DeleteMessage(resp.ReceiptHandle); e != nil {
//						logs.Error(e)
//					}
//				}
//			case err := <-errChan:
//				{
//					logs.Error(err)
//				}
//			}
//		}
//
//	}()
//
//	queue.ReceiveMessage(respChan, errChan)
//	for {
//		time.Sleep(time.Second * 1)
//	}
//
//}