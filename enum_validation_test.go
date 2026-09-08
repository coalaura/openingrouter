package openingrouter

import (
	"testing"
)

func TestValidateApiKeyEnums(t *testing.T) {
	apiKeyLimitResetValues := []ApiKeyLimitReset{
		ApiKeyLimitResetDaily,
		ApiKeyLimitResetWeekly,
		ApiKeyLimitResetMonthly,
	}

	tAssertEnumValidation(t, apiKeyLimitResetValues, IsValidApiKeyLimitReset)
}

func TestValidateCompletionsEnums(t *testing.T) {
	completionObjectValues := []CompletionObject{
		CompletionObjectTextCompletion,
	}

	tAssertEnumValidation(t, completionObjectValues, IsValidCompletionObject)
}

func TestValidateErrorEnums(t *testing.T) {
	errorTypeValues := []ErrorType{
		ErrorTypeContextLengthExceeded,
		ErrorTypeMaxTokensExceeded,
		ErrorTypeTokenLimitExceeded,
		ErrorTypeStringTooLong,
		ErrorTypeAuthentication,
		ErrorTypePermissionDenied,
		ErrorTypePaymentRequired,
		ErrorTypeRateLimitExceeded,
		ErrorTypeProviderOverloaded,
		ErrorTypeProviderUnavailable,
		ErrorTypeInvalidRequest,
		ErrorTypeInvalidPrompt,
		ErrorTypeNotFound,
		ErrorTypePreconditionFailed,
		ErrorTypePayloadTooLarge,
		ErrorTypeUnprocessable,
		ErrorTypeContentPolicy,
		ErrorTypeRefusal,
		ErrorTypeInvalidImage,
		ErrorTypeImageTooLarge,
		ErrorTypeImageTooSmall,
		ErrorTypeUnsupportedImage,
		ErrorTypeImageNotFound,
		ErrorTypeImageDownloadFailed,
		ErrorTypeServer,
		ErrorTypeTimeout,
		ErrorTypeUnmapped,
	}

	tAssertEnumValidation(t, errorTypeValues, IsValidErrorType)
}

func TestValidateImageModelsEnums(t *testing.T) {
	imageCapabilityTypeValues := []ImageCapabilityType{
		ImageCapabilityTypeBoolean,
		ImageCapabilityTypeEnum,
		ImageCapabilityTypeRange,
	}

	tAssertEnumValidation(t, imageCapabilityTypeValues, IsValidImageCapabilityType)
}

func TestValidateTTSEnums(t *testing.T) {
	speechResponseFormatValues := []SpeechResponseFormat{
		SpeechResponseFormatMP3,
		SpeechResponseFormatPCM,
	}

	tAssertEnumValidation(t, speechResponseFormatValues, IsValidSpeechResponseFormat)
}

func TestValidateSTTEnums(t *testing.T) {
	sttResponseFormatValues := []STTResponseFormat{
		STTResponseFormatJSON,
		STTResponseFormatVerboseJSON,
	}

	tAssertEnumValidation(t, sttResponseFormatValues, IsValidSTTResponseFormat)

	sttTimestampGranularityValues := []STTTimestampGranularity{
		STTTimestampGranularityWord,
		STTTimestampGranularitySegment,
	}

	tAssertEnumValidation(t, sttTimestampGranularityValues, IsValidSTTTimestampGranularity)
}

func TestValidateEmbeddingsEnums(t *testing.T) {
	embeddingEncodingFormatValues := []EmbeddingEncodingFormat{
		EmbeddingEncodingFormatFloat,
		EmbeddingEncodingFormatBase64,
	}

	tAssertEnumValidation(t, embeddingEncodingFormatValues, IsValidEmbeddingEncodingFormat)

	embeddingObjectValues := []EmbeddingObject{
		EmbeddingObjectList,
		EmbeddingObjectEmbedding,
	}

	tAssertEnumValidation(t, embeddingObjectValues, IsValidEmbeddingObject)

	embeddingContentPartTypeValues := []EmbeddingContentPartType{
		EmbeddingContentPartTypeText,
		EmbeddingContentPartTypeImageURL,
		EmbeddingContentPartTypeInputAudio,
		EmbeddingContentPartTypeInputVideo,
		EmbeddingContentPartTypeInputFile,
	}

	tAssertEnumValidation(t, embeddingContentPartTypeValues, IsValidEmbeddingContentPartType)
}

func TestValidateImageEnums(t *testing.T) {
	imageAspectRatioValues := []ImageAspectRatio{
		ImageAspectRatio1x1,
		ImageAspectRatio1x2,
		ImageAspectRatio1x4,
		ImageAspectRatio1x8,
		ImageAspectRatio2x1,
		ImageAspectRatio2x3,
		ImageAspectRatio3x2,
		ImageAspectRatio3x4,
		ImageAspectRatio4x1,
		ImageAspectRatio4x3,
		ImageAspectRatio4x5,
		ImageAspectRatio5x4,
		ImageAspectRatio8x1,
		ImageAspectRatio9x16,
		ImageAspectRatio16x9,
		ImageAspectRatio9x19Point5,
		ImageAspectRatio19Point5x9,
		ImageAspectRatio9x20,
		ImageAspectRatio20x9,
		ImageAspectRatio9x21,
		ImageAspectRatio21x9,
		ImageAspectRatioAuto,
	}

	tAssertEnumValidation(t, imageAspectRatioValues, IsValidImageAspectRatio)

	imageBackgroundValues := []ImageBackground{
		ImageBackgroundAuto,
		ImageBackgroundTransparent,
		ImageBackgroundOpaque,
	}

	tAssertEnumValidation(t, imageBackgroundValues, IsValidImageBackground)

	imageOutputFormatValues := []ImageOutputFormat{
		ImageOutputFormatPNG,
		ImageOutputFormatJPEG,
		ImageOutputFormatWebP,
		ImageOutputFormatSVG,
	}

	tAssertEnumValidation(t, imageOutputFormatValues, IsValidImageOutputFormat)

	imageQualityValues := []ImageQuality{
		ImageQualityAuto,
		ImageQualityLow,
		ImageQualityMedium,
		ImageQualityHigh,
	}

	tAssertEnumValidation(t, imageQualityValues, IsValidImageQuality)

	imageResolutionValues := []ImageResolution{
		ImageResolution512,
		ImageResolution1K,
		ImageResolution2K,
		ImageResolution4K,
	}

	tAssertEnumValidation(t, imageResolutionValues, IsValidImageResolution)

	imageStreamEventTypeValues := []ImageStreamEventType{
		ImageStreamEventTypePartialImage,
		ImageStreamEventTypeTextChunk,
		ImageStreamEventTypeCompleted,
		ImageStreamEventTypeError,
	}

	tAssertEnumValidation(t, imageStreamEventTypeValues, IsValidImageStreamEventType)

	imageStreamPhaseValues := []ImageStreamPhase{
		ImageStreamPhaseContent,
		ImageStreamPhaseReasoning,
		ImageStreamPhaseDraft,
	}

	tAssertEnumValidation(t, imageStreamPhaseValues, IsValidImageStreamPhase)
}

func TestValidateModelsEnums(t *testing.T) {
	parameterValues := []Parameter{
		ParameterTemperature,
		ParameterTopP,
		ParameterTopK,
		ParameterMinP,
		ParameterTopA,
		ParameterFrequencyPenalty,
		ParameterPresencePenalty,
		ParameterRepetitionPenalty,
		ParameterMaxTokens,
		ParameterMaxCompletionTokens,
		ParameterLogitBias,
		ParameterLogprobs,
		ParameterTopLogprobs,
		ParameterPrediction,
		ParameterSeed,
		ParameterResponseFormat,
		ParameterStructuredOutputs,
		ParameterStop,
		ParameterTools,
		ParameterToolChoice,
		ParameterParallelToolCalls,
		ParameterIncludeReasoning,
		ParameterReasoning,
		ParameterReasoningEffort,
		ParameterWebSearchOptions,
		ParameterVerbosity,
	}

	tAssertEnumValidation(t, parameterValues, IsValidParameter)

	instructTypeValues := []InstructType{
		InstructTypeNone,
		InstructTypeAiroboros,
		InstructTypeAlpaca,
		InstructTypeAlpacaModif,
		InstructTypeChatML,
		InstructTypeClaude,
		InstructTypeCodeLlama,
		InstructTypeGemma,
		InstructTypeLlama2,
		InstructTypeLlama3,
		InstructTypeMistral,
		InstructTypeNemotron,
		InstructTypeNeural,
		InstructTypeOpenChat,
		InstructTypePhi3,
		InstructTypeRWKV,
		InstructTypeVicuna,
		InstructTypeZephyr,
		InstructTypeDeepSeekR1,
		InstructTypeDeepSeekV31,
		InstructTypeQwQ,
		InstructTypeQwen3,
	}

	tAssertEnumValidation(t, instructTypeValues, IsValidInstructType)

	modelGroupValues := []ModelGroup{
		ModelGroupRouter,
		ModelGroupMedia,
		ModelGroupOther,
		ModelGroupGPT,
		ModelGroupClaude,
		ModelGroupGemini,
		ModelGroupGemma,
		ModelGroupGrok,
		ModelGroupCohere,
		ModelGroupNova,
		ModelGroupQwen,
		ModelGroupYi,
		ModelGroupDeepSeek,
		ModelGroupMistral,
		ModelGroupLlama2,
		ModelGroupLlama3,
		ModelGroupLlama4,
		ModelGroupPaLM,
		ModelGroupRWKV,
		ModelGroupQwen3,
	}

	tAssertEnumValidation(t, modelGroupValues, IsValidModelGroup)

	reasoningEffortValues := []ReasoningEffort{
		ReasoningEffortMax,
		ReasoningEffortXHigh,
		ReasoningEffortHigh,
		ReasoningEffortMedium,
		ReasoningEffortLow,
		ReasoningEffortMinimal,
		ReasoningEffortNone,
	}

	tAssertEnumValidation(t, reasoningEffortValues, IsValidReasoningEffort)

	modelCategoryValues := []ModelCategory{
		ModelCategoryProgramming,
		ModelCategoryRoleplay,
		ModelCategoryMarketing,
		ModelCategoryMarketingSEO,
		ModelCategoryTechnology,
		ModelCategoryScience,
		ModelCategoryTranslation,
		ModelCategoryLegal,
		ModelCategoryFinance,
		ModelCategoryHealth,
		ModelCategoryTrivia,
		ModelCategoryAcademia,
	}

	tAssertEnumValidation(t, modelCategoryValues, IsValidModelCategory)

	modelSortValues := []ModelSort{
		ModelSortMostPopular,
		ModelSortNewest,
		ModelSortTopWeekly,
		ModelSortPricingLowToHigh,
		ModelSortPricingHighToLow,
		ModelSortContextHighToLow,
		ModelSortThroughputHighToLow,
		ModelSortLatencyLowToHigh,
		ModelSortIntelligenceHighToLow,
		ModelSortCodingHighToLow,
		ModelSortAgenticHighToLow,
		ModelSortDesignArenaELOHighToLow,
	}

	tAssertEnumValidation(t, modelSortValues, IsValidModelSort)

	modelRegionValues := []ModelRegion{
		ModelRegionEU,
	}

	tAssertEnumValidation(t, modelRegionValues, IsValidModelRegion)
}

func TestValidateCommonEnums(t *testing.T) {
	inputModalityValues := []InputModality{
		InputModalityText,
		InputModalityImage,
		InputModalityFile,
		InputModalityAudio,
		InputModalityVideo,
	}

	tAssertEnumValidation(t, inputModalityValues, IsValidInputModality)

	outputModalityValues := []OutputModality{
		OutputModalityText,
		OutputModalityImage,
		OutputModalityEmbeddings,
		OutputModalityAudio,
		OutputModalityVideo,
		OutputModalityRerank,
		OutputModalitySpeech,
		OutputModalityTranscription,
	}

	tAssertEnumValidation(t, outputModalityValues, IsValidOutputModality)

	contentPartTypeValues := []ContentPartType{
		ContentPartTypeImageURL,
	}

	tAssertEnumValidation(t, contentPartTypeValues, IsValidContentPartType)

	providerDataCollectionValues := []ProviderDataCollection{
		ProviderDataCollectionAllow,
		ProviderDataCollectionDeny,
	}

	tAssertEnumValidation(t, providerDataCollectionValues, IsValidProviderDataCollection)

	quantizationValues := []Quantization{
		QuantizationInt4,
		QuantizationInt8,
		QuantizationFP4,
		QuantizationMXFP4,
		QuantizationNVFP4,
		QuantizationFP6,
		QuantizationFP8,
		QuantizationMXFP8,
		QuantizationFP16,
		QuantizationBF16,
		QuantizationFP32,
		QuantizationUnknown,
	}

	tAssertEnumValidation(t, quantizationValues, IsValidQuantization)

	providerSortValues := []ProviderSort{
		ProviderSortPrice,
		ProviderSortThroughput,
		ProviderSortLatency,
		ProviderSortExacto,
	}

	tAssertEnumValidation(t, providerSortValues, IsValidProviderSort)

	providerSortPartitionValues := []ProviderSortPartition{
		ProviderSortPartitionModel,
		ProviderSortPartitionNone,
	}

	tAssertEnumValidation(t, providerSortPartitionValues, IsValidProviderSortPartition)

	anthropicSpeedValues := []AnthropicSpeed{
		AnthropicSpeedFast,
		AnthropicSpeedStandard,
	}

	tAssertEnumValidation(t, anthropicSpeedValues, IsValidAnthropicSpeed)

	anthropicUsageIterationTypeValues := []AnthropicUsageIterationType{
		AnthropicUsageIterationTypeCompaction,
		AnthropicUsageIterationTypeMessage,
		AnthropicUsageIterationTypeAdvisorMessage,
	}

	tAssertEnumValidation(t, anthropicUsageIterationTypeValues, IsValidAnthropicUsageIterationType)
}

func TestValidateChatEnums(t *testing.T) {
	chatMetadataLevelValues := []ChatMetadataLevel{
		ChatMetadataLevelDisabled,
		ChatMetadataLevelEnabled,
	}

	tAssertEnumValidation(t, chatMetadataLevelValues, IsValidChatMetadataLevel)

	chatRoleValues := []ChatRole{
		ChatRoleSystem,
		ChatRoleDeveloper,
		ChatRoleUser,
		ChatRoleAssistant,
		ChatRoleTool,
	}

	tAssertEnumValidation(t, chatRoleValues, IsValidChatRole)

	chatContentPartTypeValues := []ChatContentPartType{
		ChatContentPartTypeText,
		ChatContentPartTypeImageURL,
		ChatContentPartTypeInputAudio,
		ChatContentPartTypeVideoURL,
		ChatContentPartTypeFile,
		ChatContentPartTypeInputVideo,
	}

	tAssertEnumValidation(t, chatContentPartTypeValues, IsValidChatContentPartType)

	chatImageDetailValues := []ChatImageDetail{
		ChatImageDetailAuto,
		ChatImageDetailLow,
		ChatImageDetailHigh,
		ChatImageDetailOriginal,
	}

	tAssertEnumValidation(t, chatImageDetailValues, IsValidChatImageDetail)

	anthropicCacheControlTypeValues := []AnthropicCacheControlType{
		AnthropicCacheControlTypeEphemeral,
	}

	tAssertEnumValidation(t, anthropicCacheControlTypeValues, IsValidAnthropicCacheControlType)

	anthropicCacheTTLValues := []AnthropicCacheTTL{
		AnthropicCacheTTL5M,
		AnthropicCacheTTL1H,
	}

	tAssertEnumValidation(t, anthropicCacheTTLValues, IsValidAnthropicCacheTTL)

	chatPromptCacheModeValues := []ChatPromptCacheMode{
		ChatPromptCacheModeExplicit,
	}

	tAssertEnumValidation(t, chatPromptCacheModeValues, IsValidChatPromptCacheMode)

	chatServiceTierValues := []ChatServiceTier{
		ChatServiceTierAuto,
		ChatServiceTierDefault,
		ChatServiceTierFlex,
		ChatServiceTierPriority,
		ChatServiceTierScale,
	}

	tAssertEnumValidation(t, chatServiceTierValues, IsValidChatServiceTier)

	chatReasoningSummaryValues := []ChatReasoningSummary{
		ChatReasoningSummaryAuto,
		ChatReasoningSummaryConcise,
		ChatReasoningSummaryDetailed,
	}

	tAssertEnumValidation(t, chatReasoningSummaryValues, IsValidChatReasoningSummary)

	chatPredictionTypeValues := []ChatPredictionType{
		ChatPredictionTypeContent,
	}

	tAssertEnumValidation(t, chatPredictionTypeValues, IsValidChatPredictionType)

	chatResponseFormatTypeValues := []ChatResponseFormatType{
		ChatResponseFormatTypeText,
		ChatResponseFormatTypeJSONObject,
		ChatResponseFormatTypeJSONSchema,
		ChatResponseFormatTypeGrammar,
		ChatResponseFormatTypePython,
	}

	tAssertEnumValidation(t, chatResponseFormatTypeValues, IsValidChatResponseFormatType)

	chatRouteValues := []ChatRoute{
		ChatRouteFallback,
		ChatRouteSort,
	}

	tAssertEnumValidation(t, chatRouteValues, IsValidChatRoute)

	chatStopConditionTypeValues := []ChatStopConditionType{
		ChatStopConditionTypeStepCountIs,
		ChatStopConditionTypeHasToolCall,
		ChatStopConditionTypeMaxTokensUsed,
		ChatStopConditionTypeMaxCost,
		ChatStopConditionTypeFinishReasonIs,
	}

	tAssertEnumValidation(t, chatStopConditionTypeValues, IsValidChatStopConditionType)

	chatToolTypeValues := []ChatToolType{
		ChatToolTypeFunction,
		ChatToolTypeAdvisor,
		ChatToolTypeBash,
		ChatToolTypeDatetime,
		ChatToolTypeFiles,
		ChatToolTypeFusion,
		ChatToolTypeImageGeneration,
		ChatToolTypeSearchModels,
		ChatToolTypeSubagent,
		ChatToolTypeWebFetch,
		ChatToolTypeWebSearch,
		ChatToolTypeWebSearchShorthand,
		ChatToolTypeWebSearchPreview,
		ChatToolTypeWebSearchPreview20250311,
		ChatToolTypeWebSearch20250826,
	}

	tAssertEnumValidation(t, chatToolTypeValues, IsValidChatToolType)

	chatToolChoiceModeValues := []ChatToolChoiceMode{
		ChatToolChoiceModeNone,
		ChatToolChoiceModeAuto,
		ChatToolChoiceModeRequired,
	}

	tAssertEnumValidation(t, chatToolChoiceModeValues, IsValidChatToolChoiceMode)

	chatBashEngineValues := []ChatBashEngine{
		ChatBashEngineAuto,
		ChatBashEngineNative,
		ChatBashEngineOpenRouter,
	}

	tAssertEnumValidation(t, chatBashEngineValues, IsValidChatBashEngine)

	chatBashEnvironmentTypeValues := []ChatBashEnvironmentType{
		ChatBashEnvironmentTypeContainerAuto,
		ChatBashEnvironmentTypeContainerReference,
	}

	tAssertEnumValidation(t, chatBashEnvironmentTypeValues, IsValidChatBashEnvironmentType)

	chatWebSearchEngineValues := []ChatWebSearchEngine{
		ChatWebSearchEngineAuto,
		ChatWebSearchEngineNative,
		ChatWebSearchEngineExa,
		ChatWebSearchEngineParallel,
		ChatWebSearchEngineFirecrawl,
		ChatWebSearchEnginePerplexity,
	}

	tAssertEnumValidation(t, chatWebSearchEngineValues, IsValidChatWebSearchEngine)

	chatWebFetchEngineValues := []ChatWebFetchEngine{
		ChatWebFetchEngineAuto,
		ChatWebFetchEngineNative,
		ChatWebFetchEngineOpenRouter,
		ChatWebFetchEngineExa,
		ChatWebFetchEngineParallel,
		ChatWebFetchEngineFirecrawl,
	}

	tAssertEnumValidation(t, chatWebFetchEngineValues, IsValidChatWebFetchEngine)

	chatSearchContextSizeValues := []ChatSearchContextSize{
		ChatSearchContextSizeLow,
		ChatSearchContextSizeMedium,
		ChatSearchContextSizeHigh,
	}

	tAssertEnumValidation(t, chatSearchContextSizeValues, IsValidChatSearchContextSize)

	chatUserLocationTypeValues := []ChatUserLocationType{
		ChatUserLocationTypeApproximate,
	}

	tAssertEnumValidation(t, chatUserLocationTypeValues, IsValidChatUserLocationType)

	chatPluginIDValues := []ChatPluginID{
		ChatPluginIDAutoRouter,
		ChatPluginIDAutoBetaRouter,
		ChatPluginIDModeration,
		ChatPluginIDWeb,
		ChatPluginIDWebFetch,
		ChatPluginIDFileParser,
		ChatPluginIDResponseHealing,
		ChatPluginIDContextCompression,
		ChatPluginIDParetoRouter,
		ChatPluginIDFusion,
	}

	tAssertEnumValidation(t, chatPluginIDValues, IsValidChatPluginID)

	chatCostTierValues := []ChatCostTier{
		ChatCostTierLow,
		ChatCostTierMedium,
		ChatCostTierHigh,
		ChatCostTierXHigh,
		ChatCostTierMax,
	}

	tAssertEnumValidation(t, chatCostTierValues, IsValidChatCostTier)

	chatContextCompressionEngineValues := []ChatContextCompressionEngine{
		ChatContextCompressionEngineMiddleOut,
	}

	tAssertEnumValidation(t, chatContextCompressionEngineValues, IsValidChatContextCompressionEngine)

	chatPDFParserEngineValues := []ChatPDFParserEngine{
		ChatPDFParserEngineMistralOCR,
		ChatPDFParserEngineNative,
		ChatPDFParserEngineCloudflareAI,
		ChatPDFParserEnginePDFText,
	}

	tAssertEnumValidation(t, chatPDFParserEngineValues, IsValidChatPDFParserEngine)

	chatFusionPresetValues := []ChatFusionPreset{
		ChatFusionPresetGeneralHigh,
		ChatFusionPresetGeneralBudget,
		ChatFusionPresetGeneralFast,
	}

	tAssertEnumValidation(t, chatFusionPresetValues, IsValidChatFusionPreset)

	chatParetoPriceSourceValues := []ChatParetoPriceSource{
		ChatParetoPriceSourcePrompt,
		ChatParetoPriceSourceWeightedAvg,
	}

	tAssertEnumValidation(t, chatParetoPriceSourceValues, IsValidChatParetoPriceSource)

	chatObjectValues := []ChatObject{
		ChatObjectCompletion,
		ChatObjectCompletionChunk,
	}

	tAssertEnumValidation(t, chatObjectValues, IsValidChatObject)

	chatFinishReasonValues := []ChatFinishReason{
		ChatFinishReasonToolCalls,
		ChatFinishReasonStop,
		ChatFinishReasonLength,
		ChatFinishReasonContentFilter,
		ChatFinishReasonError,
	}

	tAssertEnumValidation(t, chatFinishReasonValues, IsValidChatFinishReason)

	chatReasoningDetailTypeValues := []ChatReasoningDetailType{
		ChatReasoningDetailTypeSummary,
		ChatReasoningDetailTypeEncrypted,
		ChatReasoningDetailTypeText,
		ChatReasoningDetailTypeServerToolCall,
	}

	tAssertEnumValidation(t, chatReasoningDetailTypeValues, IsValidChatReasoningDetailType)

	chatReasoningFormatValues := []ChatReasoningFormat{
		ChatReasoningFormatUnknown,
		ChatReasoningFormatOpenAIResponsesV1,
		ChatReasoningFormatAzureOpenAIResponsesV1,
		ChatReasoningFormatBedrockOpenAIResponsesV1,
		ChatReasoningFormatXAIResponsesV1,
		ChatReasoningFormatMetaResponsesV1,
		ChatReasoningFormatAnthropicClaudeV1,
		ChatReasoningFormatGoogleGeminiV1,
	}

	tAssertEnumValidation(t, chatReasoningFormatValues, IsValidChatReasoningFormat)

	chatErrorTypeValues := []ChatErrorType{
		ChatErrorTypeContextLengthExceeded,
		ChatErrorTypeMaxTokensExceeded,
		ChatErrorTypeTokenLimitExceeded,
		ChatErrorTypeStringTooLong,
		ChatErrorTypeAuthentication,
		ChatErrorTypePermissionDenied,
		ChatErrorTypePaymentRequired,
		ChatErrorTypeRateLimitExceeded,
		ChatErrorTypeProviderOverloaded,
		ChatErrorTypeProviderUnavailable,
		ChatErrorTypeInvalidRequest,
		ChatErrorTypeInvalidPrompt,
		ChatErrorTypeNotFound,
		ChatErrorTypePreconditionFailed,
		ChatErrorTypePayloadTooLarge,
		ChatErrorTypeUnprocessable,
		ChatErrorTypeContentPolicyViolation,
		ChatErrorTypeRefusal,
		ChatErrorTypeInvalidImage,
		ChatErrorTypeImageTooLarge,
		ChatErrorTypeImageTooSmall,
		ChatErrorTypeUnsupportedImageFormat,
		ChatErrorTypeImageNotFound,
		ChatErrorTypeImageDownloadFailed,
		ChatErrorTypeServer,
		ChatErrorTypeTimeout,
		ChatErrorTypeUnmapped,
	}

	tAssertEnumValidation(t, chatErrorTypeValues, IsValidChatErrorType)

	routingStrategyValues := []RoutingStrategy{
		RoutingStrategyDirect,
		RoutingStrategyAuto,
		RoutingStrategyFree,
		RoutingStrategyLatest,
		RoutingStrategyAlias,
		RoutingStrategyFallback,
		RoutingStrategyPareto,
		RoutingStrategyBodybuilder,
		RoutingStrategyFusion,
	}

	tAssertEnumValidation(t, routingStrategyValues, IsValidRoutingStrategy)

	pipelineStageTypeValues := []PipelineStageType{
		PipelineStageTypeGuardrail,
		PipelineStageTypePlugin,
		PipelineStageTypeServerTools,
		PipelineStageTypeResponseHealing,
		PipelineStageTypeContextCompression,
	}

	tAssertEnumValidation(t, pipelineStageTypeValues, IsValidPipelineStageType)
}
