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
	return fmt.Sprintf("Given an existing codebase, use the %s UI framework to create a %s.\nWrite the complete code for any files that needs to be changed, which should look like: \"###FILENAME:\nfilename\n###CODE:\ncode\".\n If there is an API key involved, make a centralized .env file with all the keys needed, and read from that file in your new code. Use the UI framework whereever fit. Design the UI for a desktop webapp, and use the UI framework to make the components beautiful. package-lock.json is redacted due to its length. Here is the source code: %s", framework, useCase, encodedFiles)
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

func coldStartHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Framework string `json:"framework"`
		UseCase   string `json:"useCase"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//coldStartPromptRequest := coldStartPrompt(requestData.Framework, requestData.UseCase)
	//codeChanges := singleQueryLLM(coldStartPromptRequest)
	codeChanges := "###FILENAME:\n.env\n###CODE:\nREACT_APP_WEATHER_API_KEY=your_api_key_here\n\n###FILENAME:\nsrc/App.tsx\n###CODE:\nimport React, { useState } from 'react';\nimport { Layout, Input, Button, Card, Typography } from 'antd';\nimport './App.css';\n\nconst { Header, Content } = Layout;\nconst { Title } = Typography;\n\nconst App: React.FC = () => {\n  const [city, setCity] = useState('');\n  const [weatherData, setWeatherData] = useState<any>(null);\n\n  const fetchWeatherData = async () => {\n    const response = await fetch(\n      `https://api.openweathermap.org/data/2.5/weather?q=${city}&appid=${process.env.REACT_APP_WEATHER_API_KEY}&units=metric`\n    );\n    const data = await response.json();\n    setWeatherData(data);\n  };\n\n  return (\n    <Layout className=\"layout\">\n      <Header>\n        <Title level={2} style={{ color: 'white' }}>\n          Weather App\n        </Title>\n      </Header>\n      <Content className=\"content\">\n        <Input\n          placeholder=\"Enter city name\"\n          value={city}\n          onChange={(e) => setCity(e.target.value)}\n          onPressEnter={fetchWeatherData}\n        />\n        <Button type=\"primary\" onClick={fetchWeatherData}>\n          Search\n        </Button>\n        {weatherData && (\n          <Card title={weatherData.name} className=\"weather-card\">\n            <p>Temperature: {weatherData.main.temp}°C</p>\n            <p>Feels like: {weatherData.main.feels_like}°C</p>\n            <p>Humidity: {weatherData.main.humidity}%</p>\n            <p>Wind speed: {weatherData.wind.speed} m/s</p>\n          </Card>\n        )}\n      </Content>\n    </Layout>\n  );\n};\n\nexport default App;\n\n###FILENAME:\nsrc/App.css\n###CODE:\n.layout {\n  height: 100vh;\n}\n\n.content {\n  padding: 50px;\n  display: flex;\n  flex-direction: column;\n  align-items: center;\n  justify-content: center;\n}\n\n.weather-card {\n  margin-top: 20px;\n  width: 300px;\n}"
	// If response is empty, return an error
	if codeChanges == "" {
		http.Error(w, "Error querying LLM, did you export the API key to $OPENAI_API_KEY?", http.StatusInternalServerError)
		return
	}

	newChanges := decodeFilesFromResponse(codeChanges)

	for _, change := range newChanges {
		writeFiles(change)
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

func decodeFilesFromResponse(encodedFiles string) [][]string {
	// Define the regular expression pattern to match the filename and code portions
	filenamePattern := "###FILENAME:\n"
	codePattern := "###CODE:\n"

	fileCodeCombined := strings.Split(encodedFiles, filenamePattern)

	var result [][]string
	for _, fileCode := range fileCodeCombined {

		fileCodeSplit := strings.Split(fileCode, codePattern)
		if len(fileCodeSplit) != 2 {
			fmt.Println("Length of array is:", len(fileCodeSplit))
			fmt.Println(fileCodeSplit)
			continue
		} else {
			// Trim the trailing newline character and spaces
			fileCodeSplit[0] = strings.TrimSpace(fileCodeSplit[0])
			fileCodeSplit[1] = strings.TrimSpace(fileCodeSplit[1])
		}
		// Append the array to the result
		result = append(result, fileCodeSplit)
	}
	return result
}

func writeFiles(inputCode []string) error {
	filename := strings.TrimSpace(inputCode[0])
	code := strings.TrimSpace(inputCode[1])
	// Write the file
	fullPath := filepath.Join(tempDir, filename)
	err := ioutil.WriteFile(fullPath, []byte(code), 0644)
	if err != nil {
		return fmt.Errorf("error writing file %s: %v", fullPath, err)
	}

	return nil
}

func createFiles(filePaths []string) error {
	for _, filePath := range filePaths {
		// Create the directory structure
		dirPath := filepath.Join(tempDir, filepath.Dir(filePath))
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}

		// Create an empty file
		fullPath := filepath.Join(tempDir, filePath)
		err = ioutil.WriteFile(fullPath, []byte{}, 0644)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
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
	fmt.Println("Response choices amount:", len(response.Choices))

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
