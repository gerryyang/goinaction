package main

import (
	"context"
	//"fmt"
	"gerry/httpc"
	"net/url"
	"log"
)

//[http]
//address="l5://543745:327680"
//host="api.qpay.qq.com"
//path="/cgi-bin/hongbao/qpay_hb_mch_send.cgi"

var ctx = context.Background()

func main() {

	log.SetFlags(log.Ldate | log.Ltime |log.Lshortfile)

	o := httpc.New()
	o.SetFormData(url.Values{"gerry": {"1", "2", "3"}})

	rsp, err := o.Do(ctx)
	if err != nil {
		log.Printf("http Do error:%s\n", err)
		return
	}

	log.Printf("rsp[%v]\n", rsp)

	/*var res struct {
		Retcode int
		Retmsg  string
		Result  map[string]interface{} //自己根据业务需求定义返回数据结构体
	}
	err := o.GetJSONBody(&res)
	fmt.Println("hello http")
	fmt.Println(o)
	fmt.Println(res)
	fmt.Println(err)*/
	//output:
	//hello http
}
