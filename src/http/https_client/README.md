



# test

```
$curl -kv "https://127.0.0.1/hello"
* Hostname was NOT found in DNS cache
*   Trying 127.0.0.1...
* Connected to 127.0.0.1 (127.0.0.1) port 443 (#0)
* successfully set certificate verify locations:
*   CAfile: /etc/pki/tls/certs/ca-bundle.crt
  CApath: none
* SSLv3, TLS handshake, Client hello (1):
* SSLv3, TLS handshake, Server hello (2):
* SSLv3, TLS handshake, CERT (11):
* SSLv3, TLS handshake, Server key exchange (12):
* SSLv3, TLS handshake, Server finished (14):
* SSLv3, TLS handshake, Client key exchange (16):
* SSLv3, TLS change cipher, Client hello (1):
* SSLv3, TLS handshake, Finished (20):
* SSLv3, TLS change cipher, Client hello (1):
* SSLv3, TLS handshake, Finished (20):
* SSL connection using TLSv1.2 / ECDHE-RSA-AES128-GCM-SHA256
* Server certificate:
*        subject: C=XX; L=Default City; O=Default Company Ltd
*        start date: 2018-11-12 03:54:53 GMT
*        expire date: 2028-11-09 03:54:53 GMT
*        issuer: C=XX; L=Default City; O=Default Company Ltd
*        SSL certificate verify result: self signed certificate (18), continuing anyway.
> GET /hello HTTP/1.1
> User-Agent: curl/7.38.0
> Host: 127.0.0.1
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: text/plain
< Date: Mon, 12 Nov 2018 05:12:44 GMT
< Content-Length: 27
< 
This is an example server.
* Connection #0 to host 127.0.0.1 left intact
```

