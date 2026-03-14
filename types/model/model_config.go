package model

import "github.com/fanlv/deep-agent-demo/pkg/modelbuilder"

type ModelInstance struct {
	ID              int64                        `json:"id"`
	ModelClass      modelbuilder.ModelClass      `json:"model_class"`
	DisplayName     string                       `json:"display_name"`
	Connection      *modelbuilder.ConnectionInfo `json:"connection"`
	ThinkingType    modelbuilder.ThinkingType    `json:"thinking_type,omitempty"`
	EnableBase64URL bool                         `json:"enable_base64_url,omitempty"`
	Status          int                          `json:"status"`
	CreatedAt       int64                        `json:"created_at"`
	UpdatedAt       int64                        `json:"updated_at"`
	DeletedAt       int64                        `json:"deleted_at,omitempty"`
}

type ProviderInfo struct {
	ModelClass  modelbuilder.ModelClass `json:"model_class"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	IconURL     string                  `json:"icon_url"`
}

type ProviderModelList struct {
	Provider  *ProviderInfo    `json:"provider"`
	ModelList []*ModelInstance `json:"model_list"`
}

type CreateModelRequest struct {
	ModelClass      modelbuilder.ModelClass      `json:"model_class"`
	DisplayName     string                       `json:"display_name"`
	Connection      *modelbuilder.ConnectionInfo `json:"connection"`
	ThinkingType    modelbuilder.ThinkingType    `json:"thinking_type,omitempty"`
	EnableBase64URL bool                         `json:"enable_base64_url,omitempty"`
}

var DefaultProviders = []ProviderInfo{
	{ModelClass: modelbuilder.ModelClassArk, Name: "豆包模型", Description: "火山引擎 Ark 大模型服务", IconURL: "https://lf-cdn.marscode.com.cn/obj/marscode-bucket-cn/images/doubao_icon.png"},
	{ModelClass: modelbuilder.ModelClassOpenAI, Name: "OpenAI", Description: "OpenAI GPT 系列模型", IconURL: "https://lf-cdn.marscode.com.cn/obj/marscode-bucket-cn/images/openai_icon.png"},
	{ModelClass: modelbuilder.ModelClassClaude, Name: "Claude", Description: "Anthropic Claude 系列模型", IconURL: "https://lf-cdn.marscode.com.cn/obj/marscode-bucket-cn/images/claude_icon.png"},
	{ModelClass: modelbuilder.ModelClassDeepSeek, Name: "DeepSeek", Description: "DeepSeek 深度求索", IconURL: "https://lf-cdn.marscode.com.cn/obj/marscode-bucket-cn/images/deepseek_icon.png"},
	{ModelClass: modelbuilder.ModelClassGemini, Name: "Gemini", Description: "Google Gemini 系列模型", IconURL: "https://lf-cdn.marscode.com.cn/obj/marscode-bucket-cn/images/gemini_icon.png"},
	{ModelClass: modelbuilder.ModelClassOllama, Name: "Ollama", Description: "本地部署模型", IconURL: "https://lf-cdn.marscode.com.cn/obj/marscode-bucket-cn/images/ollama_icon.png"},
	{ModelClass: modelbuilder.ModelClassQwen, Name: "通义千问", Description: "阿里云通义千问", IconURL: "https://lf-cdn.marscode.com.cn/obj/marscode-bucket-cn/images/qwen_icon.png"},
}
