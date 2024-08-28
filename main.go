package main

import (
	"log"
	"sanjieke/config"
)

func main() {

	//if !checkInputConfig() {
	//	return
	//}

	config.Authorization = "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiIxMDAwMDAwMSIsImp0aSI6ImEwMGUxYWYyODNiNDY1NGNjZGJjZjYyOGY3ZWFjNjdkMzlkM2JhZTY5NTQyNzNhOTBiNjE1ZjEwZDk3M2EyOWJhOTY3N2E3M2JlZGIxZmYwIiwiaWF0IjoxNzI0ODQ4NzE0LCJuYmYiOjE3MjQ4NDg3MTQsImV4cCI6MTcyNjE0NDcxNCwic3ViIjoiMjMzODU4MjMiLCJpc3MiOiJzYW5qaWVrZS1vbmxpbmUiLCJzaWQiOiJhMDBlMWFmMjgzYjQ2NTRjY2RiY2Y2MjhmN2VhYzY3ZDM5ZDNiYWU2OTU0MjczYTkwYjYxNWYxMGQ5NzNhMjliYTk2NzdhNzNiZWRiMWZmMCIsInNjb3BlcyI6W119.Cqm9DPxkyRqXSqY3m6urVe8yu4gUDZRVr54SqcxBnJhpW3jzl96kGbeEvDzXHoxbuXSjlTaX-C0BhmFMzXOjUQ"
	config.Cookie = "sajssdk_2015_cross_new_user=1; PHPSESSID=3mms67kuht5c1la2mkeip7rcg3; _sjk_jwt=eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiIxMDAwMDAwMSIsImp0aSI6ImEwMGUxYWYyODNiNDY1NGNjZGJjZjYyOGY3ZWFjNjdkMzlkM2JhZTY5NTQyNzNhOTBiNjE1ZjEwZDk3M2EyOWJhOTY3N2E3M2JlZGIxZmYwIiwiaWF0IjoxNzI0ODQ4NzE0LCJuYmYiOjE3MjQ4NDg3MTQsImV4cCI6MTcyNjE0NDcxNCwic3ViIjoiMjMzODU4MjMiLCJpc3MiOiJzYW5qaWVrZS1vbmxpbmUiLCJzaWQiOiJhMDBlMWFmMjgzYjQ2NTRjY2RiY2Y2MjhmN2VhYzY3ZDM5ZDNiYWU2OTU0MjczYTkwYjYxNWYxMGQ5NzNhMjliYTk2NzdhNzNiZWRiMWZmMCIsInNjb3BlcyI6W119.Cqm9DPxkyRqXSqY3m6urVe8yu4gUDZRVr54SqcxBnJhpW3jzl96kGbeEvDzXHoxbuXSjlTaX-C0BhmFMzXOjUQ; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%2223385823%22%2C%22first_id%22%3A%2219198fe58e713f0-033f47b54599e18-26001f51-3686400-19198fe58e813fd%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%2C%22%24latest_referrer%22%3A%22%22%7D%2C%22identities%22%3A%22eyIkaWRlbnRpdHlfY29va2llX2lkIjoiMTkxOThmZTU4ZTcxM2YwLTAzM2Y0N2I1NDU5OWUxOC0yNjAwMWY1MS0zNjg2NDAwLTE5MTk4ZmU1OGU4MTNmZCIsIiRpZGVudGl0eV9sb2dpbl9pZCI6IjIzMzg1ODIzIn0%3D%22%2C%22history_login_id%22%3A%7B%22name%22%3A%22%24identity_login_id%22%2C%22value%22%3A%2223385823%22%7D%2C%22%24device_id%22%3A%2219198fe58e713f0-033f47b54599e18-26001f51-3686400-19198fe58e813fd%22%7D"
	config.ApiKey = "cDpJh7SuWGFZCFfSjvByc34PNSBrNVrB"
	config.CourseId = "34003473"
	config.OutDirectory = "./AAA"

	//检测目录是否存在
	err := ensureDirExists(config.OutDirectory)
	if err != nil {
		log.Fatalf("创建目录失败[%v]失败[%v]", config.OutDirectory, err.Error())
		return
	}

	err = loadCourse()
	if err != nil {
		log.Println("loadCourse 错误", err.Error())
		return
	}

}
