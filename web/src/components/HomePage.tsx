import { useState, useEffect, useRef } from 'react';
import './HomePage.css';
import { fetchModelOptions, ModelOption } from '../utils/models';

interface HomePageProps {
  onStartChat: (message: string, modelId: number) => void;
  isInitializing?: boolean;
}

export function HomePage({ onStartChat, isInitializing }: HomePageProps) {
  const [input, setInput] = useState('');
  const [models, setModels] = useState<ModelOption[]>([]);
  const [selectedModelId, setSelectedModelId] = useState<number | null>(null);
  const [showDropdown, setShowDropdown] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    fetchModelOptions().then((allModels) => {
      setModels(allModels);
      setSelectedModelId((prev) => {
        if (prev || allModels.length === 0) {
          return prev;
        }
        return allModels[0].model.id;
      });
    }).catch((err) => {
      console.error('Failed to fetch model options:', err);
    });
  }, []);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setShowDropdown(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const selectedModel = models.find((m) => m.model.id === selectedModelId);

  const handleSubmit = () => {
    if (input.trim() && !isInitializing && selectedModelId) {
      onStartChat(input.trim(), selectedModelId);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit();
    }
  };

  const suggestions = [
    { text: '今天北京天气怎么样？', emoji: '☀️' },
    { text: '帮我搜索 AI Agent 最新进展', emoji: '🔍' },
    { text: '你能帮我做什么？', emoji: '🤔' },
  ];

  const handleSuggestionClick = (text: string) => {
    if (!isInitializing && selectedModelId) {
      onStartChat(text, selectedModelId);
    }
  };

  return (
    <div className="home-page">
      <div className="home-content">
        <div className="home-header">
          <div className="home-logo">🤖</div>
          <h1 className="home-title">What's on your mind?</h1>
          <p className="home-subtitle">Trusted AI Work Platform for Every Developer</p>
        </div>

        <div className="home-input-wrapper">
          <textarea
            className="home-input"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Ask anything (Press Shift + Enter for a new line)"
            disabled={isInitializing}
            rows={1}
          />
          <div className="home-input-footer">
            <div className="home-input-options" ref={dropdownRef}>
              <div
                className="model-tag"
                onClick={() => setShowDropdown(!showDropdown)}
              >
                {selectedModel && (
                  <img src={selectedModel.iconUrl} alt="" className="model-tag-icon" referrerPolicy="no-referrer" />
                )}
                <span>{selectedModel ? selectedModel.model.display_name : 'Select Model'}</span>
                <svg className="model-tag-arrow" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M6 9l6 6 6-6" />
                </svg>
              </div>
              {showDropdown && (
                <div className="model-dropdown">
                  {models.length === 0 ? (
                    <div className="model-dropdown-empty">No models available</div>
                  ) : (
                    models.map((item) => (
                      <div
                        key={item.model.id}
                        className={`model-dropdown-item ${item.model.id === selectedModelId ? 'active' : ''}`}
                        onClick={() => {
                          setSelectedModelId(item.model.id);
                          setShowDropdown(false);
                        }}
                      >
                        <img src={item.iconUrl} alt="" className="model-dropdown-icon" referrerPolicy="no-referrer" />
                        <div className="model-dropdown-info">
                          <span className="model-dropdown-name">{item.model.display_name}</span>
                          <span className="model-dropdown-provider">{item.providerName}</span>
                        </div>
                        {item.model.id === selectedModelId && (
                          <svg className="model-dropdown-check" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                            <path d="M20 6L9 17l-5-5" />
                          </svg>
                        )}
                      </div>
                    ))
                  )}
                </div>
              )}
            </div>
            <button
              className="home-send-btn"
              onClick={handleSubmit}
              disabled={!input.trim() || isInitializing || !selectedModelId}
            >
              <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M22 2L11 13" />
                <path d="M22 2L15 22L11 13L2 9L22 2Z" />
              </svg>
            </button>
          </div>
        </div>

        <div className="home-suggestions">
          {suggestions.map((item, index) => (
            <span
              key={index}
              className="home-suggestion-tag"
              onClick={() => handleSuggestionClick(item.text)}
            >
              {item.emoji} {item.text}
            </span>
          ))}
        </div>
      </div>
    </div>
  );
}
