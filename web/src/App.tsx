import { useState, useCallback, useEffect } from 'react';
import { ChatBot, HomePage, Sidebar, Settings } from './components';
import './App.css';

const API_BASE_URL = '/api/v1';

function getSessionIdFromUrl(): string | undefined {
  const params = new URLSearchParams(window.location.search);
  return params.get('sessionId') || undefined;
}

function updateUrlWithSessionId(sessionId: string) {
  const url = new URL(window.location.href);
  url.searchParams.set('sessionId', sessionId);
  window.history.pushState({}, '', url.toString());
}

function clearSessionIdFromUrl() {
  const url = new URL(window.location.href);
  url.searchParams.delete('sessionId');
  window.history.pushState({}, '', url.toString());
}

function App() {
  const [showChat, setShowChat] = useState(() => !!getSessionIdFromUrl());
  const [initialMessage, setInitialMessage] = useState<string | null>(null);
  const [currentSessionId, setCurrentSessionId] = useState<string | undefined>(() => getSessionIdFromUrl());
  const [isInitializing, setIsInitializing] = useState(false);
  const [showSettings, setShowSettings] = useState(false);

  useEffect(() => {
    const handlePopState = () => {
      const sessionId = getSessionIdFromUrl();
      setCurrentSessionId(sessionId);
      setShowChat(!!sessionId);
      if (!sessionId) {
        setInitialMessage(null);
      }
    };

    window.addEventListener('popstate', handlePopState);
    return () => window.removeEventListener('popstate', handlePopState);
  }, []);

  const handleStartChat = useCallback(async (message: string, modelId: number) => {
    setIsInitializing(true);
    try {
      const response = await fetch(`${API_BASE_URL}/agent/init`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ modelId }),
      });
      if (!response.ok) {
        const errData = await response.json().catch(() => null);
        const errMsg = errData?.error || `HTTP ${response.status}`;
        throw new Error(errMsg);
      }
      const data = await response.json();
      const sessionId = data.sessionId;

      localStorage.setItem(`session_model_${sessionId}`, String(modelId));
      updateUrlWithSessionId(sessionId);
      setCurrentSessionId(sessionId);
      setInitialMessage(message);
      setShowChat(true);
    } catch (err) {
      console.error('Failed to init agent:', err);
      const msg = err instanceof Error ? err.message : String(err);
      alert(`Failed to start chat: ${msg}`);
    } finally {
      setIsInitializing(false);
    }
  }, []);

  const handleNewChat = useCallback(() => {
    clearSessionIdFromUrl();
    setInitialMessage(null);
    setCurrentSessionId(undefined);
    setShowChat(false);
  }, []);

  const handleSelectSession = useCallback((sessionId: string) => {
    updateUrlWithSessionId(sessionId);
    setCurrentSessionId(sessionId);
    setInitialMessage(null);
    setShowChat(true);
  }, []);

  const handleSessionCreated = useCallback((sessionId: string) => {
    updateUrlWithSessionId(sessionId);
    setCurrentSessionId(sessionId);
  }, []);

  const handleOpenSettings = useCallback(() => {
    setShowSettings(true);
  }, []);

  const handleCloseSettings = useCallback(() => {
    setShowSettings(false);
  }, []);

  return (
    <div className="app-layout">
      {showChat && currentSessionId ? (
        <div className="app-main">
          <ChatBot
            baseUrl={`${API_BASE_URL}/agent`}
            initialMessage={initialMessage}
            existingSessionId={currentSessionId}
            onBack={handleNewChat}
            onSessionCreated={handleSessionCreated}
          />
        </div>
      ) : (
        <>
          <Sidebar
            currentSessionId={currentSessionId}
            onNewChat={handleNewChat}
            onSelectSession={handleSelectSession}
            onOpenSettings={handleOpenSettings}
            baseUrl={`${API_BASE_URL}/sessions`}
          />
          <div className="app-main">
            <HomePage onStartChat={handleStartChat} isInitializing={isInitializing} />
          </div>
        </>
      )}
      {showSettings && <Settings onClose={handleCloseSettings} />}
    </div>
  );
}

export default App;
