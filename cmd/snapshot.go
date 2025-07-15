package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var namespace string
var allNamespaces bool

// snapshotCmd represents the snapshot command
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Snapshot Kubernetes cluster resources",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("📂 Using namespace: %s\n", namespace)
		// Load kubeconfig (update if needed)
		kubeconfig := "/home/sameerk/.kube/config"
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("❌ Failed to load kubeconfig: %v", err)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("❌ Failed to create clientset: %v", err)
		}

		targetNamespace := namespace
		if allNamespaces {
			targetNamespace = ""
			fmt.Println("🌍 Snapshotting ALL namespaces")
		} else {
			fmt.Printf("📂 Using namespace: %s\n", targetNamespace)
		}

		// Create timestamped snapshot directory
		timestamp := time.Now().Format("2006-01-02_150405")
		dir := filepath.Join("snapshots", timestamp)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("❌ Failed to create directory: %v", err)
		}

		// Fetch Pods
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("❌ Failed to list pods: %v", err)
		}
		saveYAML(pods, filepath.Join(dir, "pods.yaml"))

		// Fetch Deployments
		deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("❌ Failed to list deployments: %v", err)
		}
		saveYAML(deployments, filepath.Join(dir, "deployments.yaml"))

		// Fetch Services
		services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("❌ Failed to list services: %v", err)
		}
		saveYAML(services, filepath.Join(dir, "services.yaml"))

		// Fetch ConfigMaps
		configMaps, err := clientset.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("❌ Failed to list configmaps: %v", err)
		}
		saveYAML(configMaps, filepath.Join(dir, "configmaps.yaml"))

		fmt.Printf("✅ Snapshot saved to folder: %s\n", dir)
	},
}

func saveYAML(obj interface{}, path string) {
	data, err := yaml.Marshal(obj)
	if err != nil {
		log.Printf("⚠️ Failed to marshal YAML for %s: %v", path, err)
		return
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		log.Printf("⚠️ Failed to write file %s: %v", path, err)
	}
}

func init() {
	snapshotCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Namespace to snapshot")
	snapshotCmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "Snapshot all namespaces")
	rootCmd.AddCommand(snapshotCmd)

}
