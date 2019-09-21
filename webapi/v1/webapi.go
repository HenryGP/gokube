package webapi

import (
	"encoding/json"
	"net/http"

	clientv1 "github.com/10gen/dredd/clientset/v1"

	"github.com/gorilla/mux"
)

type WebAPIHandler struct {
	kubeClient *clientv1.KubeClient
	namespace  string
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func ServeAPI(endpoint string, clientSet *clientv1.KubeClient, namespace string) error {
	handler := &WebAPIHandler{
		kubeClient: clientSet,
		namespace:  namespace,
	}
	router := mux.NewRouter()
	InitialiseMongoDBRoutes(router, handler)
	InitialiseCoreRoutes(router, handler)
	return http.ListenAndServe(endpoint, router)
}
