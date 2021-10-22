package main


import (
	"net/http"
	"./hawkservice"
)


func main(){


	hawkservice.Init()
	http.ListenAndServe(":4040", nil)



}
