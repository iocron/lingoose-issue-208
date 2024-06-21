# Lingoose Issue #208
Using lingoose with RAG produces 2 system messages in the Thread (+sometimes returns random output/answer, depending on the chosen RAG/DB/Model, but that might be some other issue).

Additionally the .withModel() method on the embedder does not seem to have any effect or returns empty.

Link to Issue: https://github.com/henomis/lingoose/issues/208

## To reproduce
```
go test -v ./...
```

This will output a thread with 2 system messages (one from lingoose/assistant/prompt.go and one from the user defined system message).

## Test Output
```
TestOllamaAssistantWithRag
    ollama_test.go:34: ragQueryData:  [this is a little side story in paris about a little mermaid]
    ollama_test.go:47: Thread:
        system:
        	Type: text
        	Text: You are a helpful AI Assistant. Count the words and always return the answer as a number without verbose information. Answer 0 if you are unsure.
        system:
        	Type: text
        	Text: You name is AI assistant, and you are a helpful and polite assistant . Your task is to assist humans with their questions.
        user:
        	Type: text
        	Text: Use the following pieces of retrieved context to answer the question.

        Question: It's about a story of a little mermaid.
        Context:
        this is a little side story in paris about a little mermaid


        assistant:
        	Type: text
        	Text:

    ollama_test.go:50:
        Answer needs to be a number!
```


