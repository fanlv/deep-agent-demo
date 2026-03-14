import { useEffect, useRef, useCallback } from 'react';
import { Message } from '../types';
import { MessageItem } from './MessageItem';
import './MessageList.css';

interface MessageListProps {
  messages: Message[];
  isLoading: boolean;
  onSendMessage?: (message: string) => void;
}

export function MessageList({ messages, isLoading, onSendMessage }: MessageListProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const prevMessagesLengthRef = useRef(messages.length);

  const checkIfNearBottom = useCallback(() => {
    if (!containerRef.current) return true;
    const { scrollTop, scrollHeight, clientHeight } = containerRef.current;
    const threshold = 20;
    return scrollHeight - scrollTop - clientHeight < threshold;
  }, []);

  useEffect(() => {
    if (!containerRef.current) return;
    const hasNewMessages = messages.length > prevMessagesLengthRef.current;
    if (hasNewMessages || checkIfNearBottom()) {
      containerRef.current.scrollTop = containerRef.current.scrollHeight;
    }

    prevMessagesLengthRef.current = messages.length;
  }, [messages, checkIfNearBottom]);

  const suggestions = [
    { text: '查询北京天气', emoji: '☀️' },
    { text: '搜索 AI 最新进展', emoji: '' },
    { text: '你好', emoji: '' },
  ];

  const handleSuggestionClick = (text: string) => {
    if (onSendMessage && !isLoading) {
      onSendMessage(text);
    }
  };

  return (
    <div className="message-list" ref={containerRef}>
      {messages.length === 0 ? (
        <div className="empty-state">
          <div className="empty-icon">⚙️</div>
          <h3>Agent UI Protocol Demo</h3>
          <p>支持任意响应 AgentUIProtocol 事件的 SSE 接口</p>
          <div className="suggestions">
            {suggestions.map((item, index) => (
              <span
                key={index}
                className="suggestion-tag"
                onClick={() => handleSuggestionClick(item.text)}
              >
                {item.text} {item.emoji}
              </span>
            ))}
          </div>
        </div>
      ) : (
        <>
          {messages.map((message) => (
            <MessageItem key={message.id} message={message} />
          ))}
          {isLoading && messages.length > 0 && (
            <div className="loading-indicator">
              <div className="loading-dots">
                <span />
                <span />
                <span />
              </div>
              <span>AI 正在思考...</span>
            </div>
          )}
        </>
      )}
    </div>
  );
}
