package commands

import (
    "bytes"
    "fmt"
    "os/exec"
    "strings"

    "github.com/spf13/cobra"
)

// GetDeployCommand returns the deploy command with docker and kubernetes targets
func GetDeployCommand() *cobra.Command {
    var (
        target      string
        environment string
        image       string
        tag         string
        namespace   string
        replicas    int
        ingressHost string
        ports       string
        detach      bool
        envFile     string
        rollback    bool
    )

    cmd := &cobra.Command{
        Use:   "deploy",
        Short: "Deploy CMP application",
        Long:  "Deploy the current CMP application to the configured environment.",
        RunE: func(cmd *cobra.Command, args []string) error {
            switch target {
            case "docker":
                return deployDocker(imageWithTag(image, tag), ports, detach, envFile)
            case "kubernetes":
                if rollback {
                    return kubectl("rollout", "undo", "deployment/contexis-app", "-n", namespace)
                }
                // Render and apply basic manifests
                files := []string{
                    "src/core/deployment/kubernetes/configmap.yaml",
                    "src/core/deployment/kubernetes/secret.yaml",
                    "src/core/deployment/kubernetes/deployment.yaml",
                    "src/core/deployment/kubernetes/service.yaml",
                }
                // Replace image if provided
                if image != "" || tag != "" {
                    img := imageWithTag(image, tag)
                    if err := kubectlWithImage(namespace, files, img); err != nil {
                        return err
                    }
                } else {
                    if err := kubectlApply(namespace, files); err != nil {
                        return err
                    }
                }
                // Optional ingress
                if ingressHost != "" {
                    if err := kubectlApply(namespace, []string{"src/core/deployment/kubernetes/ingress.yaml"}); err != nil {
                        return err
                    }
                }
                // Optional HPA
                _ = kubectlApply(namespace, []string{"src/core/deployment/kubernetes/hpa.yaml"})
                return nil
            default:
                return fmt.Errorf("unsupported target: %s", target)
            }
        },
    }

    cmd.Flags().StringVar(&target, "target", "docker", "Deployment target: docker|kubernetes")
    cmd.Flags().StringVar(&environment, "environment", "production", "Deployment environment")
    cmd.Flags().StringVar(&image, "image", "contexis-cmp/contexis", "Container image name")
    cmd.Flags().StringVar(&tag, "tag", "latest", "Image tag")
    cmd.Flags().StringVar(&namespace, "namespace", "default", "Kubernetes namespace")
    cmd.Flags().IntVar(&replicas, "replicas", 2, "Kubernetes replicas")
    cmd.Flags().StringVar(&ingressHost, "ingress-host", "", "Kubernetes ingress host (optional)")
    cmd.Flags().StringVar(&ports, "ports", "8000:8000", "Docker port mapping, e.g., 8000:8000")
    cmd.Flags().BoolVar(&detach, "detach", true, "Run Docker container in detached mode")
    cmd.Flags().StringVar(&envFile, "env-file", "", "Path to .env file for Docker (optional)")
    cmd.Flags().BoolVar(&rollback, "rollback", false, "Rollback last Kubernetes deployment")
    return cmd
}

func imageWithTag(image, tag string) string {
    if tag == "" || strings.Contains(image, ":") {
        return image
    }
    return fmt.Sprintf("%s:%s", image, tag)
}

func deployDocker(image, ports string, detach bool, envFile string) error {
    args := []string{"run", "--rm", "--name", "contexis", "-p", ports}
    if detach {
        args = append(args, "-d")
    }
    if envFile != "" {
        args = append(args, "--env-file", envFile)
    }
    args = append(args, image, "serve", "--addr", ":8000")
    return docker(args...)
}

func docker(args ...string) error {
    cmd := exec.Command("docker", args...)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("docker %v failed: %v\n%s", args, err, out.String())
    }
    return nil
}

func kubectl(args ...string) error {
    cmd := exec.Command("kubectl", args...)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("kubectl %v failed: %v\n%s", args, err, out.String())
    }
    return nil
}

func kubectlApply(namespace string, files []string) error {
    args := []string{"apply", "-n", namespace, "-f"}
    for _, f := range files {
        if err := kubectl(append(args, f)...); err != nil {
            return err
        }
    }
    return nil
}

func kubectlWithImage(namespace string, files []string, image string) error {
    // Simple strategy: apply manifests, then set image
    if err := kubectlApply(namespace, files); err != nil {
        return err
    }
    return kubectl("set", "image", "deployment/contexis-app", fmt.Sprintf("contexis=%s", image), "-n", namespace)
}


