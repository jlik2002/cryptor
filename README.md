# Welcome to Crypton!

## Examples request for this API

For upload file to server
```
 curl -X POST http://localhost:8070/upload \
 -F "file=@/home/jasmin/ex_go/cryptor/t123" \
 -F "passPhrase=test_pass" \
 -H "Content-Type: multipart/form-data"
```
For decrypt and download file
```
curl -X POST http://localhost:8070/decrypt \
-d '{"fileName":"t123","passPhrase":"test_pass"}' \
-H 'Content-Type: application/json'
```

For only download file
```
curl -X POST http://localhost:8070/download \
-d '{"fileName":"t123","passPhrase":"test_pass"}' \
-H 'Content-Type: application/json'
``` 