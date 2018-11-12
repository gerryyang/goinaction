package httpc_test

import (
	"context"
	"fmt"
	"gerry/httpc"
	"net/url"
)

//[http]
//address="l5://543745:327680"
//host="api.qpay.qq.com"
//path="/cgi-bin/hongbao/qpay_hb_mch_send.cgi"

var ctx = context.Background()

func Example() {

	o := httpc.New()
	o.SetFormData(url.Values{"gerry": {"1", "2", "3"}})

	o.Do(ctx)

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
