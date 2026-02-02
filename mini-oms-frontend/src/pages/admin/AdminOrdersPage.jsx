import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { orderAPI } from '../../api';
import '../../styles/admin.css';

const AdminOrdersPage = () => {
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        fetchOrders();
    }, []);

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
            completed: 'status-completed',
            canceled: 'status-canceled',
        };
        return <span className={`status-badge ${statusClasses[status]}`}>{status}</span>;
    };

    if (loading) return <div className="loading">Loading...</div>;
    if (error) return <div className="error">{error}</div>;

    return (
        <div className="admin-page">
            <div className="admin-header">
                <h1>Manage Orders</h1>
                <p className="subtitle">View and manage all orders</p>
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
                                <td>{order.order_number}</td>
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
