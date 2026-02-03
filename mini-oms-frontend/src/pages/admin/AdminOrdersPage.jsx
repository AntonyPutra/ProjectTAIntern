import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { orderAPI } from '../../api';
import '../../styles/admin.css';

const AdminOrdersPage = () => {
    const [orders, setOrders] = useState([]);
    const [stats, setStats] = useState({ total_orders: 0, total_revenue: 0, pending_payments: 0 });
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        fetchOrders();
        fetchStats();
    }, []);

    const fetchStats = async () => {
        try {
            const response = await orderAPI.getStats();
            setStats(response.data || { total_orders: 0, total_revenue: 0, pending_payments: 0 });
        } catch (err) {
            console.error("Failed to fetch stats:", err);
        }
    };

    const fetchOrders = async () => {
        try {
            const response = await orderAPI.getAll();
            setOrders(response.data || []);
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    const formatPrice = (price) => {
        return new Intl.NumberFormat('id-ID', {
            style: 'currency',
            currency: 'IDR',
            minimumFractionDigits: 0,
        }).format(price);
    };

    const formatDate = (dateString) => {
        return new Date(dateString).toLocaleDateString('id-ID', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    const getStatusBadge = (status) => {
        const statusClasses = {
            created: 'status-created',
            processing: 'status-processing',
            paid: 'status-paid', // Add this
            success: 'status-success', // Add this
            completed: 'status-completed',
            canceled: 'status-canceled',
        };
        return <span className={`status-badge ${statusClasses[status]}`}>{status}</span>;
    };

    if (loading) return <div className="loading">Loading...</div>;
    if (error) return <div className="error">{error}</div>;

    return (
        <div className="admin-page container">
            <div className="admin-header">
                <h1>Manage Orders</h1>
                <p className="subtitle">Overview of store performance</p>
            </div>

            {/* Stats Dashboard */}
            <div className="stats-grid" style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
                gap: '1.5rem',
                marginBottom: '2rem'
            }}>
                <div className="stat-card card" style={{ padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                    <h3 style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', textTransform: 'uppercase' }}>Total Revenue</h3>
                    <p style={{ fontSize: '1.5rem', fontWeight: '700', color: 'var(--success-color)' }}>
                        {formatPrice(stats.total_revenue)}
                    </p>
                </div>
                <div className="stat-card card" style={{ padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                    <h3 style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', textTransform: 'uppercase' }}>Total Orders</h3>
                    <p style={{ fontSize: '1.5rem', fontWeight: '700', color: 'var(--primary-color)' }}>
                        {stats.total_orders}
                    </p>
                </div>
                <div className="stat-card card" style={{ padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                    <h3 style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', textTransform: 'uppercase' }}>Pending Payments</h3>
                    <p style={{ fontSize: '1.5rem', fontWeight: '700', color: 'var(--warning-color)' }}>
                        {stats.pending_payments}
                    </p>
                </div>
            </div>

            <div className="admin-table">
                <table>
                    <thead>
                        <tr>
                            <th>Order Number</th>
                            <th>Customer</th>
                            <th>Date</th>
                            <th>Total</th>
                            <th>Status</th>
                            <th>Payment Status</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {orders.map((order) => (
                            <tr key={order.id}>
                                <td>
                                    <span style={{ fontFamily: 'monospace', fontWeight: 'bold' }}>
                                        {order.order_number || order.id.substring(0, 8).toUpperCase()}
                                    </span>
                                </td>
                                <td>{order.user?.name || order.user?.email || 'N/A'}</td>
                                <td>{formatDate(order.created_at)}</td>
                                <td>{formatPrice(order.total_amount)}</td>
                                <td>{getStatusBadge(order.status)}</td>
                                <td>
                                    {order.payment ? (
                                        <span className={`status-badge status-${order.payment.status}`}>
                                            {order.payment.status}
                                        </span>
                                    ) : (
                                        <span className="status-badge status-none">No Payment</span>
                                    )}
                                </td>
                                <td>
                                    <Link to={`/orders/${order.id}`} className="btn-view-small">View</Link>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default AdminOrdersPage;
