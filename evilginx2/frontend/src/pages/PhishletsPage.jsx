import React, { useState, useEffect } from 'react';
import { phishletAPI } from '../api/client';
import '../styles/phishlets.css';

export default function PhishletsPage() {
  const [phishlets, setPhishlets] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchPhishlets();
  }, []);

  const fetchPhishlets = async () => {
    try {
      setLoading(true);
      const response = await phishletAPI.getAll();
      setPhishlets(response.data.data || []);
    } catch (err) {
      setError('Failed to load phishlets');
    } finally {
      setLoading(false);
    }
  };

  const handleToggle = async (name, isEnabled) => {
    try {
      if (isEnabled) {
        await phishletAPI.disable(name);
      } else {
        await phishletAPI.enable(name);
      }
      fetchPhishlets();
    } catch (err) {
      setError(`Failed to ${isEnabled ? 'disable' : 'enable'} phishlet`);
    }
  };

  if (loading) return <div className="loading">Loading phishlets...</div>;

  return (
    <div className="phishlets-page">
      <h2>Phishlets</h2>
      
      {error && <div className="error-message">{error}</div>}
      
      <button onClick={fetchPhishlets} className="btn-primary">Refresh</button>

      <div className="phishlets-grid">
        {phishlets.map(phishlet => (
          <div key={phishlet.name} className="phishlet-card">
            <h3>{phishlet.name}</h3>
            <p className="version">Version: {phishlet.version}</p>
            <button
              onClick={() => handleToggle(phishlet.name, phishlet.enabled)}
              className={phishlet.enabled ? 'btn-danger' : 'btn-primary'}
            >
              {phishlet.enabled ? 'Disable' : 'Enable'}
            </button>
            <div className={`status ${phishlet.enabled ? 'enabled' : 'disabled'}`}>
              {phishlet.enabled ? '● Active' : '○ Inactive'}
            </div>
          </div>
        ))}
      </div>

      {phishlets.length === 0 && (
        <div className="empty-state">No phishlets loaded</div>
      )}
    </div>
  );
}
