import { useState } from 'react';
import { AccountSettings } from './AccountSettings';
import { ModelSettings } from './ModelSettings';
import { PromptSettings } from './PromptSettings';
import './Settings.css';

type SettingsTab = 'account' | 'model' | 'prompt';

interface SettingsProps {
  onClose: () => void;
}

const tabs: { key: SettingsTab; label: string; icon: string }[] = [
  { key: 'account', label: '账号设置', icon: '👤' },
  { key: 'model', label: '模型管理', icon: '🤖' },
  { key: 'prompt', label: 'Prompt 设置', icon: '📝' },
];

export function Settings({ onClose }: SettingsProps) {
  const [activeTab, setActiveTab] = useState<SettingsTab>('account');

  const renderContent = () => {
    switch (activeTab) {
      case 'account':
        return <AccountSettings />;
      case 'model':
        return <ModelSettings />;
      case 'prompt':
        return <PromptSettings />;
      default:
        return null;
    }
  };

  return (
    <div className="settings-overlay" onClick={onClose}>
      <div className="settings-modal" onClick={(e) => e.stopPropagation()}>
        <div className="settings-header">
          <h2 className="settings-title">设置</h2>
          <button className="settings-close-btn" onClick={onClose}>
            ×
          </button>
        </div>

        <div className="settings-body">
          <nav className="settings-nav">
            {tabs.map((tab) => (
              <div
                key={tab.key}
                className={`settings-nav-item ${activeTab === tab.key ? 'active' : ''}`}
                onClick={() => setActiveTab(tab.key)}
              >
                <span className="settings-nav-icon">{tab.icon}</span>
                <span className="settings-nav-label">{tab.label}</span>
              </div>
            ))}
          </nav>

          <div className="settings-content">{renderContent()}</div>
        </div>
      </div>
    </div>
  );
}
