package webapi

import (
	"encoding/json"
	"net/http"

	typesv1 "github.com/10gen/dredd/crdapi/types/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (eh *WebAPIHandler) findMongoDBHandler(w http.ResponseWriter, r *http.Request) {
	zap.S().Debugf("GET /mongodbs/{name}")
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		msg := "No MongoDB deployment name was specified in the request"
		RespondWithError(w, http.StatusBadRequest, msg)
		zap.S().Warnf(msg)
		return
	}
	mongodb, err := eh.kubeClient.MongoDBs(eh.namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, &mongodb)
}

func (eh *WebAPIHandler) allMongoDBHandler(w http.ResponseWriter, r *http.Request) {
	zap.S().Debugf("GET /mongodbs")
	mongodbs, err := eh.kubeClient.MongoDBs(eh.namespace).List(metav1.ListOptions{})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, &mongodbs)
}

func (eh *WebAPIHandler) newMongoDBHandler(w http.ResponseWriter, r *http.Request) {
	zap.S().Debugf("POST /mongodbs")
	mongodb := typesv1.MongoDB{}
	err := json.NewDecoder(r.Body).Decode(&mongodb)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var result *typesv1.MongoDB
	result, err = eh.kubeClient.MongoDBs(eh.namespace).Create(&mongodb)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, &result)
}

func (eh *WebAPIHandler) deleteMongoDBHandler(w http.ResponseWriter, r *http.Request) {
	zap.S().Debugf("DELETE /mongodbs/{name}")
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		msg := "No MongoDB deployment name was specified in the request"
		RespondWithError(w, http.StatusBadRequest, msg)
		return
	}
	err := eh.kubeClient.MongoDBs(eh.namespace).Delete(name)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func InitialiseMongoDBRoutes(r *mux.Router, handler *WebAPIHandler) {
	mongodbsrouter := r.PathPrefix("/mongodbs").Subrouter()
	mongodbsrouter.Methods("GET").Path("/{name}").HandlerFunc(handler.findMongoDBHandler)
	mongodbsrouter.Methods("GET").Path("").HandlerFunc(handler.allMongoDBHandler)
	mongodbsrouter.Methods("POST").Path("").HandlerFunc(handler.newMongoDBHandler)
	mongodbsrouter.Methods("DELETE").Path("/{name}").HandlerFunc(handler.deleteMongoDBHandler)
}
