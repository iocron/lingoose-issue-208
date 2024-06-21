package llm

import (
	"context"
	"errors"
	"fmt"
	"github.com/henomis/lingoose/assistant"
	ollamaembedder "github.com/henomis/lingoose/embedder/ollama"
	"github.com/henomis/lingoose/index"
	"github.com/henomis/lingoose/index/vectordb/jsondb"
	"github.com/henomis/lingoose/llm/ollama"
	"github.com/henomis/lingoose/rag"
	"github.com/henomis/lingoose/thread"
)

type OllamaAssistantOptions struct {
	Endpoint      string
	UserMessage   string
	Model         string
	Rag           assistant.RAG
	SystemMessage string
	Temperature   float64
}

func OllamaAssistantLastMessage(assistant *assistant.Assistant) (string, error) {
	if len(assistant.Thread().LastMessage().Contents) <= 0 {
		return "", errors.New("GetOllamaAssistantLastMessage(): There is no content.")
	}
	return assistant.Thread().LastMessage().Contents[0].AsString(), nil
}

func OllamaRagNew(dbPath string) *rag.RAG {
	// BUG: There will be 2 system messages at the end
	return rag.New(
		index.New(
			jsondb.New().WithPersist(dbPath),
			ollamaembedder.New().WithModel("nomic-embed-text"), // BUG: Using .WithModel seems not to do anything (sometimes returns empty)
			// openaiembedder.New(openaiembedder.AdaEmbeddingV2),
		),
	).WithChunkSize(1024).WithChunkOverlap(0)
}

func OllamaAssistantNew(options OllamaAssistantOptions) (*assistant.Assistant, error) {
	// LLM Init
	llm := ollama.New().WithModel(options.Model)

	// LLM Endpoint
	if options.Endpoint != "" {
		llm.WithEndpoint(options.Endpoint)
	}

	// LLM Temperature
	if options.Temperature >= 0 {
		llm.WithTemperature(options.Temperature)
	}

	// LLM Assistant / LLM Thread
	llmAssistant := assistant.New(llm) // TODO: Adding a RAG Query/Text as well?!
	llmThread := thread.New()

	// LLM Assistant - RAG
	if options.Rag != nil {
		llmAssistant.WithRAG(options.Rag)
	}

	// LLM Thread - Add System Message
	if options.SystemMessage != "" {
		llmThread.AddMessage(thread.NewSystemMessage().AddContent(thread.NewTextContent(options.SystemMessage)))
	}

	// LLM Thread - Add User Message
	if options.UserMessage != "" {
		llmThread.AddMessage(thread.NewUserMessage().AddContent(thread.NewTextContent(options.UserMessage)))
	}

	// RUN LLM Assistant / LLM Thread
	err := llmAssistant.RunWithThread(context.Background(), llmThread)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// fmt.Println("llmThread: ", llmThread)
	// fmt.Println("llmAssistant: ", llmAssistant)

	return llmAssistant, nil
}
