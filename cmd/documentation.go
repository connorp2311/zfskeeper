/*
Copyright © 2023 Connor Parsons

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/connorp2311/zfsTools/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	docDir string
)

// documentationCmd represents the documentation command
var documentationCmd = &cobra.Command{
	Use:   "documentation",
	Short: "Generate documentation for the tool",
	Long:  `Generate documentation for the tool`,
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := utils.NewLogger(logFile, "DOC")
		if err != nil {
			fmt.Printf("Error creating logger: %s\n", err)
			os.Exit(1)
		}
		defer logger.Close()

		docDir, err := filepath.Abs(docDir)
		if err != nil {
			logger.Log("Error converting relative path to absolute path: " + err.Error())
		}

		err = os.MkdirAll(docDir, 0755)
		if err != nil {
			logger.Log("Error creating directory: " + err.Error())
			os.Exit(1)
		}

		err = doc.GenMarkdownTree(rootCmd, docDir)
		if err != nil {
			panic(err)
		} else {
			logger.Log("Documentation generated in " + docDir)
		}
	},
}

func init() {
	rootCmd.AddCommand(documentationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// documentationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// documentationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// add a directory flag to specify where to generate the documentation
	documentationCmd.Flags().StringVarP(&docDir, "doc-dir", "d", "", "Output directory for documentation")
	documentationCmd.MarkFlagRequired("doc-dir")
}
