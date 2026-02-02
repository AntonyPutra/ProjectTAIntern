import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { orderAPI } from '../../api';
import '../../styles/orders.css';

const OrderListPage = () => {
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
            setError(err.message || 'Failed to fetch orders');
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
            month: 'long',
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

    if (loading) return <div className="loading">Loading orders...</div>;
    if (error) return <div className="error">{error}</div>;

    return (
        <div className="orders-page">
            <h1>My Orders</h1>

            {orders.length === 0 ? (
                <div className="empty-state">
                    <p>You haven't placed any orders yet</p>
                    <Link to="/products" className="btn-primary">Browse Products</Link>
                </div>
            ) : (
                <div className="orders-table">
                    <table>
                        <thead>
                            <tr>
                                <th>Order Number</th>
                                <th>Date</th>
                                <th>Total</th>
                                <th>Status</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {orders.map((order) => (
                                <tr key={order.id}>
                                    <td>{order.order_number}</td>
                                    <td>{formatDate(order.created_at)}</td>
                                    <td>{formatPrice(order.total_amount)}</td>
                                    <td>{getStatusBadge(order.status)}</td>
                                    <td>
                                        <Link to={`/orders/${order.id}`} className="btn-view-small">
                                            View Details
                                        </Link>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
};

export default OrderListPage;
