import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8000/api';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests if available
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const locationsAPI = {
  getAll: (params) => apiClient.get('/locations/', { params }),
  getById: (id) => apiClient.get(`/locations/${id}/`),
  search: (query) => apiClient.get('/locations/', { params: { search: query } }),
};

export const routesAPI = {
  getAll: (params) => apiClient.get('/routes/', { params }),
  getById: (id) => apiClient.get(`/routes/${id}/`),
  create: (data) => apiClient.post('/routes/', data),
  update: (id, data) => apiClient.put(`/routes/${id}/`, data),
};

export const aiAPI = {
  generateRoute: (data) => apiClient.post('/ai/generate-route/', data),
  getRecommendations: () => apiClient.get('/ai/recommendations/'),
};

export const authAPI = {
  login: (username, password) =>
    apiClient.post('/auth/token/', { username, password }),
  refreshToken: (refresh) =>
    apiClient.post('/auth/token/refresh/', { refresh }),
};

export default apiClient;
