import React, { useState, useEffect } from 'react';
import { configAPI } from '../api/client';
import '../styles/config.css';

export default function ConfigPage() {
  const [config, setConfig] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [editingKey, setEditingKey] = useState(null);
  const [editValue, setEditValue] = useState('');

  useEffect(() => {
    fetchConfig();
  }, []);

  const fetchConfig = async () => {
    try {
      setLoading(true);
      const response = await configAPI.get();
      setConfig(response.data.data || {});
    } catch (err) {
      setError('Failed to load configuration');
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = (key, value) => {
    setEditingKey(key);
    setEditValue(value);
  };

  const handleSave = async () => {
    try {
      await configAPI.update({ [editingKey]: editValue });
      setConfig(prev => ({ ...prev, [editingKey]: editValue }));
      setEditingKey(null);
      setError('');
    } catch (err) {
      setError('Failed to update configuration');
    }
  };

  if (loading) return <div className="loading">Loading configuration...</div>;

  return (
    <div className="config-page">
      <h2>Configuration</h2>
      
      {error && <div className="error-message">{error}</div>}
      
      <button onClick={fetchConfig} className="btn-primary">Refresh</button>

      <table className="config-table">
        <thead>
          <tr>
            <th>Key</th>
            <th>Value</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          {Object.entries(config).map(([key, value]) => (
            <tr key={key}>
              <td className="key">{key}</td>
              <td className="value">
                {editingKey === key ? (
                  <input
                    type="text"
                    value={editValue}
                    onChange={(e) => setEditValue(e.target.value)}
                    autoFocus
                  />
                ) : (
                  <span>{String(value)}</span>
                )}
              </td>
              <td className="action">
                {editingKey === key ? (
                  <>
                    <button onClick={handleSave} className="btn-small btn-primary">Save</button>
                    <button onClick={() => setEditingKey(null)} className="btn-small btn-secondary">Cancel</button>
                  </>
                ) : (
                  <button onClick={() => handleEdit(key, value)} className="btn-small btn-primary">Edit</button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {Object.keys(config).length === 0 && (
        <div className="empty-state">No configuration loaded</div>
      )}
    </div>
  );
}
