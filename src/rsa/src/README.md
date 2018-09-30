
# 公私钥生成

* 生成RSA私钥

openssl genrsa -out privatekey.pem 1024

* 从RSA私钥导出公钥

openssl rsa -in privatekey.pem -out public.pem -outform PEM -pubout

* 转换为pkcs8格式私钥

openssl pkcs8 -topk8 -inform PEM -in privatekey.pem -outform PEM -nocrypt

# refer

1. https://blog.csdn.net/xz_studying/article/details/80314111
2. https://segmentfault.com/a/1190000004151272
3. https://github.com/89hmdys/toast
4. https://studygolang.com/articles/5257
