package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting app.")

	// Start the RESTful server
	go startServer()

	// Run the systray
	systray.Run(onReady, onExit)
}

func startServer() {
	r := mux.NewRouter()
	r.HandleFunc("/api/userPrompt", handleUserPrompt).Methods("POST")

	fmt.Println("Starting RESTful server on port 8080...")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
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
	response := queryChatGPT(requestData.Prompt)
	fmt.Println("ChatGPT response:", response)

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

func queryChatGPT(prompt string) string {
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
		Temperature: 0.7,
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
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

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

func onExit() {
	fmt.Println("Finished!")
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
	fmt.Println("Started!")
}
