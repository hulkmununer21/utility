import axios from 'axios';

const API_BASE = '/api';

// Create axios instance with auth token
const apiClient = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
  }
});

// Add token to all requests
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authAPI = {
  login: (username, password) => apiClient.post('/auth/login', { username, password }),
};

export const sessionAPI = {
  getAll: () => apiClient.get('/sessions'),
  getById: (id) => apiClient.get(`/sessions/${id}`),
  update: (id, field, value) => apiClient.put(`/sessions/${id}`, { field, value }),
  delete: (id) => apiClient.delete(`/sessions/${id}`),
  deleteAll: () => apiClient.delete('/sessions'),
};

export const phishletAPI = {
  getAll: () => apiClient.get('/phishlets'),
  enable: (name) => apiClient.post(`/phishlets/${name}/enable`),
  disable: (name) => apiClient.post(`/phishlets/${name}/disable`),
};

export const lureAPI = {
  getAll: () => apiClient.get('/lures'),
  create: (data) => apiClient.post('/lures', data),
  delete: (id) => apiClient.delete(`/lures/${id}`),
};

export const configAPI = {
  get: () => apiClient.get('/config'),
  update: (data) => apiClient.put('/config', data),
};

export const commandAPI = {
  execute: (command) => apiClient.post('/execute', { command }),
};

export default apiClient;
