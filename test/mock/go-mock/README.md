# Gotut - gomock

## golang/mock 설치하기

```bash
GO111MODULE=on go get github.com/golang/mock/mockgen@v1.5.0
```

## mock 만들기

패키지 경로로 이동해서 만드는 것이 관리하기에 편할 것 같다. 이 예제의 패키지 경로는 test/mock/go-mock 으로 하자.

```bash
mockgen -destination=port/user/mocks/mock_user_repo.go -package=mocks -mock_names="Repository=MockUserRepository" -source=port/user/user_repo.go port/user Repository
```