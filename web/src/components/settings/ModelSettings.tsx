import { useState, useEffect, useCallback } from 'react';
import './ModelSettings.css';

interface ConnectionInfo {
  api_key?: string;
  base_url?: string;
  model: string;
  ark?: { region?: string };
  openai?: { by_azure?: boolean; api_version?: string };
  gemini?: { backend?: string; project?: string; location?: string };
}

interface ModelInstance {
  id: number;
  model_class: string;
  display_name: string;
  connection: ConnectionInfo;
  thinking_type?: string;
  enable_base64_url?: boolean;
  status: number;
  created_at: number;
}

interface ProviderInfo {
  model_class: string;
  name: string;
  description: string;
  icon_url: string;
}

interface ProviderModelList {
  provider: ProviderInfo;
  model_list: ModelInstance[];
}

const API_BASE = '/api/v1/config/model';

const PROVIDER_ICONS: Record<string, string> = {
  ark: 'https://upload-images.jianshu.io/upload_images/12321605-0ece441a9983a40d.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240',
  openai: 'https://upload-images.jianshu.io/upload_images/12321605-91a8106e59f7126f.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240',
  claude: 'https://upload-images.jianshu.io/upload_images/12321605-2fc28d63c089a216.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240',
  deepseek: 'https://upload-images.jianshu.io/upload_images/12321605-6a3bdc5e184a6e04.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240',
  gemini: 'https://upload-images.jianshu.io/upload_images/12321605-21f811ad1bed58bd.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240',
  ollama: 'https://upload-images.jianshu.io/upload_images/12321605-ee4bd5afa8598a64.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240',
  qwen: 'https://upload-images.jianshu.io/upload_images/12321605-2763958be48a880a.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240',
};

const PROVIDER_ICONS_FALLBACK: Record<string, string> = {
  ark: '🔥',
  openai: '🤖',
  claude: '🎭',
  deepseek: '🔍',
  gemini: '✨',
  ollama: '🦙',
  qwen: '☁️',
};

const THINKING_TYPE_OPTIONS = [
  { value: '', label: '默认' },
  { value: 'enable', label: '启用' },
  { value: 'disable', label: '禁用' },
  { value: 'auto', label: '自动' },
];

const SUPPORTS_THINKING = ['ark', 'claude', 'gemini', 'qwen', 'ollama'];

const PROVIDER_PLACEHOLDERS: Record<string, {
  displayName: string;
  model: string;
  apiKey: string;
}> = {
  ark: {
    displayName: '例如：Doubao Seed 1.6',
    model: '例如：doubao-seed-1-6-250615',
    apiKey: '密钥',
  },
  openai: {
    displayName: '例如：GPT-4o',
    model: '例如：gpt-4o',
    apiKey: '密钥',
  },
  claude: {
    displayName: '例如：Claude 3.5 Sonnet',
    model: '例如：claude-3-5-sonnet-20241022',
    apiKey: '密钥',
  },
  deepseek: {
    displayName: '例如：DeepSeek V3',
    model: '例如：deepseek-chat',
    apiKey: '密钥',
  },
  gemini: {
    displayName: '例如：Gemini 2.0 Flash',
    model: '例如：gemini-2.0-flash-exp',
    apiKey: '密钥',
  },
  ollama: {
    displayName: '例如：Qwen 2.5 7B',
    model: '例如：qwen2.5:7b',
    apiKey: '',
  },
  qwen: {
    displayName: '例如：Qwen Max',
    model: '例如：qwen-max',
    apiKey: '密钥',
  },
};

function maskApiKey(key?: string): string {
  if (!key) return '未配置';
  if (key.length <= 10) return '***';
  return key.slice(0, 6) + '***' + key.slice(-4);
}

export function ModelSettings() {
  const [providers, setProviders] = useState<ProviderModelList[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [addingProvider, setAddingProvider] = useState<ProviderInfo | null>(null);
  const [formData, setFormData] = useState({
    displayName: '',
    model: '',
    apiKey: '',
    baseUrl: '',
    enableBase64Url: false,
    thinkingType: '',
    arkRegion: '',
    openaiByAzure: false,
    openaiApiVersion: '',
    geminiBackend: '',
    geminiProject: '',
    geminiLocation: '',
  });
  const [saving, setSaving] = useState(false);

  const loadModels = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const resp = await fetch(`${API_BASE}/list`);
      const data = await resp.json();
      if (data.code !== 0) {
        throw new Error(data.msg || '加载失败');
      }
      setProviders(data.provider_model_list || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载模型列表失败');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadModels();
  }, [loadModels]);

  const handleDelete = async (id: number) => {
    if (!confirm('确认删除该模型？删除后不可恢复')) return;
    try {
      const resp = await fetch(`${API_BASE}/delete`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: String(id) }),
      });
      const data = await resp.json();
      if (data.code !== 0) {
        throw new Error(data.msg || '删除失败');
      }
      loadModels();
    } catch (err) {
      alert(err instanceof Error ? err.message : '删除失败');
    }
  };

  const handleOpenAddModal = (provider: ProviderInfo) => {
    setAddingProvider(provider);
    setFormData({
      displayName: '',
      model: '',
      apiKey: '',
      baseUrl: '',
      enableBase64Url: false,
      thinkingType: '',
      arkRegion: '',
      openaiByAzure: false,
      openaiApiVersion: '',
      geminiBackend: '',
      geminiProject: '',
      geminiLocation: '',
    });
  };

  const handleCloseModal = () => {
    setAddingProvider(null);
  };

  const handleSave = async () => {
    if (!addingProvider) return;
    if (!formData.displayName || !formData.model) {
      alert('请填写显示名称和模型名称');
      return;
    }
    if (addingProvider.model_class !== 'ollama' && !formData.apiKey) {
      alert('请填写 API Key');
      return;
    }

    setSaving(true);
    try {
      const connection: Record<string, unknown> = {
        model: formData.model,
      };
      if (formData.apiKey) connection.api_key = formData.apiKey;
      if (formData.baseUrl) connection.base_url = formData.baseUrl;

      if (addingProvider.model_class === 'ark' && formData.arkRegion) {
        connection.ark = { region: formData.arkRegion };
      }
      if (addingProvider.model_class === 'openai') {
        if (formData.openaiByAzure || formData.openaiApiVersion) {
          connection.openai = {
            by_azure: formData.openaiByAzure,
            api_version: formData.openaiApiVersion || undefined,
          };
        }
      }
      if (addingProvider.model_class === 'gemini') {
        if (formData.geminiBackend || formData.geminiProject || formData.geminiLocation) {
          connection.gemini = {
            backend: formData.geminiBackend || undefined,
            project: formData.geminiProject || undefined,
            location: formData.geminiLocation || undefined,
          };
        }
      }

      const payload = {
        model_class: addingProvider.model_class,
        display_name: formData.displayName,
        connection,
        thinking_type: formData.thinkingType || undefined,
        enable_base64_url: formData.enableBase64Url || undefined,
      };

      const resp = await fetch(`${API_BASE}/create`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });
      const data = await resp.json();
      if (data.code !== 0) {
        throw new Error(data.msg || '创建失败');
      }
      handleCloseModal();
      loadModels();
    } catch (err) {
      alert(err instanceof Error ? err.message : '创建失败');
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="model-settings">
        <div className="model-loading">加载中...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="model-settings">
        <div className="model-error">
          <p>{error}</p>
          <button className="settings-btn settings-btn-primary" onClick={loadModels}>
            重试
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="model-settings">
      <section className="settings-section">
        <div className="model-settings-header">
          <h3 className="settings-section-title">模型列表</h3>
          <span className="model-provider-count">共 {providers.length} 个提供商</span>
        </div>

        <div className="provider-list">
          {providers.map((pm) => (
            <div key={pm.provider.model_class} className="provider-card">
              <div className="provider-header">
                <div className="provider-info">
                  <span className="provider-icon">
                    <img 
                      src={PROVIDER_ICONS[pm.provider.model_class]} 
                      alt={pm.provider.name}
                      referrerPolicy="no-referrer"
                      style={{ width: '32px', height: '32px', objectFit: 'contain' }}
                      onError={(e) => {
                        e.currentTarget.style.display = 'none';
                        const next = e.currentTarget.nextElementSibling as HTMLSpanElement | null;
                        if (next) {
                          next.style.display = 'inline';
                        }
                      }}
                    />
                    <span style={{ display: 'none', fontSize: '32px', lineHeight: '32px' }}>
                      {PROVIDER_ICONS_FALLBACK[pm.provider.model_class] || '🤖'}
                    </span>
                  </span>
                  <div className="provider-meta">
                    <span className="provider-name">{pm.provider.name}</span>
                    <span className="provider-desc">{pm.provider.description}</span>
                  </div>
                </div>
                <button
                  className="settings-btn settings-btn-primary"
                  onClick={() => handleOpenAddModal(pm.provider)}
                >
                  + 添加模型
                </button>
              </div>

              <div className="provider-body">
                <div className="provider-model-count">
                  已配置 {pm.model_list.length} 个模型
                </div>
                {pm.model_list.map((model) => (
                  <div key={model.id} className="model-card">
                    <div className="model-card-header">
                      <span className="model-name">{model.display_name}</span>
                      <div className="model-actions">
                        <span
                          className={`model-status ${model.status === 1 ? 'active' : ''}`}
                        >
                          {model.status === 1 ? '启用' : '禁用'}
                        </span>
                        <button
                          className="model-action-btn model-action-btn-danger"
                          onClick={() => handleDelete(model.id)}
                        >
                          删除
                        </button>
                      </div>
                    </div>
                    <div className="model-card-details">
                      <div className="model-detail-row">
                        <span className="detail-label">ID:</span>
                        <span className="detail-value">{model.id}</span>
                      </div>
                      <div className="model-detail-row">
                        <span className="detail-label">模型:</span>
                        <span className="detail-value">{model.connection.model}</span>
                      </div>
                      {pm.provider.model_class !== 'ollama' && (
                        <div className="model-detail-row">
                          <span className="detail-label">API Key:</span>
                          <span className="detail-value">
                            {maskApiKey(model.connection.api_key)}
                          </span>
                        </div>
                      )}
                      {model.connection.base_url && (
                        <div className="model-detail-row">
                          <span className="detail-label">Endpoint:</span>
                          <span className="detail-value detail-url">
                            {model.connection.base_url}
                          </span>
                        </div>
                      )}
                      {model.enable_base64_url && (
                        <div className="model-detail-row">
                          <span className="detail-label">Base64 URL:</span>
                          <span className="detail-value">启用</span>
                        </div>
                      )}
                      {SUPPORTS_THINKING.includes(pm.provider.model_class) &&
                        model.thinking_type && (
                          <div className="model-detail-row">
                            <span className="detail-label">思考模式:</span>
                            <span className="detail-value">{model.thinking_type}</span>
                          </div>
                        )}
                      {pm.provider.model_class === 'ark' && model.connection.ark?.region && (
                        <div className="model-detail-row">
                          <span className="detail-label">Region:</span>
                          <span className="detail-value">{model.connection.ark.region}</span>
                        </div>
                      )}
                      {pm.provider.model_class === 'openai' && model.connection.openai && (
                        <>
                          <div className="model-detail-row">
                            <span className="detail-label">Azure:</span>
                            <span className="detail-value">
                              {model.connection.openai.by_azure ? '是' : '否'}
                            </span>
                          </div>
                          {model.connection.openai.api_version && (
                            <div className="model-detail-row">
                              <span className="detail-label">API Version:</span>
                              <span className="detail-value">
                                {model.connection.openai.api_version}
                              </span>
                            </div>
                          )}
                        </>
                      )}
                      {pm.provider.model_class === 'gemini' && model.connection.gemini && (
                        <>
                          {model.connection.gemini.backend && (
                            <div className="model-detail-row">
                              <span className="detail-label">Backend:</span>
                              <span className="detail-value">
                                {model.connection.gemini.backend}
                              </span>
                            </div>
                          )}
                          {model.connection.gemini.project && (
                            <div className="model-detail-row">
                              <span className="detail-label">Project:</span>
                              <span className="detail-value">
                                {model.connection.gemini.project}
                              </span>
                            </div>
                          )}
                          {model.connection.gemini.location && (
                            <div className="model-detail-row">
                              <span className="detail-label">Location:</span>
                              <span className="detail-value">
                                {model.connection.gemini.location}
                              </span>
                            </div>
                          )}
                        </>
                      )}
                    </div>
                  </div>
                ))}
                {pm.model_list.length === 0 && (
                  <div className="no-models">暂无配置的模型</div>
                )}
              </div>
            </div>
          ))}
        </div>
      </section>

      {addingProvider && (
        <div className="modal-overlay" onClick={handleCloseModal}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <span className="modal-title">添加 {addingProvider.name} 模型</span>
              <button className="modal-close" onClick={handleCloseModal}>
                ×
              </button>
            </div>
            <div className="modal-body">
              <div className="form-group">
                <label className="form-label">
                  显示名称 <span className="required">*</span>
                </label>
                <input
                  type="text"
                  className="form-input"
                  value={formData.displayName}
                  onChange={(e) =>
                    setFormData({ ...formData, displayName: e.target.value })
                  }
                  placeholder={
                    PROVIDER_PLACEHOLDERS[addingProvider.model_class]?.displayName ||
                    '例如：GPT-4o'
                  }
                />
              </div>
              <div className="form-group">
                <label className="form-label">
                  模型名称 <span className="required">*</span>
                </label>
                <input
                  type="text"
                  className="form-input"
                  value={formData.model}
                  onChange={(e) => setFormData({ ...formData, model: e.target.value })}
                  placeholder={
                    PROVIDER_PLACEHOLDERS[addingProvider.model_class]?.model ||
                    '例如：gpt-4o'
                  }
                />
              </div>
              {addingProvider.model_class !== 'ollama' && (
                <div className="form-group">
                  <label className="form-label">
                    API Key <span className="required">*</span>
                  </label>
                  <input
                    type="password"
                    className="form-input"
                    value={formData.apiKey}
                    onChange={(e) => setFormData({ ...formData, apiKey: e.target.value })}
                    placeholder={
                      PROVIDER_PLACEHOLDERS[addingProvider.model_class]?.apiKey || '密钥'
                    }
                  />
                </div>
              )}
              <div className="form-group">
                <label className="form-label">Base URL</label>
                <input
                  type="text"
                  className="form-input"
                  value={formData.baseUrl}
                  onChange={(e) => setFormData({ ...formData, baseUrl: e.target.value })}
                  placeholder="选填"
                />
              </div>
              <div className="form-group form-group-checkbox">
                <input
                  type="checkbox"
                  id="enableBase64Url"
                  checked={formData.enableBase64Url}
                  onChange={(e) =>
                    setFormData({ ...formData, enableBase64Url: e.target.checked })
                  }
                />
                <label htmlFor="enableBase64Url">启用 Base64 URL</label>
              </div>
              {SUPPORTS_THINKING.includes(addingProvider.model_class) && (
                <div className="form-group">
                  <label className="form-label">思考模式</label>
                  <select
                    className="form-select"
                    value={formData.thinkingType}
                    onChange={(e) =>
                      setFormData({ ...formData, thinkingType: e.target.value })
                    }
                  >
                    {THINKING_TYPE_OPTIONS.map((opt) => (
                      <option key={opt.value} value={opt.value}>
                        {opt.label}
                      </option>
                    ))}
                  </select>
                </div>
              )}
              {addingProvider.model_class === 'ark' && (
                <div className="form-group">
                  <label className="form-label">Region</label>
                  <input
                    type="text"
                    className="form-input"
                    value={formData.arkRegion}
                    onChange={(e) =>
                      setFormData({ ...formData, arkRegion: e.target.value })
                    }
                    placeholder="选填，例如：cn-beijing"
                  />
                </div>
              )}
              {addingProvider.model_class === 'openai' && (
                <>
                  <div className="form-group form-group-checkbox">
                    <input
                      type="checkbox"
                      id="openaiByAzure"
                      checked={formData.openaiByAzure}
                      onChange={(e) =>
                        setFormData({ ...formData, openaiByAzure: e.target.checked })
                      }
                    />
                    <label htmlFor="openaiByAzure">使用 Azure OpenAI</label>
                  </div>
                  <div className="form-group">
                    <label className="form-label">API Version</label>
                    <input
                      type="text"
                      className="form-input"
                      value={formData.openaiApiVersion}
                      onChange={(e) =>
                        setFormData({ ...formData, openaiApiVersion: e.target.value })
                      }
                      placeholder="选填，例如：2024-06-01"
                    />
                  </div>
                </>
              )}
              {addingProvider.model_class === 'gemini' && (
                <>
                  <div className="form-group">
                    <label className="form-label">Backend</label>
                    <select
                      className="form-select"
                      value={formData.geminiBackend}
                      onChange={(e) =>
                        setFormData({ ...formData, geminiBackend: e.target.value })
                      }
                    >
                      <option value="">默认 (Gemini API)</option>
                      <option value="gemini">Gemini API</option>
                      <option value="vertex">Vertex AI</option>
                    </select>
                  </div>
                  <div className="form-group">
                    <label className="form-label">Project</label>
                    <input
                      type="text"
                      className="form-input"
                      value={formData.geminiProject}
                      onChange={(e) =>
                        setFormData({ ...formData, geminiProject: e.target.value })
                      }
                      placeholder="选填，GCP 项目 ID"
                    />
                  </div>
                  <div className="form-group">
                    <label className="form-label">Location</label>
                    <input
                      type="text"
                      className="form-input"
                      value={formData.geminiLocation}
                      onChange={(e) =>
                        setFormData({ ...formData, geminiLocation: e.target.value })
                      }
                      placeholder="选填，例如：us-central1"
                    />
                  </div>
                </>
              )}
            </div>
            <div className="modal-footer">
              <button
                className="settings-btn settings-btn-secondary"
                onClick={handleCloseModal}
              >
                取消
              </button>
              <button
                className="settings-btn settings-btn-primary"
                onClick={handleSave}
                disabled={saving}
              >
                {saving ? '保存中...' : '保存'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
