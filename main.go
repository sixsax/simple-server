package main
import(
    "encoding/json"
	"net/http"
	log "github.com/sirupsen/logrus"
    "simple-server/source/controller"
)

func main(){
	//Create controller
    c := controller.New()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		//Handle action
		response, err := c.Execute(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal("%v",err)
			return
		}

		js, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal("%v",err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

		log.Info(string(js))

	})

	err := http.ListenAndServeTLS(":8443", "/kube_certificates/server.crt","/kube_certificates/server.key",nil)

	if err != nil {
		log.Fatal("%v",err)
	}
}

