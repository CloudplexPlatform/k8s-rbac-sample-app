package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
	"os"
)

var port = "3550"

func mapEnv(target *string, envKey string) {
	v := os.Getenv(envKey)
	if v != "" {
		//panic(fmt.Sprintf("environment variable %q not set", envKey))
		*target = v
	}
}

func main() {

	mapEnv(&port, "PORT")
	g := gin.Default()
	g.GET("/permissions", getPermissions)
	panic(g.Run(fmt.Sprintf(":%s", port)))

}

func getPermissions(r *gin.Context) {
	config, err1 := rest.InClusterConfig()
	if err1 != nil {
		r.JSON(http.StatusInternalServerError, err1.Error())
		return
	}
	// creates the clientset
	clientset, err1 := kubernetes.NewForConfig(config)
	if err1 != nil {
		r.JSON(http.StatusInternalServerError, err1.Error())
		return
	}

	secrets, err := clientset.CoreV1().Secrets("default").List(metav1.ListOptions{})
	if errors.IsNotFound(err) {
		r.JSON(http.StatusNotFound, err.Error())
		return
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		r.JSON(http.StatusForbidden, statusError.ErrStatus.Message)
		return
	} else if err != nil {
		r.JSON(http.StatusInternalServerError, err.Error())
		return
	} else {
		r.JSON(http.StatusOK, fmt.Sprintf("There are %d secrets in the cluster\n", len(secrets.Items)))
		return
	}
}
