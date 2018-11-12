package main

import (
    // "io"
    "net/http"
    "log"
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
)

type myhandler struct {
}

func (h *myhandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    log.Printf("ServeHTTP")

    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte("This is an example server.\n"))
    // fmt.Fprintf(w, "This is an example server.\n")
    // io.WriteString(w, "This is an example server.\n")

    log.Printf("ServeHTTP end")
}

func main() {

    pool := x509.NewCertPool()
    caCertPath := "../cert/server.crt"

    caCrt, err := ioutil.ReadFile(caCertPath)
    if err != nil {
      log.Println("ReadFile err:", err)
      return
    }
    pool.AppendCertsFromPEM(caCrt)

    s := &http.Server{
      Addr:    ":443",
      Handler: &myhandler{},
      TLSConfig: &tls.Config{
          ClientCAs:  pool,
          //ClientAuth: tls.RequireAndVerifyClientCert,
          ClientAuth: tls.NoClientCert, // https://golang.org/pkg/crypto/tls/
      },
    }

    err = s.ListenAndServeTLS("../cert/server.crt", "../cert/server.key")
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}