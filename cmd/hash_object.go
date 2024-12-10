package cmd

import (
	"armanVersionControl/storage/objectstore"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	filePath string
	content  string
	write    bool
)
var hashObjectCmd = &cobra.Command{
	Use:   "hash-object [-w | --write] { {--file-path | -f} | {--content | -c} }",
	Short: "Computes the object ID value for a specified file path or provided content.",
	Long: "Computes the object ID value for a specified file path or provided content " +
		"and optionally writes the resulting object into the object database.\n" +
		"Reports its object ID into standard output.",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		s, err := func() (string, error) {
			c, err := func() ([]byte, error) {
				if content != "" {
					return []byte(content), nil
				}

				return os.ReadFile(filePath)
			}()

			if err != nil {
				return "", err
			}

			if write {
				return objectstore.Store(c)
			}

			return objectstore.ComputeHash(c), nil
		}()

		if err != nil {
			return err
		}

		fmt.Println(s)

		return nil
	},
}

func init() {
	hashObjectCmd.Flags().StringVarP(&filePath, "file-path", "f", "", "Hash content of file located at the given path.")
	hashObjectCmd.Flags().StringVarP(&content, "content", "c", "", "Hash the provided content from initCmd line input.")
	hashObjectCmd.Flags().BoolVarP(&write, "write", "w", false, "When specified will actually store the resulting object into object database.")

	RootCmd.AddCommand(hashObjectCmd)
}
