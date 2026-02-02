import apiClient from './http';

// Auth API
export const authAPI = {
    register: (data) => apiClient.post('/api/auth/register', data),
    login: (data) => apiClient.post('/api/auth/login', data),
};

// Product API
export const productAPI = {
    getAll: () => apiClient.get('/api/products'),
    getById: (id) => apiClient.get(`/api/products/${id}`),
    create: (data) => apiClient.post('/api/products', data),
    update: (id, data) => apiClient.put(`/api/products/${id}`, data),
    delete: (id) => apiClient.delete(`/api/products/${id}`),
};

// Order API
export const orderAPI = {
    getAll: () => apiClient.get('/api/orders'),
    getById: (id) => apiClient.get(`/api/orders/${id}`),
    create: (data) => apiClient.post('/api/orders', data),
};

// Payment API
export const paymentAPI = {
    create: (data) => apiClient.post('/api/payments', data),
    getByOrderId: (orderId) => apiClient.get(`/api/payments/order/${orderId}`),
};
