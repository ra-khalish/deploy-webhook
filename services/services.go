package services

import (
	"log"
	"os/exec"
	"strings"
	"time"
)

func waiting() {
	log.Println("hold up...")
	time.Sleep(5 * time.Second)
	log.Println("go!")
}

func All(namespaces string) []byte {
	ns := namespaces
	log.Println("start execute kubectl get all command")
	cmd := exec.Command("kubectl", "get", "pods", "-n", ns, "-o", "name")
	stdout, err := cmd.Output()

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Completed")

	return stdout
}

func Apply(namespaces *string, path *string) string {
	waiting()
	setUrl := "https://minio.this.id/" + *path
	log.Printf(setUrl)
	log.Printf("start execute kubectl apply command to %s, with file path: %s", *namespaces, *path)
	cmd := exec.Command("kubectl", "apply", "-f", setUrl, "-n", *namespaces)
	stdout, err := cmd.Output()

	if err != nil {
		log.Fatal("fail execute command", err.Error())
	}
	check := strings.Contains(string(stdout), "configured")
	if check != true {
		log.Fatalf("apply namespace %s, want configured", string(stdout))
	}

	log.Println("completed")
	return string(stdout)
}

func Delete(namespaces *string) string {
	log.Println("start execute kubectl delete command")
	cmd := exec.Command("kubectl", "delete", "-f", "https://minio.this.id/deployment/ingress.yaml", "-n", *namespaces)
	stdout, err := cmd.Output()

	if err != nil {
		log.Fatalf("faild execute command %s", err.Error())
	}

	log.Println("completed")
	return string(stdout)
}
