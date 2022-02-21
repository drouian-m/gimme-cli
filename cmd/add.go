package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a module to the CDN",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cdnUrl, ok := os.LookupEnv("GIMME_URL")
		if !ok {
			fmt.Println("Missing GIMME_URL environment variable")
			return
		}
		token, ok := os.LookupEnv("GIMME_TOKEN")
		if !ok {
			fmt.Println("Missing GIMME_TOKEN environment variable")
			return
		}

		method := "POST"

		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		filePath, err := cmd.Flags().GetString("file")
		moduleName, err := cmd.Flags().GetString("name")
		version, err := cmd.Flags().GetString("version")
		if err != nil {
			fmt.Println("Error: File does not exists")
			return
		}

		fmt.Printf("Uploading module %s@%s\n", moduleName, version)

		file, err := os.Open(filePath)
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println("Fail to close file")
			}
		}(file)

		if err != nil {
			fmt.Println("Error", err)
			return
		}

		formFile,
			err := writer.CreateFormFile("file", filepath.Base(filePath))
		_, err = io.Copy(formFile, file)
		if err != nil {
			fmt.Println("Error", err)
			return
		}
		_ = writer.WriteField("name", moduleName)
		_ = writer.WriteField("version", version)
		err = writer.Close()
		if err != nil {
			fmt.Println("Error", err)
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest(method, fmt.Sprintf("%s%s", cdnUrl, "/packages"), payload)

		if err != nil {
			fmt.Println("Error", err)
			return
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		req.Header.Set("Content-Type", writer.FormDataContentType())
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error", err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println("Fail to close file")
			}
		}(res.Body)

		_, err = ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error", err)
			return
		}
		fmt.Printf("Module %s@%s has been successfully uploaded. You can retrieve it from %s/gimme/%s@%s/<file>\n", moduleName, version, cdnUrl, moduleName, version)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().String("name", "", "Module name")
	addCmd.PersistentFlags().String("version", "", "Module version (ex: 1.2.3)")
	addCmd.PersistentFlags().String("file", "", "Module file to upload (zip only)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
