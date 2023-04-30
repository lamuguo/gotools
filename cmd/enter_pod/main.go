package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/exec"
)

func main() {
	var (
		kubeConfig string
		podName    string
		shellCmd   string
	)
	flag.StringVar(&kubeConfig, "kubeconfig", "/home/lamuguo/.kube/config", "path to the kubeconfig file")
	flag.StringVar(&podName, "pod", "", "pod name")
	flag.StringVar(&podName, "p", "", "pod name")
	flag.StringVar(&shellCmd, "command", "", "shell command")
	flag.StringVar(&shellCmd, "c", "sh", "shell command")
	flag.Parse()

	var config *rest.Config
	var err error

	if kubeConfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	pod, err := clientset.CoreV1().Pods("default").Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	var containerName string
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Ready {
			containerName = containerStatus.Name
			break
		}
	}

	if containerName == "" {
		fmt.Println("No ready containers found in the specified pod.")
		os.Exit(1)
	}

	execCmd := fmt.Sprintf("kubectl exec -it %s -c %s -- %s", podName, containerName, shellCmd)
	cmd := exec.Command("sh", "-c", execCmd)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting the command:", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error waiting for the command to finish:", err)
	}
}
