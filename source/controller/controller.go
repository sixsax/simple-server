package controller

import(
	"net/http"
    "net/http/httputil"
	"encoding/json"
	i "simple-server/source/includes"
	chain "simple-server/source/validator"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
)

type controller struct {
}

func New() controller{
    log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	c := controller{}
	return c
}

func (c controller) Execute(r *http.Request) (i.ResponseJSON, error) {

    requestDump, err := httputil.DumpRequest(r, true)
    if err != nil {
      return nil, err
    }

    log.Info(string(requestDump))

	action := r.URL.Path

	switch action{
		case "/validate":

            b, err := ioutil.ReadAll(r.Body)
            if err != nil {
                return nil, err
            }
            var reqJSON i.RequestJSON
            json.Unmarshal(b, &reqJSON)
            if err != nil {
                return nil, err
            }

			//The chain that will test the k8s resource
            ch, err := chain.New(reqJSON)
            if err != nil {
                return nil, err
            }

			//Execute the chain
			res := ch.Execute()

			return res, nil
		default:
			jsonRes := []byte(`{
		        "response":{
                        	"message": "Action not recognized!"
                	}
		        }`)

            var res i.ResponseJSON
            err := json.Unmarshal(jsonRes, &res)
            if err != nil {
                return nil, err
            }

            return res, nil
        }
}
