package injection

import (
	"encoding/json"
	"gomodules.xyz/jsonpatch/v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func InjectPod(originPod *corev1.Pod, originPodBytes []byte) ([]byte, error) {
	c := &corev1.Pod{}
	c = originPod

	//注解标识一下
	if c.Annotations == nil {
		c.Annotations = make(map[string]string)
	}
	annotations := c.Annotations
	annotations["syf"] = "operator"

	//注入container
	if newContainer, err := containerTemplate(); err != nil {
		return nil, err
	} else {
		c.Spec.Containers = append(c.Spec.Containers, *newContainer)
	}

	cBytes, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.CreatePatch(originPodBytes, cBytes)
	if err != nil {
		return nil, err
	}
	return json.Marshal(patch)
}

func PotentialPodName(metadata metav1.ObjectMeta) string {
	if metadata.Name != "" {
		return metadata.Name
	}
	if metadata.GenerateName != "" {
		return metadata.GenerateName + "***** (actual name not yet known)"
	}
	return ""
}

func containerTemplate() (*corev1.Container, error) {
	con := &corev1.Container{
		Name:            "syf-inject",
		Image:           "nginx:latest",
		ImagePullPolicy: corev1.PullIfNotPresent,
		SecurityContext: &corev1.SecurityContext{
			Privileged: func() *bool {
				b := true
				return &b
			}(),
		},
	}
	return con, nil
}
