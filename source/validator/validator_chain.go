package validator_chain

import(
	log "github.com/sirupsen/logrus"
	i "simple-server/source/includes"
)

type handler interface {
	Request(req i.RequestJSON) i.ResponseJSON
}

//ConfigMap test
type configMapTestHandler struct {
    next handler
}

func (h *configMapTestHandler) Request(req i.RequestJSON) i.ResponseJSON {
    log.Info("Handler configMapTestHandler")

    if(req["request"].(map[string]interface{})["kind"].(map[string]interface{})["kind"] != "ConfigMap"){
        if (h.next != nil){
            return h.next.Request(req)
        }
    }

    var res = i.ResponseJSON{
		"response": i.ResponseJSON{
			"status": i.ResponseJSON{
			},
		},
	}

    if ( h.analyze(req) ) {
        res["response"].(i.ResponseJSON)["allowed"] = true
        res["response"].(i.ResponseJSON)["status"].(i.ResponseJSON)["status"] = "Success"
    }else{
        res["response"].(i.ResponseJSON)["allowed"] = false
        res["response"].(i.ResponseJSON)["status"].(i.ResponseJSON)["status"] = "Failure"
        res["response"].(i.ResponseJSON)["status"].(i.ResponseJSON)["message"] = "ConfigMap not valid!"
        res["response"].(i.ResponseJSON)["status"].(i.ResponseJSON)["reason"] = "ConfigMap not valid!"
    }

    return res
}

//Test the resource
func (h *configMapTestHandler) analyze(req i.RequestJSON) bool {
    _, ok := req["request"].(map[string]interface{})["object"].(map[string]interface{})["data"].(map[string]interface{})["test"]
    return ok
}

//Default behaviour
type defaultBehaviourHandler struct {
        next handler
}

func (h *defaultBehaviourHandler) Request(req i.RequestJSON) i.ResponseJSON {
    log.Info("Handler defaultBehaviourHandler")
    log.Info(req)

    var res = i.ResponseJSON{
		"response": i.ResponseJSON{
			"status": i.ResponseJSON{
			},
		},
	}

    res["response"].(i.ResponseJSON)["allowed"] = false
    res["response"].(i.ResponseJSON)["status"].(i.ResponseJSON)["status"] = "Failure"
    res["response"].(i.ResponseJSON)["status"].(i.ResponseJSON)["message"] = "Resource is not recognized"
    res["response"].(i.ResponseJSON)["status"].(i.ResponseJSON)["reason"] = "Resource is not recognized"

    return res
}


type chain struct {
	req i.RequestJSON
    chain handler
}

func New(req i.RequestJSON) (*chain, error){

    log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

    ec := &configMapTestHandler{new(defaultBehaviourHandler)}
    //ec := &defaultBehaviourHandler{nil}
	c := chain{req,ec}
    return &c, nil
}

func (c chain) Execute() i.ResponseJSON{
	return c.chain.Request(c.req)
}

