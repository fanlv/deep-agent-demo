import { useState } from 'react';
import { Message, AssistantMessage, ToolMessage, MessageRoleEnum, MessageStatusEnum, ToolCallStatusEnum } from '../types';
import './MessageItem.css';

function showToast(message: string) {
  const existing = document.querySelector('.copy-toast');
  if (existing) {
    existing.remove();
  }

  const toast = document.createElement('div');
  toast.className = 'copy-toast';
  toast.textContent = message;
  document.body.appendChild(toast);

  setTimeout(() => {
    toast.classList.add('show');
  }, 10);

  setTimeout(() => {
    toast.classList.remove('show');
    setTimeout(() => toast.remove(), 300);
  }, 2000);
}

interface MessageItemProps {
  message: Message;
}

function UserMessageContent({ message }: { message: Message }) {
  return (
    <div className="message-item user-message">
      <div className="message-content">
        <div className="message-bubble user-bubble">
          {message.content}
        </div>
      </div>
    </div>
  );
}

function AssistantMessageContent({ message }: { message: AssistantMessage }) {
  return (
    <div className="message-item assistant-message">
      <div className="assistant-content-wrapper">
        {message.thinkingContent && (
          <div className="thinking-block">
            <div className="thinking-header">
              <span className="thinking-icon">💭</span>
              <span>深度思考</span>
              {message.isThinking && <span className="thinking-indicator" />}
            </div>
            <div className="thinking-content">
              {message.thinkingContent}
            </div>
          </div>
        )}
        {message.content && (
          <div className="assistant-bubble">
            <div className="markdown-content">
              {message.content.split('\n').map((line, i) => (
                <p key={i}>{formatLine(line)}</p>
              ))}
              {message.status === MessageStatusEnum.Started && (
                <span className="typing-indicator" />
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

function getToolIcon(toolName: string): string {
  const iconMap: Record<string, string> = {
    Read: '📖',
    Edit: '✏️',
    Write: '📝',
    Glob: '📂',
    Grep: '🔍',
    WebSearch: '🌐',
    WebFetch: '🌍',
    Bash: '💻',
    browser_click: '🖱️',
    browser_evaluate: '⚙️',
    browser_get_html: '📄',
    browser_get_page_info: 'ℹ️',
    browser_get_title: '📑',
    browser_get_url: '🔗',
    browser_navigate: '🧭',
    browser_pdf: '📋',
    browser_screenshot: '📸',
    browser_scroll: '📜',
    browser_type: '⌨️',
    browser_wait_visible: '👁️',
  };
  return iconMap[toolName] || '🔧';
}

function ToolMessageContent({ message }: { message: ToolMessage }) {
  const isCompleted = message.toolCallStatus !== ToolCallStatusEnum.Processing;
  const [isExpanded, setIsExpanded] = useState(!isCompleted);

  const statusText = {
    [ToolCallStatusEnum.Processing]: 'Running',
    [ToolCallStatusEnum.Success]: 'Completed',
    [ToolCallStatusEnum.Error]: 'Error',
  };

  const statusClass = {
    [ToolCallStatusEnum.Processing]: 'processing',
    [ToolCallStatusEnum.Success]: 'success',
    [ToolCallStatusEnum.Error]: 'error',
  };

  let parsedArgs: string | object | null = null;
  let parsedResult: string | object | null = null;

  try {
    parsedArgs = message.toolCallArgs ? JSON.parse(message.toolCallArgs) : null;
  } catch {
    parsedArgs = message.toolCallArgs;
  }

  try {
    parsedResult = message.content ? JSON.parse(message.content) : null;
  } catch {
    parsedResult = message.content;
  }

  const copyToClipboard = (text: string, e: React.MouseEvent) => {
    e.stopPropagation();
    navigator.clipboard.writeText(text).then(() => {
      showToast('复制成功');
    }).catch(() => {
      showToast('复制失败');
    });
  };

  const argsStr = parsedArgs
    ? typeof parsedArgs === 'string'
      ? parsedArgs
      : JSON.stringify(parsedArgs, null, 2)
    : '';

  const resultStr = parsedResult
    ? typeof parsedResult === 'string'
      ? parsedResult
      : JSON.stringify(parsedResult, null, 2)
    : '';

  return (
    <div className="message-item tool-message">
      <div className="message-content">
        <div className={`tool-call-card ${isExpanded ? 'expanded' : 'collapsed'}`}>
          <div className="tool-call-header" onClick={() => setIsExpanded(!isExpanded)}>
            <span className="tool-icon">{getToolIcon(message.toolCallName)}</span>
            <span className="tool-name">{message.toolCallName}</span>
            <span className={`tool-status-badge ${statusClass[message.toolCallStatus]}`}>
              {message.toolCallStatus === ToolCallStatusEnum.Success && (
                <svg className="status-check" viewBox="0 0 12 12" fill="none">
                  <path d="M10 3L4.5 8.5L2 6" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
                </svg>
              )}
              {statusText[message.toolCallStatus]}
            </span>
            <span className="expand-icon">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2">
                <polyline points={isExpanded ? "18 15 12 9 6 15" : "6 9 12 15 18 9"} />
              </svg>
            </span>
          </div>

          {isExpanded && (
            <>
              {parsedArgs && (
                <div className="tool-section">
                  <div className="tool-section-title">PARAMETERS</div>
                  <div className="tool-code-wrapper">
                    <pre className="tool-code">{argsStr}</pre>
                    <button className="copy-btn" onClick={(e) => copyToClipboard(argsStr, e)} title="复制">
                      <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" strokeWidth="2">
                        <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                        <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                      </svg>
                    </button>
                  </div>
                </div>
              )}

              {parsedResult && (
                <div className="tool-section">
                  <div className="tool-section-title">RESULT</div>
                  <div className="tool-code-wrapper result-wrapper">
                    <pre className="tool-code result-code">{resultStr}</pre>
                    <button className="copy-btn" onClick={(e) => copyToClipboard(resultStr, e)} title="复制">
                      <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" strokeWidth="2">
                        <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                        <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                      </svg>
                    </button>
                  </div>
                </div>
              )}
            </>
          )}

          {message.toolCallStatus === ToolCallStatusEnum.Processing && (
            <div className="tool-loading">
              <div className="tool-loading-bar" />
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

function formatLine(line: string): React.ReactNode {
  const patterns = [
    { regex: /\[([^\]]+)\]\((https?:\/\/[^\s)]+)\)/g, type: 'mdlink' },
    { regex: /\*\*(.*?)\*\*/g, type: 'bold' },
    { regex: /(https?:\/\/[^\s<>[\]()\]]+)/g, type: 'url' },
  ];

  type Token = { index: number; end: number; type: string; text: string; url?: string };
  const tokens: Token[] = [];

  for (const { regex, type } of patterns) {
    let match;
    regex.lastIndex = 0;
    while ((match = regex.exec(line)) !== null) {
      const overlaps = tokens.some(
        (t) => (match!.index >= t.index && match!.index < t.end) || (t.index >= match!.index && t.index < match!.index + match![0].length)
      );
      if (!overlaps) {
        if (type === 'mdlink') {
          tokens.push({ index: match.index, end: match.index + match[0].length, type, text: match[1], url: match[2] });
        } else if (type === 'url') {
          tokens.push({ index: match.index, end: match.index + match[0].length, type, text: match[1], url: match[1] });
        } else {
          tokens.push({ index: match.index, end: match.index + match[0].length, type, text: match[1] });
        }
      }
    }
  }

  tokens.sort((a, b) => a.index - b.index);

  const parts: React.ReactNode[] = [];
  let lastIndex = 0;

  for (const token of tokens) {
    if (token.index > lastIndex) {
      parts.push(line.slice(lastIndex, token.index));
    }
    if (token.type === 'bold') {
      parts.push(<strong key={token.index}>{token.text}</strong>);
    } else if (token.type === 'mdlink' || token.type === 'url') {
      parts.push(
        <a key={token.index} href={token.url} target="_blank" rel="noopener noreferrer" className="inline-link">
          {token.text}
        </a>
      );
    }
    lastIndex = token.end;
  }

  if (lastIndex < line.length) {
    parts.push(line.slice(lastIndex));
  }

  return parts.length > 0 ? parts : line || '\u00A0';
}

export function MessageItem({ message }: MessageItemProps) {
  switch (message.role) {
    case MessageRoleEnum.USER:
      return <UserMessageContent message={message} />;
    case MessageRoleEnum.ASSISTANT:
      return <AssistantMessageContent message={message as AssistantMessage} />;
    case MessageRoleEnum.TOOL:
      return <ToolMessageContent key={`${message.id}:${(message as ToolMessage).toolCallStatus}`} message={message as ToolMessage} />;
    default:
      return null;
  }
}
