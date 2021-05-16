### 토큰 요청
```bash
# 액세스 토큰 요청
curl \
-X GET \
'http://localhost:9096/token?grant_type=client_credentials&client_id=b16bb655&client_secret=b16bb655-9568-4faa-82c0-4d152ed33035&scope=all'

# 액세스 토큰 재발급 요청
curl \
-X GET \
'http://127.0.0.1:9096/token?grant_type=refresh_token&refresh_token=ZWFJZTFIMZKTZTCZMS01NWVILWEXMTQTOGZIY2I1NWRMYJY0&client_id=b16bb655&client_secret=b16bb655-9568-4faa-82c0-4d152ed33035'


curl \
-X GET \
'http://127.0.0.1:9096/token?grant_type=client_credentials&client_id=b16bb655&client_secret=b16bb655-9568-4faa-82c0-4d152ed33035&scope=all'

```

### validate 토큰
```bash

curl \
-X GET \
'http://127.0.0.1:9096/validate?grant_type=client_credentials&access_token=eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJiMTZiYjY1NSIsImV4cCI6MTYxNzYwOTA3Nn0.HeP-IslApwQu5SzvwhmsasMcdqOuOt_NuVgZbHfItYBiarU-OGJhDYKo1FGGMg5smgS3NXddNy4VLSVi1wHsIw'

curl \
-X GET \
'http://127.0.0.1:9096/validate?grant_type=client_credentials&access_token=O-ZGEF4WOJU4UYXO88RIXG'

```


curl \
-X GET \
'http://127.0.0.1:9097/protected?access_token=eyJhbGciOiJIUzUxMiIsImtpZCI6IlRFU1RfS0VZX0lEIiwidHlwIjoiSldUIn0.eyJhdWQiOiJiMTZiYjY1NSIsImV4cCI6MTYxNzYxNjE5MCwiZG9tYWluIjoiaHR0cDovL2xvY2FsaG9zdDo5MDk0Iiwic2NvcGUiOiJhbGwifQ.wQ9Q_kAkd44cx9YY-JODDYYgz6AZJm0ZGpIZzLWR6gKNeWxbjlejVD3c2dPh1Gnu1CAi2F49oHmksYAC8G_khg'


http://localhost:9096/oauth/authorize?client_id=222222&code_challenge=Qn3Kywp0OiU4NK_AFzGPlmrcYJDJ13Abj_jdL08Ahg8%3D&code_challenge_method=S256&redirect_uri=http%3A%2F%2Flocalhost%3A9094%2Foauth2&response_type=code&scope=all&state=xyz

http://localhost:9096/oauth/token?grant_type=client_credentials&client_id=222222&client_secret=22222222&scope=all

http://localhost:9096/oauth/authorize?grant_type=client_credentials&client_id=222222&client_secret=22222222&scope=all


curl \
-H 'Authorization: Basic MjIyMjIyOjIyMjIyMjIy' \
-X POST \
-d 'grant_type=client_credentials&client_id=222222&client_secret=22222222&scope=all' \
'http://localhost:9096/oauth/token'