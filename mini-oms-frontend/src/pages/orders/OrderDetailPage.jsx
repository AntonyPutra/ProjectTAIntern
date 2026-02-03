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

    const [verifying, setVerifying] = useState(false);

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

    const isAdmin = () => {
        try {
            const userStr = localStorage.getItem('user');
            if (userStr) {
                const user = JSON.parse(userStr);
                return user?.role === 'admin';
            }
        } catch (e) {
            console.error("Error parsing user from local storage", e);
            return false;
        }
        return false;
    };

    const handleVerifyPayment = async (paymentId) => {
        if (!window.confirm('Are you sure you want to verify this payment?')) return;

        setVerifying(true);
        try {
            await paymentAPI.verify(paymentId);
            alert('Payment verified successfully!');
            fetchOrderDetails(); // Refresh data
        } catch (err) {
            alert('Failed to verify payment: ' + (err.response?.data?.message || err.message));
        } finally {
            setVerifying(false);
        }
    };

    const handleCancelOrder = async () => {
        if (!window.confirm('Are you sure you want to cancel this order? Stock will be restored.')) return;

        try {
            await orderAPI.cancel(id);
            alert('Order canceled successfully!');
            fetchOrderDetails(); // Refresh data
        } catch (err) {
            alert('Failed to cancel order: ' + (err.response?.data?.message || err.message));
        }
    };

    const formatPrice = (price) => {
        if (price === undefined || price === null) return 'Rp 0';
        try {
            return new Intl.NumberFormat('id-ID', {
                style: 'currency',
                currency: 'IDR',
                minimumFractionDigits: 0,
            }).format(price);
        } catch (e) {
            console.error("Format price error", e);
            return 'Rp ' + price;
        }
    };

    const formatDate = (dateString) => {
        try {
            if (!dateString) return '-';
            return new Date(dateString).toLocaleDateString('id-ID', {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
            });
        } catch (e) {
            console.error("Format date error", e);
            return dateString || '-';
        }
    };

    if (loading) return <div className="loading">Loading...</div>;
    if (error) return <div className="error">{error}</div>;
    if (!order) return <div className="error">Order not found</div>;

    // Debug logging
    console.log("Rendering Order:", order);
    console.log("Rendering Payment:", payment);

    return (
        <div className="order-detail-page">
            <button onClick={() => navigate('/orders')} className="btn-back">
                ← Back to Orders
            </button>

            <div className="order-header">
                <h1>Order Details</h1>
                <div className="order-number">{order?.order_number || order?.id}</div>
            </div>

            <div className="order-info-grid">
                <div className="info-card">
                    <h3>Order Information</h3>
                    <p><strong>Date:</strong> {formatDate(order?.created_at)}</p>
                    <p><strong>Status:</strong> <span className={`status-${order?.status}`}>{order?.status || 'unknown'}</span></p>
                    <p><strong>Total Amount:</strong> {formatPrice(order?.total_amount)}</p>
                    {order?.notes && <p><strong>Notes:</strong> {order.notes}</p>}
                </div>

                {payment && (
                    <div className="info-card">
                        <h3>Payment Information</h3>
                        <p><strong>Payment Number:</strong> {payment?.payment_number || '-'}</p>
                        <p><strong>Method:</strong> {payment?.payment_method || '-'}</p>
                        <p><strong>Status:</strong> <span className={`status-${payment?.status}`}>{payment?.status || 'unknown'}</span></p>
                        <p><strong>Amount:</strong> {formatPrice(payment?.amount)}</p>
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
                        {(order?.items || []).map((item, index) => (
                            <tr key={item?.id || index}>
                                <td>{item?.product_name || 'Product'}</td>
                                <td>{formatPrice(item?.product_price)}</td>
                                <td>{item?.quantity || 0}</td>
                                <td>{formatPrice(item?.subtotal)}</td>
                            </tr>
                        ))}
                    </tbody>
                    <tfoot>
                        <tr>
                            <td colSpan="3"><strong>Total</strong></td>
                            <td><strong>{formatPrice(order?.total_amount)}</strong></td>
                        </tr>
                    </tfoot>
                </table>
            </div>

            {!payment && order?.status !== 'canceled' && (
                <div className="payment-action" style={{ display: 'flex', gap: '1rem', justifyContent: 'center' }}>
                    <Link to={`/orders/${order.id}/payment`} className="btn-payment">
                        Create Payment
                    </Link>

                    {order?.status === 'created' && (
                        <button
                            onClick={handleCancelOrder}
                            className="btn-payment"
                            style={{ backgroundColor: '#ef4444' }} // Red color
                        >
                            Cancel Order
                        </button>
                    )}
                </div>
            )}

            {payment && payment.status === 'pending' && isAdmin() && (
                <div className="payment-action">
                    <button
                        onClick={() => handleVerifyPayment(payment.id)}
                        className="btn-payment btn-verify"
                        style={{ backgroundColor: '#10b981', marginTop: '10px' }}
                        disabled={verifying}
                    >
                        {verifying ? 'Verifying...' : 'Verify Payment (Admin)'}
                    </button>
                </div>
            )}
        </div>
    );
};

export default OrderDetailPage;
