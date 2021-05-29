package main

import (
	"errors"
	"log"
	"strings"
)

func main() {
	log.Println("start")
	// actions, 행위
	// read, update, write, delete, 4개의 행위로 구성
	// 각각 http method와 1:1로 매핑, (GET, POST, PUT, DELETE)

	// scopes, urn format?
	// 4개의 범위를 정의한다.
	// 	- item
	// 	- item:read
	// 	- item:write
	// 	- emp

	// map scopes to resources
	// 정의된 범위와 리소스 매핑
	mapScopes := make(map[string]string)
	mapScopes["item"] = "/item,/item/new,/item/_add,/item/_delete"
	mapScopes["item:read"] = "/item,/item/new"
	mapScopes["item:write"] = "/item/_add"
	mapScopes["emp"] = "/emp,/emp/new,/emp/_add"

	// 클라이언트에게 인가된 scopes
	clientScope := "item emp"

	// 접근 가능한 리소스와 행위
	authResources, err := GetAuthResources(mapScopes, clientScope)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(authResources)

	// url path가 주어졌을 때 접근 가능한 리소스인지 보기
	ok, err := IsAuthorized(authResources, "/item", "GET")
	if err != nil {
		log.Println(err)
		return
	}
	if ok {
		log.Println("welcome")
	} else {
		log.Println("you are not authorized")
	}

}

type AuthorizedResource struct {
	Path   string
	Action Action
}
type AuthorizedResources []AuthorizedResource
type Action struct {
	Get    bool
	Post   bool
	Put    bool
	Delete bool
}

func GetAuthResources(mapScopes map[string]string, clientScope string) (*AuthorizedResources, error) {
	if clientScope == "" {
		return nil, errors.New("no scope")
	}

	cs := strings.Split(clientScope, " ")

	var authResources AuthorizedResources
	for _, val := range cs {
		tmp := strings.Split(val, ":")
		tmpAct := Action{}
		if len(tmp) == 1 {
			tmpAct.Get = true
			tmpAct.Post = true
			tmpAct.Put = true
			tmpAct.Delete = true
		} else {
			switch tmp[1] {
			case "read":
				tmpAct.Get = true
			case "update":
				tmpAct.Post = true
			case "write":
				tmpAct.Put = true
			case "delete":
				tmpAct.Delete = true
			}
		}
		resources := strings.Split(mapScopes[tmp[0]], ",")
		for _, res := range resources {
			authResources = append(authResources, AuthorizedResource{res, tmpAct})
		}
	}

	return &authResources, nil
}

func IsAuthorized(authResources *AuthorizedResources, path, method string) (bool, error) {
	if authResources == nil {
		return false, errors.New("no authorized resources")
	}
	if path == "" {
		return false, errors.New("no path")
	}

	isAuthorized := false
	for _, res := range *authResources {
		if path == res.Path {
			if method == "GET" && res.Action.Get {
				isAuthorized = true
			} else if method == "POST" && res.Action.Post {
				isAuthorized = true
			} else if method == "PUT" && res.Action.Put {
				isAuthorized = true
			} else if method == "DELETE" && res.Action.Delete {
				isAuthorized = true
			}
		}
	}

	return isAuthorized, nil
}
