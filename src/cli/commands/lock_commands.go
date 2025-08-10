package commands

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	runtimecontext "github.com/contexis-cmp/contexis/src/runtime/context"
	"github.com/spf13/cobra"
)

type LockFile struct {
	Contexts map[string]string            `json:"contexts"` // name -> sha
	Prompts  map[string]map[string]string `json:"prompts"`  // component -> relPath -> sha
	Memory   map[string]string            `json:"memory"`   // component -> sha of content files
}

func GetLockCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "lock", Short: "Generate context.lock.json for reproducibility"}
	cmd.AddCommand(newLockGenerateCmd())
	return cmd
}

func newLockGenerateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Compute SHAs for contexts, prompts, and memory",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, _ := os.Getwd()
			lock := LockFile{Contexts: map[string]string{}, Prompts: map[string]map[string]string{}, Memory: map[string]string{}}

			// Contexts: directories under contexts/ (skip tenants)
			ctxSvc := runtimecontext.NewContextService(root)
			entries, _ := os.ReadDir(filepath.Join(root, "contexts"))
			for _, e := range entries {
				if !e.IsDir() {
					continue
				}
				if e.Name() == "tenants" {
					continue
				}
				name := e.Name()
				ctxModel, err := ctxSvc.ResolveContext("", name)
				if err != nil {
					continue
				}
				sha, _ := ctxModel.GetSHA()
				lock.Contexts[name] = sha
			}

			// Prompts: compute file shas per component
			promptsDir := filepath.Join(root, "prompts")
			_ = filepath.WalkDir(promptsDir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if d.IsDir() {
					return nil
				}
				rel, _ := filepath.Rel(promptsDir, path)
				parts := strings.Split(rel, string(filepath.Separator))
				if len(parts) < 2 {
					return nil
				}
				comp := parts[0]
				by, err := os.ReadFile(path)
				if err != nil {
					return nil
				}
				h := sha256.Sum256(by)
				if _, ok := lock.Prompts[comp]; !ok {
					lock.Prompts[comp] = map[string]string{}
				}
				lock.Prompts[comp][filepath.ToSlash(rel)] = hex.EncodeToString(h[:])
				return nil
			})

			// Memory: hash files under memory/<component>/ (non-recursive summary)
			memDir := filepath.Join(root, "memory")
			comps, _ := os.ReadDir(memDir)
			for _, c := range comps {
				if !c.IsDir() {
					continue
				}
				compDir := filepath.Join(memDir, c.Name())
				var files []string
				filepath.WalkDir(compDir, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return nil
					}
					if d.IsDir() {
						return nil
					}
					files = append(files, path)
					return nil
				})
				sort.Strings(files)
				h := sha256.New()
				for _, f := range files {
					by, err := os.ReadFile(f)
					if err != nil {
						continue
					}
					h.Write([]byte(f))
					h.Write([]byte{0})
					h.Write(by)
				}
				lock.Memory[c.Name()] = hex.EncodeToString(h.Sum(nil))
			}

			// Write lock file
			out := filepath.Join(root, "context.lock.json")
			by, _ := json.MarshalIndent(lock, "", "  ")
			if err := os.WriteFile(out, by, 0o644); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "wrote %s\n", out)
			return nil
		},
	}
}
