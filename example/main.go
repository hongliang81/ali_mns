package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/hongliang81/ali_mns"
	"fmt"
)

type appConf struct {
	Url             string `json:"url"`
	AccountId	string `json:"accountid"`
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

	// Topic Management

	topicManager := ali_mns.NewMNSTopicManager("1340859151301362", conf.AccessKeyId, conf.AccessKeySecret)
	err := topicManager.CreateTopic(ali_mns.Beijing, "test", 65536)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = topicManager.SetTopicAttributes(ali_mns.Beijing, "test", 65534)
	if err != nil {
		fmt.Println(err.Error())
	}

	attr, err := topicManager.GetTopicAttributes(ali_mns.Beijing, "test")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%+v\n", attr)
	}

	topics, err := topicManager.ListTopic(ali_mns.Beijing, "", 0, "te")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%+v\n", topics)
	}

	err = topicManager.DeleteTopic(ali_mns.Beijing, "test")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Topic [test] deleted")
	}

	// Topic Subscription
	err = topicManager.CreateTopic(ali_mns.Beijing, "testSub", 65536)

	err = topicManager.Subscribe(ali_mns.Beijing,
		"testSub",
		"",
		string("mail:directmail:{hongliang@neusoft.com}"),
		"testSub")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Subscription created successfully")
	}

	client := ali_mns.NewAliMNSClient(conf.Url,
		conf.AccessKeyId,
		conf.AccessKeySecret)

	topic := ali_mns.NewMNSTopic("testSub", client)

	msg := ali_mns.TopicMessageSendRequest{
		MessageBody:	[]byte("hello ali_mns"),
		MessageTag:	"",
	}
	resp, err := topic.SendMessage(msg)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%+v\n", resp)
	}


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