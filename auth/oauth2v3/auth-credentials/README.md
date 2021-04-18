### 인증정보 가져오기
```bash
# 요청
curl \
-X GET \
'http://127.0.0.1:9096/credentials?login_id=admin&login_pw=admin'

# 응답
{"CLIENT_ID":"b16bb655-9568-4faa-83c0-4d141ed33035","CLIENT_SECRET":"51affcd8-32c1-4bb9-ac89-e27f9ebc30ad"}
```

### 액세스토큰 가져오기
```bash
# 요청
curl \
-X GET \
'http://127.0.0.1:9096/token?grant_type=client_credentials&client_id=b16bb655-9568-4faa-83c0-4d141ed33035&client_secret=51affcd8-32c1-4bb9-ac89-e27f9ebc30ad&scope=all'

# 응답
{"access_token":"GFKLV-PRMFYZKXP3SNNDXQ","expires_in":7200,"scope":"all","token_type":"Bearer"}
```

### validation
```bash
# 요청
curl \
-X GET \
'http://127.0.0.1:9096/validate?grant_type=client_credentials&access_token=PNN6GHA-OQM1YSWATN-VLA'

# 응답
{"access_token":"X63XCEDSOEWWJBXLRHDKSA","expires_in":7200,"scope":"all","token_type":"Bearer"}
```

### 리소스서버에 요청
```bash
curl \
-X GET \
'http://127.0.0.1:9097/protected?access_token=PNN6GHA-OQM1YSWATN-VLA'
```