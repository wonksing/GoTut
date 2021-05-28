package main

import (
	"log"
	"strings"
)

func main() {
	log.Println("start")
	// act := [...]string{"read", "write", "update", "delete"}
	// scopes := [...]string{"item", "item:read", "item:write", "emp"}
	mapScopes := make(map[string]string)

	mapScopes["item"] = "/item,/item/new,/item/_add,/item/_delete"
	mapScopes["item:read"] = "/item,/item/new"
	mapScopes["emp"] = "/emp,/emp/new,/emp/_add"

	// clientScope := "item:read emp"
	clientScope := "item emp"

	// 부여받은 범위
	cs := strings.Split(clientScope, " ")
	// 접근 가능한 리소스와 행위
	var allowedResource []ClientResource
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
			allowedResource = append(allowedResource, ClientResource{res, tmpAct})
		}
	}

	log.Println(allowedResource)

	// url path가 주어졌을 때 접근 가능한 리소스인지 보기
	url := "/item"
	method := "POST"
	ok := false
	for _, val := range allowedResource {
		if url == val.Path {
			if method == "GET" && val.Action.Get {
				ok = true
			} else if method == "POST" && val.Action.Post {
				ok = true
			} else if method == "PUT" && val.Action.Put {
				ok = true
			} else if method == "DELETE" && val.Action.Delete {
				ok = true
			}
		}
	}

	if ok {
		log.Println("welcome", url, method)
	} else {
		log.Println("not authorized", url, method)
	}

}

type ClientResource struct {
	Path   string
	Action Action
}

type Action struct {
	Get    bool
	Post   bool
	Put    bool
	Delete bool
}
