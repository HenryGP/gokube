package webapi

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	apiv1 "k8s.io/api/core/v1"
)

type ConfigMapBody struct {
	ProjectName string `json:"projectName"`
	OrgID       string `json:"orgId"`
	BaseURL     string `json:"baseUrl"`
}

type SecretBody struct {
	SecretName string `json:"secretName"`
	ApiUser    string `json:"apiUser"`
	ApiKey     string `json:"apiKey"`
}

func (eh *WebAPIHandler) findCoreHandler(w http.ResponseWriter, r *http.Request) {
	zap.S().Debugf("GET /core/{component}/{name}")
	vars := mux.Vars(r)
	component, ok := vars["component"]
	if !ok {
		msg := "No component was specified in the request, this should either be 'configmap' or 'secret'"
		RespondWithError(w, http.StatusBadRequest, msg)
		zap.S().Warnf("%s", msg)
		return
	}
	var name string
	name, ok = vars["name"]
	if !ok {
		msg := "No name was specified in the request"
		RespondWithError(w, http.StatusBadRequest, msg)
		zap.S().Warnf("%s", msg)
		return
	}
	switch component {
	case "configmap":
		configMap, err := eh.kubeClient.Core(eh.namespace).GetConfigMap(name)
		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, &configMap)
	case "secret":
		secret, err := eh.kubeClient.Core(eh.namespace).GetSecret(name)
		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, &secret)
	default:
		msg := "Resource doesn't exist"
		RespondWithError(w, http.StatusBadRequest, msg)
		zap.S().Infof(msg)
		return
	}
	_, ok = vars["name"]
	if !ok {
		msg := "No name was specified in the request"
		RespondWithError(w, http.StatusBadRequest, msg)
		zap.S().Warnf(msg)
		return
	}
}

func (eh *WebAPIHandler) allCoreHandler(w http.ResponseWriter, r *http.Request) {
	zap.S().Debugf("GET /core/{component}")
	vars := mux.Vars(r)
	component, ok := vars["component"]
	if !ok {
		msg := "No component was specified in the request, this should either be 'configmap' or 'secret'"
		RespondWithError(w, http.StatusBadRequest, msg)
		zap.S().Warnf(msg)
		return
	}
	switch component {
	case "configmap":
		configMapList, err := eh.kubeClient.Core(eh.namespace).GetConfigMaps()
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, &configMapList)
	case "secret":
		secretList, err := eh.kubeClient.Core(eh.namespace).GetSecrets()
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, &secretList)
	default:
		msg := "Resource doesn't exist"
		RespondWithError(w, http.StatusBadRequest, msg)
		zap.S().Infof(msg)
		return
	}
}

func (eh *WebAPIHandler) newCoreHandler(w http.ResponseWriter, r *http.Request) {
	zap.S().Debugf("POST /core/{component}")
	vars := mux.Vars(r)
	component, ok := vars["component"]
	if !ok {
		msg := "No component was specified in the request, this should either be 'configmap' or 'secret'"
		RespondWithError(w, http.StatusBadRequest, msg)
		zap.S().Warnf(msg)
		return
	}
	switch component {
	case "configmap":
		cfgmap := ConfigMapBody{}
		err := json.NewDecoder(r.Body).Decode(&cfgmap)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		var result *apiv1.ConfigMap
		result, err = eh.kubeClient.Core(eh.namespace).CreateConfigMap(cfgmap.ProjectName, cfgmap.OrgID, cfgmap.BaseURL)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, &result)
	case "secret":
		secret := SecretBody{}
		err := json.NewDecoder(r.Body).Decode(&secret)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		var result *apiv1.Secret
		result, err = eh.kubeClient.Core(eh.namespace).CreateSecret(secret.SecretName, secret.ApiUser, secret.ApiKey)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, &result)
	default:
		msg := "Resource doesn't exist"
		RespondWithError(w, http.StatusBadRequest, msg)
		zap.S().Infof(msg)
		return
	}
}

func (eh *WebAPIHandler) deleteCoreHandler(w http.ResponseWriter, r *http.Request) {
	zap.S().Debugf("DELETE /core/{component}/{name}")
	vars := mux.Vars(r)
	component, ok := vars["component"]
	if !ok {
		msg := "No component was specified in the request, this should either be 'configmap' or 'secret'"
		RespondWithError(w, http.StatusBadRequest, msg)
		zap.S().Warnf(msg)
		return
	}
	name, ok := vars["name"]
	if !ok {
		msg := "No name was specified in the request"
		RespondWithError(w, http.StatusBadRequest, msg)
		zap.S().Warnf(msg)
		return
	}
	switch component {
	case "configmap":
		err := eh.kubeClient.Core(eh.namespace).DeleteConfigMap(name)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
	case "secret":
		err := eh.kubeClient.Core(eh.namespace).DeleteSecret(name)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
	default:
		msg := "Resource doesn't exist"
		RespondWithError(w, http.StatusBadRequest, msg)
		zap.S().Infof(msg)
		return
	}
}

func InitialiseCoreRoutes(r *mux.Router, handler *WebAPIHandler) {
	corerouter := r.PathPrefix("/core/{component}").Subrouter()
	corerouter.Methods("GET").Path("/{name}").HandlerFunc(handler.findCoreHandler)
	corerouter.Methods("GET").Path("").HandlerFunc(handler.allCoreHandler)
	corerouter.Methods("POST").Path("").HandlerFunc(handler.newCoreHandler)
	corerouter.Methods("DELETE").Path("/{name}").HandlerFunc(handler.deleteCoreHandler)
}
