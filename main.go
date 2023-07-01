package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-chatgpt/config"
	"io/ioutil"
	"net/http"
	"strings"
)

type Prompt struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
}

type GPTResponse struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type Ingredients struct {
	Items []string `json:"ingredients"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	Temperature float32 `json:"temperature"`
}

func getGptResponse(messages []Message) (*GPTResponse, error) {
	body := RequestBody{
		Model: "gpt-3.5-turbo",
		Messages: messages,
		Temperature: 0.7,
	}
	
	promptJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url := config.Config.ApiUrl + "/v1/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(promptJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+config.Config.ApiSecret)
	req.Header.Set("OpenAI-Organization", config.Config.ApiOrg)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))

	var response GPTResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func recipeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var ingredients Ingredients
	err := decoder.Decode(&ingredients)
	fmt.Print(&ingredients)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	messages := []Message{
		{
			Role: "user",
			Content: "NOTES:*日本語で料理名だけを3つ~5つ読点で区切って教えてください。料理の材料として" + strings.Join(ingredients.Items, ", ") + "があります。何を作れますか?",
		},
	}

	response, err := getGptResponse(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	recipe := response.Choices[0].Message.Content
	err = json.NewEncoder(w).Encode(recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/recipe", recipeHandler)
	http.ListenAndServe(":"+config.Config.Port, nil)
}