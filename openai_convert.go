package openingrouter

import (
	"time"

	"github.com/coalaura/openingrouter/internal/openai"
)

// openaiModelToModel converts a single OpenAI model into an openingrouter Model.
// OpenAI exposes far fewer fields, so only ID, Created and the shutdown date have
// a target; the rest are left at their zero values.
func openaiModelToModel(m *openai.Model) *Model {
	if m == nil {
		return nil
	}

	result := &Model{
		ID:            m.ID,
		CanonicalSlug: m.ID,
		Name:          m.ID,
		Created:       m.Created,
	}

	if m.ShutdownDate != nil {
		if ft, ok := parseFlexibleTime(*m.ShutdownDate); ok {
			result.ExpirationDate = &ft
		}
	}

	return result
}

// openaiModelsListToModels converts an OpenAI models list into openingrouter models.
func openaiModelsListToModels(list *openai.ModelsList) []Model {
	if list == nil {
		return nil
	}

	result := make([]Model, 0, len(list.Data))
	for i := range list.Data {
		result = append(result, *openaiModelToModel(&list.Data[i]))
	}

	return result
}

// parseFlexibleTime parses a timestamp in any of the supported formats.
func parseFlexibleTime(value string) (FlexibleTime, bool) {
	for _, format := range timeFormats {
		if t, err := time.Parse(format, value); err == nil {
			return FlexibleTime{Time: t}, true
		}
	}

	return FlexibleTime{}, false
}

// chatCompletionRequestToOpenAI converts an openingrouter chat completion request
// into the OpenAI wire format. OpenRouter-only fields (provider routing, plugins,
// server tools, reasoning details, video parts, ...) have no OpenAI equivalent and
// are dropped.
func chatCompletionRequestToOpenAI(req *ChatCompletionRequest) *openai.ChatCompletionRequest {
	if req == nil {
		return nil
	}

	result := &openai.ChatCompletionRequest{
		Model:               req.Model,
		Messages:            chatMessagesToOpenAI(req.Messages),
		FrequencyPenalty:    req.FrequencyPenalty,
		LogitBias:           req.LogitBias,
		Logprobs:            req.Logprobs,
		MaxCompletionTokens: req.MaxCompletionTokens,
		MaxTokens:           req.MaxTokens,
		Metadata:            req.Metadata,
		Modalities:          modalitiesToOpenAI(req.Modalities),
		ParallelToolCalls:   req.ParallelToolCalls,
		Prediction:          chatPredictionToOpenAI(req.Prediction),
		PresencePenalty:     req.PresencePenalty,
		PromptCacheKey:      req.PromptCacheKey,
		PromptCacheOptions:  promptCacheOptionsToOpenAI(req.PromptCacheOptions),
		ResponseFormat:      chatResponseFormatToOpenAI(req.ResponseFormat),
		Seed:                req.Seed,
		ServiceTier:         string(req.ServiceTier),
		Stop:                req.Stop,
		Stream:              req.Stream,
		StreamOptions:       streamOptionsToOpenAI(req.StreamOptions),
		Temperature:         req.Temperature,
		ToolChoice:          chatToolChoiceToOpenAI(req.ToolChoice),
		Tools:               chatToolsToOpenAI(req.Tools),
		TopLogprobs:         req.TopLogprobs,
		TopP:                req.TopP,
		User:                req.User,
	}

	// Reasoning effort: the shorthand field wins, otherwise the nested config.
	effort := req.ReasoningEffort
	if effort == "" && req.Reasoning != nil {
		effort = req.Reasoning.Effort
	}

	if effort != "" {
		result.ReasoningEffort = string(effort)
	}

	return result
}

// modalitiesToOpenAI keeps only the modalities OpenAI understands (text, audio).
func modalitiesToOpenAI(modalities []OutputModality) []string {
	if len(modalities) == 0 {
		return nil
	}

	var result []string

	for _, modality := range modalities {
		switch modality {
		case OutputModalityText, OutputModalityAudio:
			result = append(result, string(modality))
		}
	}

	return result
}

// chatPredictionToOpenAI converts a predicted output, dropping OpenRouter-only bits.
func chatPredictionToOpenAI(p *ChatPrediction) *openai.ChatPrediction {
	if p == nil {
		return nil
	}

	return &openai.ChatPrediction{
		Type:    string(p.Type),
		Content: chatContentToOpenAI(&p.Content),
	}
}

// promptCacheOptionsToOpenAI converts the request level prompt cache controls.
func promptCacheOptionsToOpenAI(opts *ChatPromptCacheOptions) *openai.ChatPromptCacheOptions {
	if opts == nil {
		return nil
	}

	return &openai.ChatPromptCacheOptions{
		Mode: string(opts.Mode),
		TTL:  opts.TTL,
	}
}

// streamOptionsToOpenAI converts the streaming options.
func streamOptionsToOpenAI(opts *ChatStreamOptions) *openai.ChatStreamOptions {
	if opts == nil {
		return nil
	}

	return &openai.ChatStreamOptions{
		IncludeUsage: opts.IncludeUsage,
	}
}

// chatResponseFormatToOpenAI converts the response format, dropping the
// OpenRouter-only grammar and python formats that OpenAI does not accept.
func chatResponseFormatToOpenAI(rf *ChatResponseFormat) *openai.ChatResponseFormat {
	if rf == nil {
		return nil
	}

	switch rf.Type {
	case ChatResponseFormatTypeText, ChatResponseFormatTypeJSONObject, ChatResponseFormatTypeJSONSchema:
	default:
		return nil
	}

	result := &openai.ChatResponseFormat{
		Type: string(rf.Type),
	}

	if rf.Type == ChatResponseFormatTypeJSONSchema && rf.JSONSchema != nil {
		result.JSONSchema = &openai.ChatJSONSchema{
			Name:        rf.JSONSchema.Name,
			Description: rf.JSONSchema.Description,
			Schema:      rf.JSONSchema.Schema,
			Strict:      rf.JSONSchema.Strict,
		}
	}

	return result
}

// chatToolChoiceToOpenAI converts the tool choice configuration.
func chatToolChoiceToOpenAI(tc *ChatToolChoice) *openai.ChatToolChoice {
	if tc == nil {
		return nil
	}

	result := &openai.ChatToolChoice{
		Mode: string(tc.Mode),
		Type: string(tc.Type),
	}

	if tc.Function != nil {
		result.Function = &openai.ChatToolChoiceFunc{
			Name: tc.Function.Name,
		}
	}

	return result
}

// chatToolsToOpenAI keeps only regular function tools; OpenRouter server tools
// (bash, web search, subagents, ...) have no OpenAI equivalent and are dropped.
func chatToolsToOpenAI(tools []ChatTool) []openai.ChatFunctionTool {
	if len(tools) == 0 {
		return nil
	}

	var result []openai.ChatFunctionTool

	for _, tool := range tools {
		functionTool, ok := tool.(ChatFunctionTool)
		if !ok {
			continue
		}

		result = append(result, openai.ChatFunctionTool{
			Type: string(functionTool.Type),
			Function: openai.ChatFunction{
				Name:        functionTool.Function.Name,
				Description: functionTool.Function.Description,
				Parameters:  functionTool.Function.Parameters,
				Strict:      functionTool.Function.Strict,
			},
		})
	}

	return result
}

// chatMessagesToOpenAI converts a list of messages.
func chatMessagesToOpenAI(messages []ChatMessage) []openai.ChatMessage {
	if len(messages) == 0 {
		return nil
	}

	result := make([]openai.ChatMessage, 0, len(messages))
	for i := range messages {
		result = append(result, chatMessageToOpenAI(&messages[i]))
	}

	return result
}

// chatMessageToOpenAI converts a single message, dropping OpenRouter-only fields
// (reasoning details, generated images, the serving model).
func chatMessageToOpenAI(m *ChatMessage) openai.ChatMessage {
	result := openai.ChatMessage{
		Role:       openai.ChatRole(m.Role),
		Name:       m.Name,
		Refusal:    m.Refusal,
		ToolCallID: m.ToolCallID,
	}

	if m.Content.Text != "" || len(m.Content.Parts) > 0 {
		content := chatContentToOpenAI(&m.Content)
		result.Content = &content
	}

	if m.Audio != nil {
		result.Audio = &openai.ChatAudio{
			ID: m.Audio.ID,
		}
	}

	if len(m.ToolCalls) > 0 {
		result.ToolCalls = chatToolCallsToOpenAI(m.ToolCalls)
	}

	return result
}

// chatContentToOpenAI converts message content, keeping the text or the supported
// content parts.
func chatContentToOpenAI(c *ChatContent) openai.ChatContent {
	if c == nil {
		return openai.ChatContent{}
	}

	if len(c.Parts) == 0 {
		return openai.ChatContent{
			Text: c.Text,
		}
	}

	result := openai.ChatContent{
		Parts: make([]openai.ChatContentPart, 0, len(c.Parts)),
	}

	for i := range c.Parts {
		part := chatContentPartToOpenAI(&c.Parts[i])
		if part != nil {
			result.Parts = append(result.Parts, *part)
		}
	}

	return result
}

// chatContentPartToOpenAI converts a single content part. video_url, input_video
// and cache controls have no OpenAI equivalent and are dropped.
func chatContentPartToOpenAI(p *ChatContentPart) *openai.ChatContentPart {
	if p == nil {
		return nil
	}

	switch p.Type {
	case ChatContentPartTypeText:
		return &openai.ChatContentPart{
			Type: openai.ChatContentPartTypeText,
			Text: p.Text,
		}
	case ChatContentPartTypeImageURL:
		if p.ImageURL == nil {
			return nil
		}

		return &openai.ChatContentPart{
			Type: openai.ChatContentPartTypeImageURL,
			ImageURL: &openai.ChatContentImageURL{
				URL:    p.ImageURL.URL,
				Detail: string(p.ImageURL.Detail),
			},
		}
	case ChatContentPartTypeInputAudio:
		if p.InputAudio == nil {
			return nil
		}

		return &openai.ChatContentPart{
			Type: openai.ChatContentPartTypeInputAudio,
			InputAudio: &openai.ChatContentInputAudio{
				Data:   p.InputAudio.Data,
				Format: p.InputAudio.Format,
			},
		}
	case ChatContentPartTypeFile:
		if p.File == nil {
			return nil
		}

		return &openai.ChatContentPart{
			Type: openai.ChatContentPartTypeFile,
			File: &openai.ChatContentFile{
				FileData: p.File.FileData,
				FileID:   p.File.FileID,
				Filename: p.File.Filename,
			},
		}
	}

	return nil
}

// chatToolCallsToOpenAI converts the assistant tool calls of a message.
func chatToolCallsToOpenAI(calls []ChatToolCall) []openai.ChatToolCall {
	if len(calls) == 0 {
		return nil
	}

	result := make([]openai.ChatToolCall, 0, len(calls))
	for i := range calls {
		result = append(result, openai.ChatToolCall{
			ID:   calls[i].ID,
			Type: string(calls[i].Type),
			Function: openai.ChatToolCallFunction{
				Name:      calls[i].Function.Name,
				Arguments: calls[i].Function.Arguments,
			},
		})
	}

	return result
}

// completionRequestToOpenAI converts a completions request into the OpenAI wire
// format. The two shapes are intentionally identical, so this is a field copy.
func completionRequestToOpenAI(req *CompletionRequest) *openai.CompletionRequest {
	if req == nil {
		return nil
	}

	return &openai.CompletionRequest{
		Model:            req.Model,
		Prompt:           completionInputToOpenAI(req.Prompt),
		BestOf:           req.BestOf,
		Echo:             req.Echo,
		FrequencyPenalty: req.FrequencyPenalty,
		LogitBias:        req.LogitBias,
		Logprobs:         req.Logprobs,
		MaxTokens:        req.MaxTokens,
		N:                req.N,
		PresencePenalty:  req.PresencePenalty,
		Seed:             req.Seed,
		Stop:             req.Stop,
		Stream:           req.Stream,
		StreamOptions:    streamOptionsToOpenAI(req.StreamOptions),
		Suffix:           req.Suffix,
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		User:             req.User,
	}
}

// completionInputToOpenAI converts the prompt input.
func completionInputToOpenAI(input CompletionInput) openai.CompletionInput {
	return openai.CompletionInput{
		Text:        input.Text,
		Texts:       input.Texts,
		Tokens:      input.Tokens,
		TokenArrays: input.TokenArrays,
	}
}

// embeddingRequestToOpenAI converts an embeddings request. Multimodal inputs have
// no OpenAI equivalent and are dropped.
func embeddingRequestToOpenAI(req *EmbeddingRequest) *openai.EmbeddingRequest {
	if req == nil {
		return nil
	}

	return &openai.EmbeddingRequest{
		Model:          req.Model,
		Input:          embeddingInputToOpenAI(&req.Input),
		Dimensions:     req.Dimensions,
		EncodingFormat: openai.EmbeddingEncodingFormat(req.EncodingFormat),
		User:           req.User,
	}
}

// embeddingInputToOpenAI converts the embedding input, keeping the supported
// encodings.
func embeddingInputToOpenAI(input *EmbeddingInput) openai.EmbeddingInput {
	if input == nil {
		return openai.EmbeddingInput{}
	}

	return openai.EmbeddingInput{
		Text:        input.Text,
		Texts:       input.Texts,
		Tokens:      input.Tokens,
		TokenArrays: input.TokenArrays,
	}
}

// openaiChatCompletionResponseToResponse converts an OpenAI chat completion
// response into the openingrouter type.
func openaiChatCompletionResponseToResponse(r *openai.ChatCompletionResponse) *ChatCompletionResponse {
	if r == nil {
		return nil
	}

	result := &ChatCompletionResponse{
		ID:                r.ID,
		Object:            ChatObject(r.Object),
		Created:           r.Created,
		Model:             r.Model,
		Choices:           openaiChatChoicesToChatChoices(r.Choices),
		SystemFingerprint: r.SystemFingerprint,
		Usage:             openaiChatUsageToChatUsage(r.Usage),
	}

	if r.ServiceTier != "" {
		tier := r.ServiceTier
		result.ServiceTier = &tier
	}

	return result
}

// openaiChatChoicesToChatChoices converts the response choices.
func openaiChatChoicesToChatChoices(choices []openai.ChatChoice) []ChatChoice {
	if len(choices) == 0 {
		return nil
	}

	result := make([]ChatChoice, 0, len(choices))
	for i := range choices {
		result = append(result, ChatChoice{
			Index:        choices[i].Index,
			FinishReason: ChatFinishReason(choices[i].FinishReason),
			Message:      openaiChatOutMessageToChatMessage(&choices[i].Message),
			Logprobs:     openaiChatLogprobsToChatLogprobs(choices[i].Logprobs),
		})
	}

	return result
}

// openaiChatOutMessageToChatMessage converts an assistant message.
func openaiChatOutMessageToChatMessage(m *openai.ChatOutMessage) ChatMessage {
	if m == nil {
		return ChatMessage{}
	}

	result := ChatMessage{
		Role:    ChatRole(m.Role),
		Refusal: m.Refusal,
	}

	if m.Content != nil {
		result.Content = ChatContent{
			Text: *m.Content,
		}
	}

	if m.Audio != nil {
		result.Audio = &ChatAudioOutput{
			ID:         m.Audio.ID,
			Data:       m.Audio.Data,
			Transcript: m.Audio.Transcript,
			ExpiresAt:  m.Audio.ExpiresAt,
		}
	}

	if len(m.ToolCalls) > 0 {
		result.ToolCalls = openaiChatToolCallsToChatToolCalls(m.ToolCalls)
	}

	return result
}

// openaiChatToolCallsToChatToolCalls converts the assistant tool calls.
func openaiChatToolCallsToChatToolCalls(calls []openai.ChatToolCall) []ChatToolCall {
	if len(calls) == 0 {
		return nil
	}

	result := make([]ChatToolCall, 0, len(calls))
	for i := range calls {
		result = append(result, ChatToolCall{
			ID:   calls[i].ID,
			Type: ChatToolType(calls[i].Type),
			Function: ChatToolCallFunction{
				Name:      calls[i].Function.Name,
				Arguments: calls[i].Function.Arguments,
			},
		})
	}

	return result
}

// openaiChatLogprobsToChatLogprobs converts the log probabilities.
func openaiChatLogprobsToChatLogprobs(lp *openai.ChatLogprobs) *ChatLogprobs {
	if lp == nil {
		return nil
	}

	return &ChatLogprobs{
		Content: openaiChatTokenLogprobsToChatTokenLogprobs(lp.Content),
		Refusal: openaiChatTokenLogprobsToChatTokenLogprobs(lp.Refusal),
	}
}

// openaiChatTokenLogprobsToChatTokenLogprobs converts a list of token logprobs.
func openaiChatTokenLogprobsToChatTokenLogprobs(tokens []openai.ChatTokenLogprob) []ChatTokenLogprob {
	if len(tokens) == 0 {
		return nil
	}

	result := make([]ChatTokenLogprob, 0, len(tokens))
	for i := range tokens {
		result = append(result, ChatTokenLogprob{
			Token:       tokens[i].Token,
			Logprob:     tokens[i].Logprob,
			Bytes:       tokens[i].Bytes,
			TopLogprobs: openaiChatTopLogprobsToChatTopLogprobs(tokens[i].TopLogprobs),
		})
	}

	return result
}

// openaiChatTopLogprobsToChatTopLogprobs converts the top logprob alternatives.
func openaiChatTopLogprobsToChatTopLogprobs(tokens []openai.ChatTopLogprob) []ChatTopLogprob {
	if len(tokens) == 0 {
		return nil
	}

	result := make([]ChatTopLogprob, 0, len(tokens))
	for i := range tokens {
		result = append(result, ChatTopLogprob{
			Token:   tokens[i].Token,
			Logprob: tokens[i].Logprob,
			Bytes:   tokens[i].Bytes,
		})
	}

	return result
}

// openaiChatUsageToChatUsage converts the token usage. OpenAI does not report a
// cost, so Cost and IsBYOK stay at their zero values.
func openaiChatUsageToChatUsage(u *openai.ChatUsage) *ChatUsage {
	if u == nil {
		return nil
	}

	result := &ChatUsage{
		PromptTokens:     u.PromptTokens,
		CompletionTokens: u.CompletionTokens,
		TotalTokens:      u.TotalTokens,
	}

	if u.PromptTokensDetails != nil {
		result.PromptTokensDetails = &PromptTokensDetails{
			CachedTokens:     u.PromptTokensDetails.CachedTokens,
			CacheWriteTokens: u.PromptTokensDetails.CacheWriteTokens,
			AudioTokens:      u.PromptTokensDetails.AudioTokens,
		}
	}

	if u.CompletionTokensDetails != nil {
		result.CompletionTokensDetails = &ChatCompletionTokensDetails{
			ReasoningTokens:          u.CompletionTokensDetails.ReasoningTokens,
			AudioTokens:              u.CompletionTokensDetails.AudioTokens,
			AcceptedPredictionTokens: u.CompletionTokensDetails.AcceptedPredictionTokens,
			RejectedPredictionTokens: u.CompletionTokensDetails.RejectedPredictionTokens,
		}
	}

	return result
}

// openaiChatCompletionChunkToChunk converts a single streamed chat chunk.
func openaiChatCompletionChunkToChunk(c *openai.ChatCompletionChunk) ChatStreamChunk {
	result := ChatStreamChunk{
		ID:      c.ID,
		Object:  ChatObject(c.Object),
		Created: c.Created,
		Model:   c.Model,
		Choices: openaiChatChunkChoicesToChatStreamChoices(c.Choices),
		Usage:   openaiChatUsageToChatUsage(c.Usage),
	}

	if c.SystemFingerprint != nil {
		result.SystemFingerprint = *c.SystemFingerprint
	}

	if c.ServiceTier != "" {
		tier := c.ServiceTier
		result.ServiceTier = &tier
	}

	return result
}

// openaiChatChunkChoicesToChatStreamChoices converts the streamed choices.
func openaiChatChunkChoicesToChatStreamChoices(choices []openai.ChatChunkChoice) []ChatStreamChoice {
	if len(choices) == 0 {
		return nil
	}

	result := make([]ChatStreamChoice, 0, len(choices))
	for i := range choices {
		choice := ChatStreamChoice{
			Index:    choices[i].Index,
			Delta:    openaiChatChunkDeltaToChatStreamDelta(&choices[i].Delta),
			Logprobs: openaiChatLogprobsToChatLogprobs(choices[i].Logprobs),
		}

		if choices[i].FinishReason != nil {
			choice.FinishReason = ChatFinishReason(*choices[i].FinishReason)
		}

		result = append(result, choice)
	}

	return result
}

// openaiChatChunkDeltaToChatStreamDelta converts a streamed delta.
func openaiChatChunkDeltaToChatStreamDelta(d *openai.ChatChunkDelta) ChatStreamDelta {
	if d == nil {
		return ChatStreamDelta{}
	}

	result := ChatStreamDelta{
		Role:    ChatRole(d.Role),
		Content: d.Content,
		Refusal: d.Refusal,
	}

	if d.Audio != nil {
		result.Audio = &ChatAudioOutput{
			ID:         d.Audio.ID,
			Data:       d.Audio.Data,
			Transcript: d.Audio.Transcript,
			ExpiresAt:  d.Audio.ExpiresAt,
		}
	}

	if len(d.ToolCalls) > 0 {
		result.ToolCalls = openaiChatStreamToolCallsToChatStreamToolCalls(d.ToolCalls)
	}

	return result
}

// openaiChatStreamToolCallsToChatStreamToolCalls converts streamed tool call
// deltas.
func openaiChatStreamToolCallsToChatStreamToolCalls(calls []openai.ChatStreamToolCall) []ChatStreamToolCall {
	if len(calls) == 0 {
		return nil
	}

	result := make([]ChatStreamToolCall, 0, len(calls))
	for i := range calls {
		call := ChatStreamToolCall{
			Index: calls[i].Index,
			ID:    calls[i].ID,
			Type:  ChatToolType(calls[i].Type),
		}

		if calls[i].Function != nil {
			call.Function = &ChatStreamToolCallFunction{
				Name:      calls[i].Function.Name,
				Arguments: calls[i].Function.Arguments,
			}
		}

		result = append(result, call)
	}

	return result
}

// openaiCompletionResponseToResponse converts an OpenAI completions response into
// the openingrouter type.
func openaiCompletionResponseToResponse(r *openai.CompletionResponse) *CompletionResponse {
	if r == nil {
		return nil
	}

	return &CompletionResponse{
		ID:                r.ID,
		Object:            CompletionObject(r.Object),
		Created:           r.Created,
		Model:             r.Model,
		Choices:           openaiCompletionChoicesToCompletionChoices(r.Choices),
		SystemFingerprint: r.SystemFingerprint,
		Usage:             openaiCompletionUsageToCompletionUsage(r.Usage),
	}
}

// openaiCompletionChoicesToCompletionChoices converts the response choices.
func openaiCompletionChoicesToCompletionChoices(choices []openai.CompletionChoice) []CompletionChoice {
	if len(choices) == 0 {
		return nil
	}

	result := make([]CompletionChoice, 0, len(choices))
	for i := range choices {
		result = append(result, CompletionChoice{
			Index:        choices[i].Index,
			FinishReason: choices[i].FinishReason,
			Text:         choices[i].Text,
			Logprobs:     openaiCompletionLogprobsToCompletionLogprobs(choices[i].Logprobs),
		})
	}

	return result
}

// openaiCompletionLogprobsToCompletionLogprobs converts the log probabilities.
func openaiCompletionLogprobsToCompletionLogprobs(lp *openai.CompletionLogprobs) *CompletionLogprobs {
	if lp == nil {
		return nil
	}

	return &CompletionLogprobs{
		TextOffset:    lp.TextOffset,
		TokenLogprobs: lp.TokenLogprobs,
		Tokens:        lp.Tokens,
		TopLogprobs:   lp.TopLogprobs,
	}
}

// openaiCompletionUsageToCompletionUsage converts the token usage.
func openaiCompletionUsageToCompletionUsage(u *openai.CompletionUsage) *CompletionUsage {
	if u == nil {
		return nil
	}

	return &CompletionUsage{
		PromptTokens:     u.PromptTokens,
		CompletionTokens: u.CompletionTokens,
		TotalTokens:      u.TotalTokens,
	}
}

// openaiCompletionChunkToChunk converts a single streamed completions chunk.
func openaiCompletionChunkToChunk(c *openai.CompletionChunk) CompletionStreamChunk {
	return CompletionStreamChunk{
		ID:                c.ID,
		Object:            CompletionObject(c.Object),
		Created:           c.Created,
		Model:             c.Model,
		Choices:           openaiCompletionChunkChoicesToCompletionStreamChoices(c.Choices),
		SystemFingerprint: c.SystemFingerprint,
		Usage:             openaiCompletionUsageToCompletionUsage(c.Usage),
	}
}

// openaiCompletionChunkChoicesToCompletionStreamChoices converts the streamed
// choices.
func openaiCompletionChunkChoicesToCompletionStreamChoices(choices []openai.CompletionChunkChoice) []CompletionStreamChoice {
	if len(choices) == 0 {
		return nil
	}

	result := make([]CompletionStreamChoice, 0, len(choices))
	for i := range choices {
		result = append(result, CompletionStreamChoice{
			Index:        choices[i].Index,
			Text:         choices[i].Text,
			FinishReason: choices[i].FinishReason,
			Logprobs:     openaiCompletionLogprobsToCompletionLogprobs(choices[i].Logprobs),
		})
	}

	return result
}

// openaiEmbeddingResponseToResponse converts an OpenAI embeddings response into
// the openingrouter type.
func openaiEmbeddingResponseToResponse(r *openai.EmbeddingResponse) *EmbeddingResponse {
	if r == nil {
		return nil
	}

	return &EmbeddingResponse{
		ID:     r.ID,
		Object: EmbeddingObject(r.Object),
		Data:   openaiEmbeddingsToEmbeddings(r.Data),
		Model:  r.Model,
		Usage:  openaiEmbeddingUsageToEmbeddingUsage(r.Usage),
	}
}

// openaiEmbeddingsToEmbeddings converts the embedding vectors.
func openaiEmbeddingsToEmbeddings(embeddings []openai.Embedding) []Embedding {
	if len(embeddings) == 0 {
		return nil
	}

	result := make([]Embedding, 0, len(embeddings))
	for i := range embeddings {
		result = append(result, Embedding{
			Object: EmbeddingObject(embeddings[i].Object),
			Embedding: EmbeddingValue{
				Floats: embeddings[i].Embedding,
			},
			Index: embeddings[i].Index,
		})
	}

	return result
}

// openaiEmbeddingUsageToEmbeddingUsage converts the token usage. OpenAI does not
// report a cost, so Cost and IsBYOK stay at their zero values.
func openaiEmbeddingUsageToEmbeddingUsage(u *openai.EmbeddingUsage) *EmbeddingUsage {
	if u == nil {
		return nil
	}

	return &EmbeddingUsage{
		PromptTokens: u.PromptTokens,
		TotalTokens:  u.TotalTokens,
	}
}
