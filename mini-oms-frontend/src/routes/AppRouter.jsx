import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from '../context/AuthContext';
import PrivateRoute from './PrivateRoute';
import AdminRoute from './AdminRoute';
import Navbar from '../components/common/Navbar';

// Auth Pages
import LoginPage from '../pages/auth/LoginPage';
import RegisterPage from '../pages/auth/RegisterPage';

// Product Pages
import ProductListPage from '../pages/products/ProductListPage';
import ProductDetailPage from '../pages/products/ProductDetailPage';

// Order Pages
import OrderListPage from '../pages/orders/OrderListPage';
import OrderDetailPage from '../pages/orders/OrderDetailPage';
import CreateOrderPage from '../pages/orders/CreateOrderPage';

// Payment Pages
import CreatePaymentPage from '../pages/orders/CreatePaymentPage';

// Admin Pages
import AdminProductsPage from '../pages/admin/AdminProductsPage';
import AdminOrdersPage from '../pages/admin/AdminOrdersPage';

const AppRouter = () => {
    return (
        <BrowserRouter>
            <AuthProvider>
                <Navbar />
                <div className="container">
                    <Routes>
                        {/* Public Routes */}
                        <Route path="/login" element={<LoginPage />} />
                        <Route path="/register" element={<RegisterPage />} />

                        {/* Redirect root to products */}
                        <Route path="/" element={<Navigate to="/products" replace />} />

                        {/* Product Routes (Public) */}
                        <Route path="/products" element={<ProductListPage />} />
                        <Route path="/products/:id" element={<ProductDetailPage />} />

                        {/* Protected Routes */}
                        <Route
                            path="/orders"
                            element={
                                <PrivateRoute>
                                    <OrderListPage />
                                </PrivateRoute>
                            }
                        />
                        <Route
                            path="/orders/create"
                            element={
                                <PrivateRoute>
                                    <CreateOrderPage />
                                </PrivateRoute>
                            }
                        />
                        <Route
                            path="/orders/:id"
                            element={
                                <PrivateRoute>
                                    <OrderDetailPage />
                                </PrivateRoute>
                            }
                        />
                        <Route
                            path="/orders/:orderId/payment"
                            element={
                                <PrivateRoute>
                                    <CreatePaymentPage />
                                </PrivateRoute>
                            }
                        />

                        {/* Admin Routes */}
                        <Route
                            path="/admin/products"
                            element={
                                <AdminRoute>
                                    <AdminProductsPage />
                                </AdminRoute>
                            }
                        />
                        <Route
                            path="/admin/orders"
                            element={
                                <AdminRoute>
                                    <AdminOrdersPage />
                                </AdminRoute>
                            }
                        />
                    </Routes>
                </div>
            </AuthProvider>
        </BrowserRouter>
    );
};

export default AppRouter;
