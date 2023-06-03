package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/gorilla/mux"
)

var OPENAI_API_KEY string
var tempDir, _ = ioutil.TempDir("", "example")

func main() {
	err := readAPIKey()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println("Temporary directory:", tempDir)

	// Start the RESTful server
	go startServer()

	// Run the systray module
	systray.Run(onReady, onExit)
}

func coldStartPrompt(framework, useCase string) string {
	return fmt.Sprintf("Use Typescript, react and vite.js with the %s UI framework to create a %s.\nAssume that Node.js, npm, vite.js are all downloaded already.\nONLY output the list of filepaths required to create %s, with \"\\n\" as a delimiter. Remember, only output the list.", framework, useCase, useCase)
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
	r.HandleFunc("/api/userPrompt", handleUserPrompt).Methods("POST")
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

	coldStartPromptRequest := coldStartPrompt(requestData.Framework, requestData.UseCase)
	filePathsResponse := singleQueryLLM(coldStartPromptRequest)
	// If response is empty, return an error
	if filePathsResponse == "" {
		http.Error(w, "Error querying LLM, did you export the API key to $OPENAI_API_KEY?", http.StatusInternalServerError)
		return
	}

	filePaths := strings.Split(filePathsResponse, "\n")
	createFiles(filePaths)
	codeGenPrompt := "Write the code for each file outputted in the chat history in sequential order.\nEach file's code should be formatted as \"###FILENAME:\nfilename\n###CODE:\ncode\"\nOnly output filepath and raw code for each file followed by the line breaks and delimiter."
	llmHistory := []string{coldStartPromptRequest, filePathsResponse, codeGenPrompt}
	combinedCodeResponse := multiQueryLLM(llmHistory)
	codeResponse := combinedCodeResponse[len(combinedCodeResponse)-1]
	writeFiles(codeResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set the HTTP status code to 200
	json.NewEncoder(w).Encode(struct {
		Message  string `json:"message"`
		Response string `json:"response"`
	}{
		Message:  "Prompt received successfully",
		Response: codeResponse,
	})
}

func writeFiles(inputCode string) error {
	entries := strings.Split(inputCode, "###FILENAME:")

	// Skip the first entry since it's empty
	entries = entries[1:]

	for _, entry := range entries {
		parts := strings.SplitN(entry, "###CODE:\n", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid input format: %s", entry)
		}

		filename := strings.TrimSpace(parts[0])
		code := strings.TrimSpace(parts[1])
		// Write the file
		fullPath := filepath.Join(tempDir, filename)
		err := ioutil.WriteFile(fullPath, []byte(code), 0644)
		if err != nil {
			return fmt.Errorf("error writing file %s: %v", fullPath, err)
		}
	}

	return nil
}

func handleUserPrompt(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Prompt string `json:"prompt"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Query ChatGPT
	response := singleQueryLLM(requestData.Prompt)
	// If response is empty, return an error
	if response == "" {
		http.Error(w, "Error querying ChatGPT, did you export the API key to $OPENAI_API_KEY?", http.StatusInternalServerError)
		return
	}

	// Send a response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set the HTTP status code to 200
	json.NewEncoder(w).Encode(struct {
		Message  string `json:"message"`
		Response string `json:"response"`
	}{
		Message:  "Prompt received successfully",
		Response: response,
	})
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
	fmt.Println("Finished!")
}
