package hawkservice

import "net/http"

func Init() {

	
	http.HandleFunc("/test", test)
	http.HandleFunc("/compareEnv",compareEnv)
	

}
