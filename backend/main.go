package main

import (
	"fmt"

	"github.com/BaiqingL/conviction-commit-team-1"
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
	response := query.QueryChatGPT(requestData.Prompt)
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
