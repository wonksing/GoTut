# Gotut - http, filedownload

파일목록을 요청하고 파일을 한개씩 다운로드한다.

## 인증서 만들기

ssl 통신을 위해 사설 인증서를 만든다.

```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt
```

## curl 예제

api 테스트용 curl 커맨드

### 파일목록 가져오기

```bash
curl --insecure \
-H "Content-Type: application/json" \
-X GET \
'https://127.0.0.1:8080/testapi/file/download'
```

### 파일 다운로드하기

```bash
curl --insecure \
-H "Content-Type: application/json" \
-X GET \
--output ./downloaded/sample.txt \
'https://127.0.0.1:8080/testapi/file/download/sample.txt'
```
