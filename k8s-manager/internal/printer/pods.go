package printer

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func PrintPods(namespace string, pods []corev1.Pod) {
	fmt.Printf("\nNamespace: %s\n", namespace)
	fmt.Printf("%-40s %-10s %-8s %s\n", "NAME", "STATUS", "RESTARTS", "AGE")
	fmt.Println(strings.Repeat("-", 75))

	for _, pod := range pods {
		restarts := int32(0)
		for _, cs := range pod.Status.ContainerStatuses {
			restarts += cs.RestartCount
		}

		age := ""
		if pod.Status.StartTime != nil {
			d := pod.Status.StartTime.Time
			_ = d
			age = pod.Status.StartTime.Time.Format("2006-01-02 15:04")
		}

		fmt.Printf("%-40s %-10s %-8d %s\n",
			pod.Name,
			string(pod.Status.Phase),
			restarts,
			age,
		)
	}
}
