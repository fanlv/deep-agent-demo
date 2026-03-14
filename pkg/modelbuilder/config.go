package modelbuilder

import "os"

func LoadConfigFromEnv() *ModelConfig {
	if modelName := os.Getenv("ARK_MODEL_ID"); modelName != "" {
		cfg := &ModelConfig{
			ModelClass: ModelClassArk,
			Connection: &ConnectionInfo{
				APIKey:  os.Getenv("ARK_API_KEY"),
				BaseURL: os.Getenv("ARK_BASE_URL"),
				Model:   modelName,
				Ark: &ArkConnectionInfo{
					Region: os.Getenv("ARK_REGION"),
				},
			},
		}
		if os.Getenv("ARK_DISABLE_THINKING") == "true" {
			cfg.ThinkingType = ThinkingTypeDisable
		}
		return cfg
	}

	if modelName := os.Getenv("OPENAI_MODEL"); modelName != "" {
		cfg := &ModelConfig{
			ModelClass: ModelClassOpenAI,
			Connection: &ConnectionInfo{
				APIKey:  os.Getenv("OPENAI_API_KEY"),
				BaseURL: os.Getenv("OPENAI_BASE_URL"),
				Model:   modelName,
				OpenAI: &OpenAIConnectionInfo{
					ByAzure: os.Getenv("OPENAI_BY_AZURE") == "true",
				},
			},
		}
		return cfg
	}

	if modelName := os.Getenv("CLAUDE_MODEL"); modelName != "" {
		return &ModelConfig{
			ModelClass: ModelClassClaude,
			Connection: &ConnectionInfo{
				APIKey:  os.Getenv("CLAUDE_API_KEY"),
				BaseURL: os.Getenv("CLAUDE_BASE_URL"),
				Model:   modelName,
			},
		}
	}

	if modelName := os.Getenv("DEEPSEEK_MODEL"); modelName != "" {
		return &ModelConfig{
			ModelClass: ModelClassDeepSeek,
			Connection: &ConnectionInfo{
				APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
				BaseURL: os.Getenv("DEEPSEEK_BASE_URL"),
				Model:   modelName,
			},
		}
	}

	if modelName := os.Getenv("GEMINI_MODEL"); modelName != "" {
		return &ModelConfig{
			ModelClass: ModelClassGemini,
			Connection: &ConnectionInfo{
				APIKey:  os.Getenv("GEMINI_API_KEY"),
				BaseURL: os.Getenv("GEMINI_BASE_URL"),
				Model:   modelName,
				Gemini: &GeminiConnectionInfo{
					Backend:  os.Getenv("GEMINI_BACKEND"),
					Project:  os.Getenv("GEMINI_PROJECT"),
					Location: os.Getenv("GEMINI_LOCATION"),
				},
			},
		}
	}

	if modelName := os.Getenv("OLLAMA_MODEL"); modelName != "" {
		baseURL := os.Getenv("OLLAMA_BASE_URL")
		if baseURL == "" {
			baseURL = "http://127.0.0.1:11434"
		}
		return &ModelConfig{
			ModelClass: ModelClassOllama,
			Connection: &ConnectionInfo{
				BaseURL: baseURL,
				Model:   modelName,
			},
		}
	}

	if modelName := os.Getenv("QWEN_MODEL"); modelName != "" {
		return &ModelConfig{
			ModelClass: ModelClassQwen,
			Connection: &ConnectionInfo{
				APIKey:  os.Getenv("QWEN_API_KEY"),
				BaseURL: os.Getenv("QWEN_BASE_URL"),
				Model:   modelName,
			},
		}
	}

	return nil
}
