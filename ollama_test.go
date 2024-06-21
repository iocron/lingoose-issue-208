package llm

import (
	"context"
	"github.com/henomis/lingoose/document"
	"github.com/henomis/lingoose/types"
	"regexp"
	"testing"
)

func TestOllamaAssistantWithRag(t *testing.T) {
	rag := OllamaRagNew("index.json")

	err := rag.AddDocuments(
		context.Background(),
		document.Document{
			Content: "this is some text about hello world",
			Metadata: types.Meta{
				"author": "Wikipedia",
			},
		},
		document.Document{
			Content: "this is a little side story in paris about a little mermaid",
			Metadata: types.Meta{
				"author": "Wikipedia",
			},
		},
	)
	if err != nil {
		t.Error("Not able to add document to RAG..\n", err)
	}

	ragQueryData, _ := rag.Retrieve(context.Background(), "where is the little mermaid")
	t.Log("ragQueryData: ", ragQueryData)

	llmAssistant, _ := OllamaAssistantNew(OllamaAssistantOptions{
		Model:         "llama3",
		UserMessage:   "It's about a story of a little mermaid.",
		Rag:           rag,
		SystemMessage: "You are a helpful AI Assistant. Count the words and always return the answer as a number without verbose information. Answer 0 if you are unsure.",
		// Temperature:   0,
	})

	answer, _ := OllamaAssistantLastMessage(llmAssistant)
	answerIsNumber, _ := regexp.MatchString(`^\d+$`, answer)

	t.Log(llmAssistant.Thread())

	if !answerIsNumber {
		t.Error("\nAnswer needs to be a number!")
	}
}
