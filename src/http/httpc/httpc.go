package httpc

import (
	//"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	//"encoding/xml"
	"fmt"
	//"going///attr"
	//"going/client"
	//"going/client/req"
	//"going/config"
	//"going///jm"
	//"going/json"
	"io"
	//"io/ioutil"
	"log"
	//"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	PostMethod = "POST"
	GetMethod  = "GET"
)

// Conf config for [http]
var Conf = struct {
	Http struct {
		Scheme         string          `default:"https"` // http协议类型 http https
		Address        string          // 请求地址 l5://1:2 cmlb://3 dns://api.qpay.qq.com:443
		Path           string          `default:"/v1/tm/portal"`// url路径 /index?k=v
		Host           string          // host
		Method         string          `default:"GET"`                             // GET POST
		ContentType    string          `default:"application/json; charset=utf-8"` // content type
		Proxy          string          // 代理 如 http://xxx.xxx.xxx.xxx:6000
		ClientCertFile string          // https需要  ClientCertFile ClientKeyFile RootCaFile 证书文件需要由服务方提供
		ClientKeyFile  string          // https需要
		RootCaFile     string          // https需要
		Timeout        time.Duration `default:"5000000000"`
		EnterAttr      int
		SuccAttr       int
		FailAttr       int
		CostAttr200    int
		CostAttr800    int
		CostAttr800p   int
	}
}{}

var once sync.Once

// httpTransport是全局仅需创建一次, 即可到处使用
var httpTransport *http.Transport

// Client httpc client
type Client struct {
	cli         *http.Client
	Scheme      string
	Method      string
	Host        string
	Cookie      string
	Address     string
	ContentType string
	Path        string
	Body        io.Reader
	Request     *http.Request
	Response    *http.Response
	header      map[string]string
	cmd         string
	addr        string
	err         error
	ModuleID    int //模调被调模块id
	InterfaceID int //模调被调接口id
	Timeout     time.Duration
	Cost        time.Duration

	EnterAttr    int
	SuccAttr     int
	FailAttr     int
	CostAttr200  int
	CostAttr800  int
	CostAttr800p int
}

// New 新建一个http请求客户端
func New() *Client {

	once.Do(func() {

		log.Printf("once do")

		// if https
		if true {

			log.Printf("into https")

			clientCertFile := "server.crt"
			clientKeyFile := "server.key"

			cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
			if err != nil {
				log.Printf("LoadX509KeyPair(%s,%s) error:%s\n", clientCertFile, clientKeyFile, err)
				return
			}

			/*certBytes, err := ioutil.ReadFile(Conf.Http.RootCaFile)
			if err != nil {
				log.Printf("Read rootCa fail, file:%s, err:%s\n", Conf.Http.RootCaFile, err)
				return
			}*/

			rootCAPool := x509.NewCertPool()

			/*ok := rootCAPool.AppendCertsFromPEM(certBytes)
			if !ok {
				log.Println("AppendCertsFromPEM fail")
				return
			}*/

			httpTransport = &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates:       []tls.Certificate{cert},
					RootCAs:            rootCAPool,
					InsecureSkipVerify: true,
				},
			}
		}

		// set ip  port
		/*log.Printf("set ip port")
		localAddr, err := net.ResolveIPAddr("ip", "10.56.58.55")
	    if err != nil {
	        panic(err)
	    }

	    localTCPAddr := net.TCPAddr{
	        IP: localAddr.IP,
	        Port: 30018,
	    }

	    d := net.Dialer{
	        LocalAddr: &localTCPAddr,
	        Timeout:   30 * time.Second,
	        KeepAlive: 30 * time.Second,
	    }

	    httpTransport = &http.Transport{
        	Dial: d.Dial,
    	}*/

		log.Printf("todo proxy")

		// 代理 如 http://xxx.xxx.xxx.xxx:6000
		if true {
			log.Printf("set proxy")

			proxy := "http://10.56.58.55:9090"
			httpProxy, err := url.Parse(proxy)
			if err != nil {
				log.Printf("parse proxy url:%s fail, err:%s\n", proxy, err)
				return
			}

			if httpTransport == nil {
				httpTransport = &http.Transport{
					Proxy: http.ProxyURL(httpProxy),
				}
			} else {
				httpTransport.Proxy = http.ProxyURL(httpProxy)
			}
		}
	})

	log.Printf("set client")
	o := &Client{
		Scheme:      "http",
		Method:      "GET",
		ContentType: "application/json; charset=utf-8",
		Timeout:     5 * time.Second,
	}

	o.Scheme       = Conf.Http.Scheme
	o.Method       = Conf.Http.Method
	o.Host         = Conf.Http.Host
	o.Address      = Conf.Http.Address
	o.ContentType  = Conf.Http.ContentType
	o.Path         = Conf.Http.Path

	//o.Timeout      = Conf.Http.Timeout
	//m, _ := time.ParseDuration("10s")
	//o.Timeout = m.Second()

	o.EnterAttr    = Conf.Http.EnterAttr
	o.SuccAttr     = Conf.Http.SuccAttr
	o.FailAttr     = Conf.Http.FailAttr
	o.CostAttr200  = Conf.Http.CostAttr200
	o.CostAttr800  = Conf.Http.CostAttr800
	o.CostAttr800p = Conf.Http.CostAttr800p
	
	o.cli = &http.Client{Timeout: o.Timeout}

	if httpTransport != nil {
		o.cli.Transport = httpTransport
	}
	log.Printf("set client over")

	return o
}

// Do do request
func Do(ctx context.Context) (*http.Response, error) {
	return New().Do(ctx)
}

// SetTransport set transport
func (c *Client) SetTransport(t *http.Transport) {
	c.cli.Transport = t
}


// Do do request
func (c *Client) Do(ctx context.Context) (*http.Response, error) {
	log.Printf("Do c.cli[%v]\n", c.cli)

	//attr.AttrAPI(c.EnterAttr, 1)
	c.cmd = fmt.Sprintf("%s.%s%s", c.Method, c.Host, c.Path)
	log.Printf("c.cmd[%s]\n", c.cmd)

	if c.err != nil {
		//attr.AttrAPI(c.FailAttr, 1)
		return nil, c.err
	}
	ctx, _ = context.WithTimeout(ctx, c.Timeout)
	d, _ := ctx.Deadline()
	begin := time.Now()
	c.Timeout = d.Sub(begin)
	c.cli.Timeout = c.Timeout

	/*addr, err := client.NewAddressing(c.Address)
	if err != nil {
		//attr.AttrAPI(c.FailAttr, 1)
		c.err = err
		return nil, err
	}*/

	/*defer func() {
		addr.Update(c.err)
		end := time.Now()
		c.Cost = end.Sub(begin)
		ec := "0"
		em := ""
		if c.err != nil {
			ec = c.err.Error()
			em = c.err.Error()
			//jm.Report(c.cmd, fmt.Sprintf("HTTP_%s", ec), c.cmd, fmt.Sprintf("HTTP_%s", ec), em, c.addr, c.Cost)
			//attr.AttrAPI(c.FailAttr, 1)
		} else {
			//jm.Report(c.cmd, ec, c.cmd, ec, em, c.addr, c.Cost)
			//attr.AttrAPI(c.SuccAttr, 1)
		}
		if c.Cost/time.Millisecond < 200 {
			//attr.AttrAPI(c.CostAttr200, 1)
		} else if c.Cost/time.Millisecond < 800 {
			//attr.AttrAPI(c.CostAttr800, 1)
		} else {
			//attr.AttrAPI(c.CostAttr800p, 1)
		}
	}()*/

	//c.addr = addr.Address()
	//query := fmt.Sprintf("%s://%s%s", c.Scheme, c.addr, c.Path)
	query := "http://10.56.58.55:30018/v1/tm/entry"

	c.Request, c.err = http.NewRequest(c.Method, query, c.Body)
	if c.err != nil {
		log.Printf("http.NewRequest err[%s]\n", c.err)
		return nil, c.err
	}
	c.Request.Host = "10.56.58.55"
	c.Request.Header.Set("Content-Type", c.ContentType)
	c.Request.Header.Set("Cookie", c.Cookie)
	for k, v := range c.header {
		c.Request.Header.Set(k, v)
	}
	c.Response, c.err = c.cli.Do(c.Request)
	if c.err != nil {
		log.Printf("c.cli.Do err[%s]\n", c.err)
		return nil, c.err
	}
	if c.Response.StatusCode >= 500 {
		c.err = fmt.Errorf("http status code:%d", c.Response.StatusCode)
	}

	log.Printf("return")
	return c.Response, c.err
}


// SetFormData set form data body to http request post body, ps: SetFormData(url.Values{"a":{"b"}}
func (c *Client) SetFormData(body url.Values) {
	log.Printf("SetFormData")
	c.Method = PostMethod
	c.ContentType = "application/x-www-form-urlencoded"

	c.Body = strings.NewReader(body.Encode())
}


// AddCookie add cookie
func (c *Client) AddCookie(key string, val string) {
	if len(c.Cookie) == 0 {
		c.Cookie = fmt.Sprintf("%s=%s", key, val)
	} else {
		c.Cookie = fmt.Sprintf("%s;%s=%s", c.Cookie, key, val)
	}
}

// SetHeader set header
func (c *Client) SetHeader(key string, val string) {
	if c.header == nil {
		c.header = make(map[string]string)
	}

	c.header[key] = val
}

// String debug string
func (c *Client) String() string {
	if c.err != nil {
		return fmt.Sprintf("http[%s], addr[%s], cost[%s], error[%v]", c.cmd, c.addr, c.Cost, c.err)
	}
	return fmt.Sprintf("http[%s], addr[%s], cost[%s]", c.cmd, c.addr, c.Cost)
}

// Success check error
func (c *Client) Success() bool {
	return c.err == nil
}

// SetTimeout set timeout like time.Second or 800 * time.Millisecond
func (c *Client) SetTimeout(d time.Duration) {
	c.Timeout = d
}

// SetAddress set address like ip://ip:port or l5://modid:cmdid or dns://kandian.qq.com:443
func (c *Client) SetAddress(s string) {
	c.Address = s
}
