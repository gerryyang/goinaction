
在[go1.10](https://golang.org/doc/go1.10)的`net/http`支持client与Proxy通过https通信。

```
On the client side, an HTTP proxy (most commonly configured by ProxyFromEnvironment) can now be specified as an https:// URL, meaning that the client connects to the proxy over HTTPS before issuing a standard, proxied HTTP request. (Previously, HTTP proxy URLs were required to begin with http:// or socks5://.)
```

# Test

```
./https_proxy_hijacker -key server.key -pem server.crt -proto http
```


# 代理涉及的类型

1. Transport

Transport is an implementation of RoundTripper that supports HTTP, HTTPS, and HTTP proxies (for either HTTP or HTTPS with CONNECT).

https://golang.org/pkg/net/http/#Transport
https://golang.org/src/net/http/transport.go?s=3628:10127#L93


2. Response

https://golang.org/src/net/http/response.go?s=696:4094#L23


# 代理涉及接口

1. http

https://golang.org/pkg/net/http/

``` golang
// ProxyURL returns a proxy function (for use in a Transport) that always returns the same URL.
func ProxyURL(fixedURL *url.URL) func(*Request) (*url.URL, error)

// NewRequest returns a new Request given a method, URL, and optional body.
func NewRequest(method, url string, body io.Reader) (*Request, error)

// Do sends an HTTP request and returns an HTTP response, following policy (such as redirects, cookies, auth) as configured on the client.
func (c *Client) Do(req *Request) (*Response, error)
```



# refer

https://golang.org/pkg/crypto/tls/

https://github.com/denji/golang-tls

https://stackoverflow.com/questions/33768557/how-to-bind-an-http-client-in-go-to-an-ip-address

[HTTP(S) Proxy in Golang in less than 100 lines of code](https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c)

[HTTPS proxies support in Go 1.10](https://medium.com/@mlowicki/https-proxies-support-in-go-1-10-b956fb501d6b)

https://golang.org/pkg/net/http/#example_Hijacker

[TLS with Go](https://ericchiang.github.io/post/go-tls/)

http://www.01happy.com/https-principle-and-golang-practice/

https://github.com/denji/golang-tls

[golang GET 出现 x509: certificate signed by unknown authority](https://studygolang.com/articles/11175)

https://github.com/snail007/goproxy

