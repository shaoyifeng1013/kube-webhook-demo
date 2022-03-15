package pkg

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	admiss "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	"syf-webhook/pkg/injection"
)

type Controller struct {
}

var (
	logger = log.Default()
)

func (c *Controller) Run(tlsPath string) error {
	g := gin.Default()

	g.POST("/inject", func(context *gin.Context) {
		bytes, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			context.JSON(http.StatusInternalServerError, &admiss.AdmissionResponse{
				Allowed: false,
				Result:  &metav1.Status{Message: err.Error()},
			})
			logger.Println(err.Error())
			return
		}
		logger.Println("into inject")
		var response *admiss.AdmissionReview
		response, err = c.executeReq(bytes)
		if err != nil {
			context.JSON(http.StatusOK, &admiss.AdmissionReview{
				Response: &admiss.AdmissionResponse{
					Allowed: false,
					Result:  &metav1.Status{Message: err.Error()},
				},
			})
			logger.Println(err.Error())
			return
		}
		context.JSON(http.StatusOK, response)
	})

	err := g.RunTLS(":443", tlsPath+"/tls.crt", tlsPath+"/tls.key")
	return err
}

func (c *Controller) executeReq(req []byte) (*admiss.AdmissionReview, error) {
	request := &admiss.AdmissionReview{}
	if err := json.Unmarshal(req, request); err != nil {
		return nil, err
	}
	logger.Println(string(request.Request.Object.Raw))
	reqPod := &corev1.Pod{}
	if err := json.Unmarshal(request.Request.Object.Raw, reqPod); err != nil {
		return nil, err
	}
	potentialPodName := injection.PotentialPodName(reqPod.ObjectMeta)
	logger.Printf("webhook prepare to inject container, pod name = %s, namespace = %s\n", potentialPodName, request.Request.Namespace)
	patch, err := injection.InjectPod(reqPod, request.Request.Object.Raw)
	if err != nil {
		return nil, err
	}
	return &admiss.AdmissionReview{
		TypeMeta: request.TypeMeta,
		Response: &admiss.AdmissionResponse{
			UID:     request.Request.UID,
			Allowed: true,
			Patch:   patch,
			PatchType: func() *admiss.PatchType {
				var pt admiss.PatchType = "JSONPatch"
				return &pt
			}(),
		},
	}, nil
}
