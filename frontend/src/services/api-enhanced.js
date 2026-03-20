import axios from 'axios';

/**
 * API Client for GOYDA
 * Configured with error handling, auth, and base URL
 */

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8000/api';

console.log('🔗 API Base URL:', API_BASE_URL);

// Create Axios instance
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

/**
 * Request interceptor - add JWT token if available
 */
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

/**
 * Response interceptor - handle token refresh
 */
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = localStorage.getItem('refresh_token');
        if (refreshToken) {
          const response = await apiClient.post('/auth/token/refresh/', {
            refresh: refreshToken,
          });

          localStorage.setItem('access_token', response.data.access);
          originalRequest.headers.Authorization = `Bearer ${response.data.access}`;

          return apiClient(originalRequest);
        }
      } catch (refreshError) {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        window.location.href = '/login';
      }
    }

    return Promise.reject(error);
  }
);

/**
 * Locations API
 */
export const locationsAPI = {
  getAll: (params) => apiClient.get('/locations/', { params }),
  getById: (id) => apiClient.get(`/locations/${id}/`),
  search: (query) => apiClient.get('/locations/', { params: { search: query } }),
  filterByCategory: (category) =>
    apiClient.get('/locations/', { params: { category } }),
  filterBySeason: (season) =>
    apiClient.get('/locations/', { params: { season } }),
};

/**
 * Routes API
 */
export const routesAPI = {
  getAll: (params) => apiClient.get('/routes/', { params }),
  getById: (id) => apiClient.get(`/routes/${id}/`),
  create: (data) => apiClient.post('/routes/', data),
  update: (id, data) => apiClient.put(`/routes/${id}/`, data),
  delete: (id) => apiClient.delete(`/routes/${id}/`),
  getFeatured: () => apiClient.get('/routes/?featured=true'),
};

/**
 * AI Recommendations API
 */
export const aiAPI = {
  generateRoute: (data) => apiClient.post('/ai/generate-route/', data),
  getRecommendations: () => apiClient.get('/ai/recommendations/'),
};

/**
 * Authentication API
 */
export const authAPI = {
  login: (username, password) =>
    apiClient.post('/auth/token/', { username, password }),
  refreshToken: (refresh) =>
    apiClient.post('/auth/token/refresh/', { refresh }),
  register: (data) => apiClient.post('/auth/register/', data),
};

/**
 * User API
 */
export const userAPI = {
  getProfile: () => apiClient.get('/users/profile/'),
  updateProfile: (data) => apiClient.put('/users/profile/', data),
  getPreferences: () => apiClient.get('/users/preferences/'),
  updatePreferences: (data) => apiClient.put('/users/preferences/', data),
};

export default apiClient;
