import React, { useState } from 'react';
import LoginPage from './pages/LoginPage';
import SessionsPage from './pages/SessionsPage';
import PhishletsPage from './pages/PhishletsPage';
import ConsolePage from './pages/ConsolePage';
import ConfigPage from './pages/ConfigPage';
import './styles/app.css';

export default function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(!!localStorage.getItem('auth_token'));
  const [currentPage, setCurrentPage] = useState('sessions');

  if (!isLoggedIn) {
    return <LoginPage onLogin={() => setIsLoggedIn(true)} />;
  }

  const handleLogout = () => {
    localStorage.removeItem('auth_token');
    setIsLoggedIn(false);
  };

  return (
    <div className="app-container">
      <nav className="navbar">
        <div className="logo">Evilginx2 Web UI</div>
        <div className="nav-links">
          <button 
            className={`nav-link ${currentPage === 'sessions' ? 'active' : ''}`}
            onClick={() => setCurrentPage('sessions')}
          >
            Sessions
          </button>
          <button 
            className={`nav-link ${currentPage === 'phishlets' ? 'active' : ''}`}
            onClick={() => setCurrentPage('phishlets')}
          >
            Phishlets
          </button>
          <button 
            className={`nav-link ${currentPage === 'config' ? 'active' : ''}`}
            onClick={() => setCurrentPage('config')}
          >
            Config
          </button>
          <button 
            className={`nav-link ${currentPage === 'console' ? 'active' : ''}`}
            onClick={() => setCurrentPage('console')}
          >
            Console
          </button>
          <button className="nav-link logout" onClick={handleLogout}>
            Logout
          </button>
        </div>
      </nav>

      <main className="main-content">
        {currentPage === 'sessions' && <SessionsPage />}
        {currentPage === 'phishlets' && <PhishletsPage />}
        {currentPage === 'config' && <ConfigPage />}
        {currentPage === 'console' && <ConsolePage />}
      </main>
    </div>
  );
}
