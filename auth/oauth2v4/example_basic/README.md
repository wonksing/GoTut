# OAuth2

## 테스트 스크립트

```bash
# 액세스 토큰 요청
# base64 of client_id:client_secret
curl \
-X GET \
-H 'Authorization: Basic MTIzNDU6MTIzNDU2Nzg=' \
'http://localhost:9096/oauth/token?grant_type=client_credentials&scope=all'

curl \
-X GET \
-H 'Authorization: Basic YWJjZGU6YWJjZDEyMzQ=' \
'http://localhost:9096/oauth/token?grant_type=client_credentials&scope=all'

curl \
-X GET \
-H 'Authorization: Basic YXNkZjpxd2Vy' \
'http://localhost:9096/oauth/token?grant_type=client_credentials&scope=all'

# 액세스 토큰 재발급 요청
curl \
-X GET \
-H 'Authorization: Basic MTIzNDU6MTIzNDU2Nzg=' \
'http://127.0.0.1:9096/oauth/token?grant_type=refresh_token&refresh_token=ZDC4NTK1MZKTMJKWMS01Y2ZILTK3NGYTYTU2MDJHMZZHYJYY'

```

### validate 토큰
```bash
curl \
-X GET \
'http://127.0.0.1:9096/test?access_token=OTCZY2M4MGETYJBIOC0ZZDY5LWE0Y2ITYZVKMTHIMWFIZJJM'

```

### client credential 추가
```bash
curl \
-X PUT \
-d 'client_id=asdf&client_secret=qwer&client_domain=localhost:8080' \
'http://127.0.0.1:9096/credentials'

```

```bash
curl \
-X GET \
'http://localhost:9096/oauth/authorize2?client_id=12345&code_challenge=Qn3Kywp0OiU4NK_AFzGPlmrcYJDJ13Abj_jdL08Ahg8%3D&code_challenge_method=S256&redirect_uri=http%3A%2F%2Flocalhost%3A9094%2Foauth2&response_type=code&scope=all&state=xyz'
```