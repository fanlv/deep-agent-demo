import { useEffect, useRef, useState, KeyboardEvent } from 'react';
import { ModelOption } from '../utils/models';
import './ChatInput.css';

interface ChatInputProps {
  onSend: (message: string) => void;
  onStop?: () => void;
  isLoading: boolean;
  disabled?: boolean;
  placeholder?: string;
  models?: ModelOption[];
  selectedModelId?: number | null;
  onSelectModel?: (modelId: number) => void;
  totalTokens?: number;
}

export function ChatInput({
  onSend,
  onStop,
  isLoading,
  disabled = false,
  placeholder = 'What can I help you?',
  models,
  selectedModelId,
  onSelectModel,
  totalTokens = 0,
}: ChatInputProps) {
  const [input, setInput] = useState('');
  const [showDropdown, setShowDropdown] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setShowDropdown(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleSend = () => {
    if (input.trim() && !isLoading && !disabled) {
      onSend(input.trim());
      setInput('');
    }
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const selectedModel = (models || []).find((m) => m.model.id === selectedModelId);
  const canSelectModel = !!models && models.length > 0 && !!onSelectModel;

  return (
    <div className="chat-input-container">
      <div className="chat-input-wrapper">
        <textarea
          className="chat-input"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={placeholder}
          disabled={isLoading || disabled}
          rows={1}
        />
        <div className="chat-input-footer">
          <div className="chat-input-options">
            {models && (
              <div className="chat-model-selector" ref={dropdownRef}>
                <div
                  className={`model-tag ${canSelectModel ? '' : 'disabled'}`}
                  onClick={() => {
                    if (canSelectModel) {
                      setShowDropdown(!showDropdown);
                    }
                  }}
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
                            onSelectModel?.(item.model.id);
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
            )}
            {totalTokens > 0 && (
              <span className="token-usage">
                Tokens: {totalTokens >= 1000 ? `${Math.floor(totalTokens / 1000)}K` : totalTokens}
              </span>
            )}
          </div>
          <div className="chat-input-actions">
            {isLoading ? (
              <button className="chat-btn stop-btn" onClick={onStop} title="停止生成">
                <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                  <rect x="6" y="6" width="12" height="12" rx="2" />
                </svg>
              </button>
            ) : (
              <button
                className="chat-btn send-btn"
                onClick={handleSend}
                disabled={!input.trim()}
                title="发送消息"
              >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                  <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" />
                </svg>
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
