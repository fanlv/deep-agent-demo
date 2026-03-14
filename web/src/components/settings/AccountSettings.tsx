import { useState } from 'react';

interface UserProfile {
  username: string;
  email: string;
  avatar: string;
}

export function AccountSettings() {
  const [profile, setProfile] = useState<UserProfile>({
    username: 'User',
    email: 'user@example.com',
    avatar: '👤',
  });

  const [notifications, setNotifications] = useState({
    email: true,
    push: false,
  });

  const handleProfileChange = (field: keyof UserProfile, value: string) => {
    setProfile((prev) => ({ ...prev, [field]: value }));
  };

  const handleSaveProfile = () => {
    console.log('Saving profile:', profile);
  };

  return (
    <div className="account-settings">
      <section className="settings-section">
        <h3 className="settings-section-title">个人信息</h3>

        <div className="settings-form-group">
          <label className="settings-label">用户名</label>
          <input
            type="text"
            className="settings-input"
            value={profile.username}
            onChange={(e) => handleProfileChange('username', e.target.value)}
            placeholder="请输入用户名"
          />
        </div>

        <div className="settings-form-group">
          <label className="settings-label">邮箱</label>
          <input
            type="email"
            className="settings-input"
            value={profile.email}
            onChange={(e) => handleProfileChange('email', e.target.value)}
            placeholder="请输入邮箱"
          />
        </div>

        <div className="settings-btn-group">
          <button className="settings-btn settings-btn-primary" onClick={handleSaveProfile}>
            保存更改
          </button>
        </div>
      </section>

      <div className="settings-divider" />

      <section className="settings-section">
        <h3 className="settings-section-title">通知设置</h3>

        <div className="settings-switch">
          <div className="settings-switch-label">
            <span className="settings-switch-title">邮件通知</span>
            <span className="settings-switch-desc">接收重要更新和消息的邮件通知</span>
          </div>
          <div
            className={`settings-toggle ${notifications.email ? 'active' : ''}`}
            onClick={() => setNotifications((prev) => ({ ...prev, email: !prev.email }))}
          />
        </div>

        <div className="settings-switch">
          <div className="settings-switch-label">
            <span className="settings-switch-title">推送通知</span>
            <span className="settings-switch-desc">接收浏览器推送通知</span>
          </div>
          <div
            className={`settings-toggle ${notifications.push ? 'active' : ''}`}
            onClick={() => setNotifications((prev) => ({ ...prev, push: !prev.push }))}
          />
        </div>
      </section>

      <div className="settings-divider" />

      <section className="settings-section">
        <h3 className="settings-section-title">账号安全</h3>

        <div className="settings-form-group">
          <label className="settings-label">修改密码</label>
          <input
            type="password"
            className="settings-input"
            placeholder="请输入当前密码"
          />
        </div>

        <div className="settings-form-group">
          <input
            type="password"
            className="settings-input"
            placeholder="请输入新密码"
          />
        </div>

        <div className="settings-form-group">
          <input
            type="password"
            className="settings-input"
            placeholder="请确认新密码"
          />
        </div>

        <div className="settings-btn-group">
          <button className="settings-btn settings-btn-secondary">
            更新密码
          </button>
        </div>
      </section>

      <div className="settings-divider" />

      <section className="settings-section">
        <h3 className="settings-section-title">危险操作</h3>
        <p className="settings-section-desc">以下操作不可逆，请谨慎操作</p>

        <div className="settings-btn-group">
          <button className="settings-btn settings-btn-danger">
            注销账号
          </button>
        </div>
      </section>
    </div>
  );
}
