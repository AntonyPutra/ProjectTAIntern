import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { orderAPI, paymentAPI } from '../../api';
import '../../styles/orders.css';

const OrderDetailPage = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [order, setOrder] = useState(null);
    const [payment, setPayment] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        fetchOrderDetails();
    }, [id]);

    const fetchOrderDetails = async () => {
        try {
            const orderResponse = await orderAPI.getById(id);
            setOrder(orderResponse.data);

            // Try to fetch payment if exists
            try {
                const paymentResponse = await paymentAPI.getByOrderId(id);
                setPayment(paymentResponse.data);
            } catch (err) {
                // Payment doesn't exist yet, that's ok
            }
        } catch (err) {
            setError(err.message || 'Failed to fetch order details');
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

    if (loading) return <div className="loading">Loading...</div>;
    if (error) return <div className="error">{error}</div>;
    if (!order) return <div className="error">Order not found</div>;

    return (
        <div className="order-detail-page">
            <button onClick={() => navigate('/orders')} className="btn-back">
                ← Back to Orders
            </button>

            <div className="order-header">
                <h1>Order Details</h1>
                <div className="order-number">{order.order_number}</div>
            </div>

            <div className="order-info-grid">
                <div className="info-card">
                    <h3>Order Information</h3>
                    <p><strong>Date:</strong> {formatDate(order.created_at)}</p>
                    <p><strong>Status:</strong> <span className={`status-${order.status}`}>{order.status}</span></p>
                    <p><strong>Total Amount:</strong> {formatPrice(order.total_amount)}</p>
                    {order.notes && <p><strong>Notes:</strong> {order.notes}</p>}
                </div>

                {payment && (
                    <div className="info-card">
                        <h3>Payment Information</h3>
                        <p><strong>Payment Number:</strong> {payment.payment_number}</p>
                        <p><strong>Method:</strong> {payment.payment_method}</p>
                        <p><strong>Status:</strong> <span className={`status-${payment.status}`}>{payment.status}</span></p>
                        <p><strong>Amount:</strong> {formatPrice(payment.amount)}</p>
                    </div>
                )}
            </div>

            <div className="order-items-section">
                <h3>Order Items</h3>
                <table className="items-table">
                    <thead>
                        <tr>
                            <th>Product</th>
                            <th>Price</th>
                            <th>Quantity</th>
                            <th>Subtotal</th>
                        </tr>
                    </thead>
                    <tbody>
                        {order.items?.map((item) => (
                            <tr key={item.id}>
                                <td>{item.product_name}</td>
                                <td>{formatPrice(item.product_price)}</td>
                                <td>{item.quantity}</td>
                                <td>{formatPrice(item.subtotal)}</td>
                            </tr>
                        ))}
                    </tbody>
                    <tfoot>
                        <tr>
                            <td colSpan="3"><strong>Total</strong></td>
                            <td><strong>{formatPrice(order.total_amount)}</strong></td>
                        </tr>
                    </tfoot>
                </table>
            </div>

            {!payment && order.status !== 'canceled' && (
                <div className="payment-action">
                    <Link to={`/orders/${order.id}/payment`} className="btn-payment">
                        Create Payment
                    </Link>
                </div>
            )}
        </div>
    );
};

export default OrderDetailPage;
