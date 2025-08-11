package commands

import (
    "bytes"
    "fmt"
    "os/exec"
    "strings"

    "github.com/spf13/cobra"
)

// GetBuildCommand returns the build command to build and optionally push a Docker image
func GetBuildCommand() *cobra.Command {
    var (
        environment string
        image       string
        tag         string
        platform    string
        push        bool
        buildArgs   []string
    )

    cmd := &cobra.Command{
        Use:   "build",
        Short: "Build container image for deployment",
        RunE: func(cmd *cobra.Command, args []string) error {
            img := image
            if tag != "" && !strings.Contains(img, ":") {
                img = fmt.Sprintf("%s:%s", image, tag)
            }
            dockerArgs := []string{"build", "-t", img, "."}
            if platform != "" {
                dockerArgs = append([]string{"buildx", "build", "--platform", platform, "-t", img, "."})
            }
            for _, ba := range buildArgs {
                dockerArgs = append([]string{"build", "-t", img}, parseBuildArg(ba)...)
                dockerArgs = append(dockerArgs, ".")
            }
            if err := runDocker(dockerArgs...); err != nil {
                return err
            }
            if push {
                if err := runDocker("push", img); err != nil {
                    return err
                }
            }
            fmt.Printf("IMAGE=%s TAG=%s\n", img, tag)
            _ = environment // reserved for future use in Dockerfile ARGs
            return nil
        },
    }

    cmd.Flags().StringVar(&environment, "environment", "production", "Build environment")
    cmd.Flags().StringVar(&image, "image", "contexis-cmp/contexis", "Image name")
    cmd.Flags().StringVar(&tag, "tag", "latest", "Image tag")
    cmd.Flags().StringVar(&platform, "platform", "", "Target platform, e.g., linux/amd64,linux/arm64 (uses buildx)")
    cmd.Flags().BoolVar(&push, "push", false, "Push image after build")
    cmd.Flags().StringArrayVar(&buildArgs, "build-arg", nil, "Build arguments (KEY=VALUE)")
    return cmd
}

func runDocker(args ...string) error {
    cmd := exec.Command("docker", args...)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("docker %v failed: %v\n%s", args, err, out.String())
    }
    return nil
}

func parseBuildArg(s string) []string {
    if s == "" {
        return nil
    }
    return []string{"--build-arg", s}
}


