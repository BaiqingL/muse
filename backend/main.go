package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/gorilla/mux"
)

var OPENAI_API_KEY string
var tempDir, _ = ioutil.TempDir("", "example")
var encodedFiles = encodeFilesToPrompt("base")

type File struct {
	Filename string
	Code     string
}

func main() {
	// First read the API key
	err := readAPIKey()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Print temp folder for debugging
	fmt.Println("Temporary directory:", tempDir)

	// Copy the base files to the temp folder
	copyFiles("base", tempDir)

	// Start the RESTful server
	go startServer()

	// Handle graceful control-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			onExit()
			systray.Quit()
			os.Exit(0)
		}
	}()

	// Run the systray module
	systray.Run(onReady, onExit)
}

func copyFiles(src, dest string) error {
	// Retrieve information about the source directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory if it doesn't exist
	err = os.MkdirAll(dest, srcInfo.Mode())
	if err != nil {
		return err
	}

	// Retrieve a list of files and directories in the source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Iterate through each entry in the source directory
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// If the entry is a directory, recursively copy its contents
			err = copyFiles(srcPath, destPath)
			if err != nil {
				return err
			}
		} else {
			// If the entry is a file, copy it to the destination directory
			err = copyFile(srcPath, destPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func coldStartPrompt(framework, useCase string) string {
	return fmt.Sprintf(`Given an existing codebase, use the %s UI framework to create a %s.
	First, output the list of all the packages that needs to be installed, such as the frameworks, supporting packages, routers, anything. It should look like:
	###PACKAGES:
	package1
	package2

	Only output top level packages, not subpackages.

	Then, write the complete code for any files that needs to be changed, which should look like:

	###FILENAME:
	filename
	###CODE:
	code

	If there is an API key involved, make a centralized .env file with all the keys needed, and read from that file in your new code.
	Use the UI framework whereever fit. Design the UI for a desktop webapp, and use the UI framework to make the components beautiful.
	package-lock.json is redacted due to its length. Remember, dependency must be listed out first. ONLY output the packages, and then the code.
	When writing code, make sure the code actually exist, do not hallucinate code if you aren't sure. Make sure the webapp can run with no errors.
	In the env file, make it clear what API the key is for.
	Here is the source code: %s`,
		framework, useCase, encodedFiles)
}

func readAPIKey() error {
	filePath := "apiKey.env"

	// Read the content of the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %v", filePath, err)
	}

	// Extract the API key from the content
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "OPENAI_API_KEY=") {
			OPENAI_API_KEY = strings.TrimPrefix(line, "OPENAI_API_KEY=")
			return nil
		}
	}

	return fmt.Errorf("API key not found in %s", filePath)
}

func startServer() {
	r := mux.NewRouter()
	r.HandleFunc("/api/coldStart", coldStartHandler).Methods("POST")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

func parsePackages(input string) []string {
	startIndex := strings.Index(input, "###PACKAGES:\n")
	endIndex := strings.Index(input, "###FILENAME:\n")
	if startIndex == -1 || endIndex == -1 {
		return nil
	}

	packages := strings.Split(input[startIndex+len("###PACKAGES:\n"):endIndex], "\n")
	return packages
}

func parseCode(input string) []File {
	files := []File{}

	for i := 0; i < len(input); {
		filenameIndex := strings.Index(input[i:], "###FILENAME:\n")
		if filenameIndex == -1 {
			break
		}
		filenameStart := i + filenameIndex + len("###FILENAME:\n")
		filenameEnd := strings.Index(input[filenameStart:], "\n")
		if filenameEnd == -1 {
			break
		}
		filename := input[filenameStart : filenameStart+filenameEnd]

		codeIndex := strings.Index(input[filenameStart+filenameEnd:], "###CODE:\n")
		if codeIndex == -1 {
			break
		}
		codeStart := filenameStart + filenameEnd + codeIndex + len("###CODE:\n")
		codeEnd := strings.Index(input[codeStart:], "###FILENAME:\n")
		if codeEnd == -1 {
			codeEnd = len(input[codeStart:])
		}

		code := input[codeStart : codeStart+codeEnd]
		files = append(files, File{Filename: filename, Code: code})

		i = codeStart + codeEnd
	}

	return files
}

func installDependencies(packages []string) error {
	for _, pkg := range packages {
		cmd := exec.Command("npm", "install", pkg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = tempDir

		fmt.Printf("Installing package: %s\n", pkg)
		err := cmd.Run()
		if err != nil {
			return err
		}

		fmt.Printf("Package %s installed successfully.\n", pkg)
	}

	return nil
}

func coldStartHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received cold start request")
	var requestData struct {
		Framework string `json:"framework"`
		UseCase   string `json:"useCase"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	coldStartCodeRequest := coldStartPrompt(requestData.Framework, requestData.UseCase)
	codeChanges := singleQueryLLM(coldStartCodeRequest)
	// If response is empty, return an error
	if codeChanges == "" {
		http.Error(w, "Error querying OpenAI, did you export the API key to $OPENAI_API_KEY?", http.StatusInternalServerError)
		return
	}

	packages := parsePackages(codeChanges)

	code := parseCode(codeChanges)
	installDependencies(packages)

	for _, file := range code {
		fmt.Printf("Filename: %s\n", file.Filename)
		writeFile(file)
	}

}

func encodeFilesToPrompt(filePath string) string {
	// Retrieve a list of files in the directory
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	resultString := ""

	// Iterate through each file in the directory
	for _, file := range files {
		// If filename is .gitignore or package-lock.json or .eslintrc.cjs, skip it
		if file.Name() == ".gitignore" || file.Name() == "package-lock.json" || file.Name() == ".eslintrc.cjs" {
			continue
		}
		// Check if the file is a regular file (not a directory)
		if file.Mode().IsRegular() {
			// Read the contents of the file
			filePath := filepath.Join(filePath, file.Name())
			contents, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			// Prepare the formatted string
			code := string(contents)
			formattedString := fmt.Sprintf("###FILENAME:\n%s\n###CODE:\n%s", filePath, code)

			// Print the formatted string
			resultString += formattedString + "\n"
		}
	}
	return resultString
}

func writeFile(inputCode File) error {
	// Check if tempDir + inputCode.Filename's dir exists
	dir := filepath.Dir(inputCode.Filename)
	fullDir := filepath.Join(tempDir, dir)
	// Check if the directory exists
	if _, err := os.Stat(fullDir); os.IsNotExist(err) {
		// Create the directory
		fmt.Println("Creating directory:", fullDir)
		err = os.MkdirAll(fullDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %v", fullDir, err)
		}
	} else {
		fmt.Println("Directory already exists:", fullDir)
	}

	// Write the file
	fullPath := filepath.Join(tempDir, inputCode.Filename)
	fmt.Println("Writing file:", fullPath)
	err := ioutil.WriteFile(fullPath, []byte(inputCode.Code), 0644)
	if err != nil {
		return fmt.Errorf("error writing file %s: %v", fullPath, err)
	}

	return nil
}

func singleQueryLLM(prompt string) string {
	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	type Request struct {
		Model       string    `json:"model"`
		Messages    []Message `json:"messages"`
		Temperature float64   `json:"temperature"`
	}

	requestData := Request{
		Model: "gpt-4",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error encoding request body:", err)
		return ""
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+OPENAI_API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	fmt.Println("ChatGPT response:", string(responseBody))

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		fmt.Println("Error parsing ChatGPT response:", err)
		return ""
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content
	}

	return ""
}

func multiQueryLLM(prompts []string) []string {
	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	type Request struct {
		Model       string    `json:"model"`
		Messages    []Message `json:"messages"`
		Temperature float64   `json:"temperature"`
	}

	requestData := Request{
		Model:       "gpt-4",
		Messages:    make([]Message, 0),
		Temperature: 0,
	}

	for i, prompt := range prompts {
		role := "user"
		if i%2 == 1 {
			role = "assistant"
		}

		requestData.Messages = append(requestData.Messages, Message{
			Role:    role,
			Content: prompt,
		})
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error encoding request body:", err)
		return nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+OPENAI_API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		fmt.Println("Error parsing ChatGPT response:", err)
		return nil
	}

	// Print response choices amount
	fmt.Println("Response body:", string(responseBody))

	results := make([]string, len(response.Choices))
	for i, choice := range response.Choices {
		results[i] = choice.Message.Content
	}

	return results
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()
}

func onExit() {
	err := os.RemoveAll(tempDir)
	if err != nil {
		fmt.Println("Error removing temporary directory:", err)
		os.Exit(1)
	}
	fmt.Println("Removed temporary directory")
}
