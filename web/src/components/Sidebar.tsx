import { useEffect, useState, useCallback } from 'react';
import './Sidebar.css';

const SESSIONS_POLL_INTERVAL_MS = 60_000;

interface SessionInfo {
  id: string;
  title: string;
  createdAt: number;
  updatedAt: number;
}

interface SidebarProps {
  currentSessionId?: string;
  onNewChat: () => void;
  onSelectSession: (sessionId: string) => void;
  onOpenSettings: () => void;
  baseUrl: string;
}

export function Sidebar({ currentSessionId, onNewChat, onSelectSession, onOpenSettings, baseUrl }: SidebarProps) {
  const [sessions, setSessions] = useState<SessionInfo[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const fetchSessions = useCallback(async () => {
    setIsLoading(true);
    try {
      const response = await fetch(baseUrl);
      if (response.ok) {
        const data = await response.json();
        const sorted = (data.sessions || []).sort(
          (a: SessionInfo, b: SessionInfo) => b.updatedAt - a.updatedAt
        );
        setSessions(sorted);
      }
    } catch (error) {
      console.error('Failed to fetch sessions:', error);
    } finally {
      setIsLoading(false);
    }
  }, [baseUrl]);

  useEffect(() => {
    fetchSessions();
    const interval = setInterval(fetchSessions, SESSIONS_POLL_INTERVAL_MS);
    return () => clearInterval(interval);
  }, [fetchSessions]);

  const handleDelete = async (e: React.MouseEvent, sessionId: string) => {
    e.stopPropagation();
    try {
      const response = await fetch(`${baseUrl}/${sessionId}`, {
        method: 'DELETE',
      });
      if (response.ok) {
        setSessions((prev) => prev.filter((t) => t.id !== sessionId));
        if (currentSessionId === sessionId) {
          onNewChat();
        }
      }
    } catch (error) {
      console.error('Failed to delete session:', error);
    }
  };

  return (
    <aside className="sidebar">
      <div className="sidebar-header">
        <div className="sidebar-logo">
          <span className="logo-icon">🤖</span>
          <span className="logo-text">Deep Agent</span>
        </div>
      </div>

      <div className="sidebar-new-chat" onClick={onNewChat}>
        <span className="new-chat-icon">+</span>
        <span>New Chat</span>
      </div>

      <div className="sidebar-section">
        <div className="sidebar-section-title">Chat History</div>
        <div className="sidebar-sessions">
          {isLoading && sessions.length === 0 ? (
            <div className="sidebar-loading">Loading...</div>
          ) : sessions.length === 0 ? (
            <div className="sidebar-empty">No conversations yet</div>
          ) : (
            sessions.map((session) => (
              <div
                key={session.id}
                className={`sidebar-session-item ${currentSessionId === session.id ? 'active' : ''}`}
                onClick={() => onSelectSession(session.id)}
              >
                <span className="session-title">{session.title}</span>
                <button
                  className="session-delete-btn"
                  onClick={(e) => handleDelete(e, session.id)}
                  title="Delete"
                >
                  ×
                </button>
              </div>
            ))
          )}
        </div>
      </div>

      <div className="sidebar-footer">
        <div className="sidebar-settings" onClick={onOpenSettings}>
          <span className="settings-icon">⚙️</span>
          <span className="settings-text">设置</span>
        </div>
        <div className="sidebar-user">
          <span className="user-avatar">👤</span>
          <span className="user-name">FanLv</span>
        </div>
      </div>
    </aside>
  );
}
