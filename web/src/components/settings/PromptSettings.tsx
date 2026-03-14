import { useState, useEffect, useCallback } from 'react';
import './PromptSettings.css';

const API_BASE = '/api/v1/prompt';

type TabKey = 'system_prompt' | 'agents_md';

const TABS: { key: TabKey; label: string }[] = [
  { key: 'system_prompt', label: '系统提示词' },
  { key: 'agents_md', label: 'AGENTS.md' },
];

export function PromptSettings() {
  const [activeTab, setActiveTab] = useState<TabKey>('system_prompt');
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [savedTip, setSavedTip] = useState(false);

  const fetchPrompt = useCallback(async (key: TabKey) => {
    try {
      setLoading(true);
      const resp = await fetch(`${API_BASE}/get`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ key }),
      });
      const data = await resp.json();
      if (data.code !== 0) {
        throw new Error('加载失败');
      }
      setContent(data.prompt || '');
    } catch (e) {
      console.error('Failed to fetch prompt:', e);
      setContent('');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchPrompt(activeTab);
  }, [activeTab, fetchPrompt]);

  const handleSave = async () => {
    try {
      setSaving(true);
      const resp = await fetch(`${API_BASE}/save`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ key: activeTab, prompt: content }),
      });
      const data = await resp.json();
      if (data.code !== 0) {
        throw new Error('保存失败');
      }
      setSavedTip(true);
      setTimeout(() => setSavedTip(false), 2000);
    } catch (e) {
      console.error('Failed to save prompt:', e);
      alert('保存失败，请重试');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="prompt-settings">
      <div className="prompt-tabs">
        {TABS.map((tab) => (
          <button
            key={tab.key}
            className={`prompt-tab ${activeTab === tab.key ? 'prompt-tab-active' : ''}`}
            onClick={() => setActiveTab(tab.key)}
          >
            {tab.label}
          </button>
        ))}
      </div>

      <section className="settings-section">
        <p className="settings-section-desc">
          {activeTab === 'system_prompt'
            ? '设置默认的系统提示词，将应用于所有新对话'
            : '配置 AGENTS.md 内容，定义 Agent 的行为和规则'}
        </p>

        <div className="settings-form-group">
          {loading ? (
            <div className="prompt-loading">加载中...</div>
          ) : (
            <textarea
              className="settings-input settings-textarea prompt-textarea"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder={
                activeTab === 'system_prompt'
                  ? '输入系统提示词...'
                  : '输入 AGENTS.md 内容...'
              }
            />
          )}
        </div>

        <div className="settings-btn-group">
          <button
            className="settings-btn settings-btn-primary"
            onClick={handleSave}
            disabled={saving || loading}
          >
            {saving ? '保存中...' : '保存更改'}
          </button>
          {savedTip && <span className="prompt-saved-tip">已保存</span>}
        </div>
      </section>
    </div>
  );
}
