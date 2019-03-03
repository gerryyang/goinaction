package main
import (
	"log"
    "crypto/tls"
    "net/http"
    "net/http/httputil"
    "net/url"
    "crypto/x509"
    "io/ioutil"
)
func main() {
	log.SetFlags(log.Ldate | log.Ltime |log.Lshortfile)

	// x509.Certificate
	pool := x509.NewCertPool()
	caCertPath := "../cert/server.crt"

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		panic(err)
	}
	pool.AppendCertsFromPEM(caCrt)
	//pool.AddCert(caCrt)

	cliCrt, err := tls.LoadX509KeyPair("../cert/server.crt", "../cert/server.key")
	if err != nil {
        panic(err)
	}

    u, err := url.Parse("http://localhost:8888")
    if err != nil {
        panic(err)
    }

    tr := &http.Transport{
        Proxy: http.ProxyURL(u),
        // Disable HTTP/2.
        TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
        TLSClientConfig:    &tls.Config{
        	//InsecureSkipVerify: true,
        	RootCAs: pool,
        	Certificates: []tls.Certificate{cliCrt},
        	},
        DisableCompression: true,
    }

    client := &http.Client{Transport: tr}
    
    // 要求生成证书时指定 -subj /CN=localhost
    resp, err := client.Get("https://localhost/hello")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    dump, err := httputil.DumpResponse(resp, true)
    if err != nil {
        panic(err)
    }
    log.Printf("%q\n", dump)
}
