

* Generate private key (.key)

```
# Key considerations for algorithm "RSA" ≥ 2048-bit
openssl genrsa -out server.key 2048
```

* Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)

```
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650 -subj /CN=localhost
```


