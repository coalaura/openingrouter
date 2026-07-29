package openingrouter

// SystemMessage creates a new system message with the given text content.
func SystemMessage(content string) ChatMessage {
	return ChatMessage{
		Role: ChatRoleSystem,
		Content: ChatContent{
			Text: content,
		},
	}
}

// UserMessage creates a new user message with the given text content.
func UserMessage(content string) ChatMessage {
	return ChatMessage{
		Role: ChatRoleUser,
		Content: ChatContent{
			Text: content,
		},
	}
}

// AssistantMessage creates a new assistant message with the given text content.
func AssistantMessage(content string) ChatMessage {
	return ChatMessage{
		Role: ChatRoleAssistant,
		Content: ChatContent{
			Text: content,
		},
	}
}

// ToolMessage creates a new tool (response) message with a call ID and content.
func ToolMessage(callID string, content string) ChatMessage {
	return ChatMessage{
		Role: ChatRoleTool,
		Content: ChatContent{
			Text: content,
		},
		ToolCallID: callID,
	}
}

// UserMessageWithImage creates a new user message with text and image URL.
func UserMessageWithImage(text, imageURL string) ChatMessage {
	return ChatMessage{
		Role: ChatRoleUser,
		Content: ChatContent{
			Parts: []ChatContentPart{
				{
					Type: ChatContentPartTypeText,
					Text: text,
				},
				{
					Type: ChatContentPartTypeImageURL,
					ImageURL: &ChatContentImageURL{
						URL: imageURL,
					},
				},
			},
		},
	}
}
