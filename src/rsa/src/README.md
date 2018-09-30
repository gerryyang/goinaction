
生成RSA私钥

openssl genrsa -out privatekey.pem 1024

从RSA私钥导出公钥

openssl rsa -in privatekey.pem -out public.pem -outform PEM -pubout

转换为pkcs8格式私钥

openssl pkcs8 -topk8 -inform PEM -in privatekey.pem -outform PEM -nocrypt

