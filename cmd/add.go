/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
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
	Short: "Add package to the CDN",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		fmt.Printf("Uploading module %s@%s\n", moduleName, version)

		if err != nil {
			fmt.Println("Error: File does not exists")
		}
		file, errFile1 := os.Open(filePath)
		defer file.Close()
		if errFile1 != nil {
			fmt.Println("Error", errFile1)
			return
		}
		part1,
			errFile1 := writer.CreateFormFile("file", filepath.Base(filePath))
		_, errFile1 = io.Copy(part1, file)
		if errFile1 != nil {
			fmt.Println("Error", errFile1)
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
		defer res.Body.Close()

		_, err = ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
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
