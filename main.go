package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/dariubs/percent"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getPods(c *kubernetes.Clientset, labelKey string, labelValue string, namespace string) (*v1.PodList, error) {
	labelSelector := fmt.Sprintf("%s=%s", labelKey, labelValue)
	opts := metav1.ListOptions{
		LabelSelector: labelSelector,
		Limit:         500,
	}

	pods, err := c.CoreV1().Pods(namespace).List(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	return pods, nil
}

/* getGroups defines number of nodes:
     g1 = 10% of nodes
     g2 ~ 30% of nodes
     g3 ~ 30% of nodes
     g4 ~ 30% of nodes
Total must be equal with total number of nodes */
func getLabelGroup(pods int, index int) (label string) {
	label = "group1"
	g1, g2, g3 := 0, 0, 0

	// NOTE: Percentages could be arguments
	if pods > 9 {
		g1 = int(percent.Percent(10, pods))
		g2 = int(percent.Percent(30, pods)) + g1
		g3 = int(percent.Percent(30, pods)) + g2
	} else {
		// if pods less than 10, don't bother with calculations
		// return just one group
		g1, g2, g3 = pods, 0, 0
	}

	if index > g1 {
		label = "group2"
	}
	if index > g2 {
		label = "group3"
	}
	if index > g3 {
		label = "group4"
	}
	return label
}

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	// variables
	namespace := flag.String("namespace", "sfu", "kubernetes namespace, default: 'sfu'")
	labelKey := flag.String("key", "app", "existing pod label key identifier, default: 'app'")
	labelValue := flag.String("value", "sfu", "existing pod label value identifier, default: 'sfu'")
	podLabelKey := flag.String("label", "deploy", "new pod label key identifier. Values are: group1, group2, ..., groupN, default 'deploy'")

	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := getPods(c, *labelKey, *labelValue, *namespace)
	if err != nil {
		panic(err.Error())
	}

	numberOfPods := len(pods.Items)

	fmt.Printf("Found %d pods featuring label '%s=%s' running on %s namespace\n", numberOfPods, *labelKey, *labelValue, *namespace)

	for i, pod := range pods.Items {
		podLabelValue := getLabelGroup(numberOfPods, i)
		patch := fmt.Sprintf(`[{"op":"add","path":"/metadata/labels/%s","value":"%s" }]`, *podLabelKey, podLabelValue)
		_, err := c.CoreV1().Pods(*namespace).Patch(context.TODO(), pod.Name, types.JSONPatchType, []byte(patch), metav1.PatchOptions{FieldManager: "JsonPatch"})
		if err == nil {
			fmt.Println(fmt.Sprintf("Label %s=%s added to pod %s", *podLabelKey, podLabelValue, pod.Name))
		} else {
			fmt.Println(err)
		}
	}
}
