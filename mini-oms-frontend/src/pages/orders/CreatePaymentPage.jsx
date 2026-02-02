import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { paymentAPI, orderAPI } from '../../api';
import '../../styles/orders.css';

const CreatePaymentPage = () => {
    const { orderId } = useParams();
    const navigate = useNavigate();
    const [order, setOrder] = useState(null);
    const [formData, setFormData] = useState({
        payment_method: 'bank_transfer',
        payment_proof_url: '',
        notes: '',
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    useEffect(() => {
        fetchOrder();
    }, [orderId]);

    const fetchOrder = async () => {
        try {
            const response = await orderAPI.getById(orderId);
            setOrder(response.data);
        } catch (err) {
            setError('Failed to fetch order details');
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        const paymentData = {
            order_id: orderId,
            ...formData,
        };

        try {
            await paymentAPI.create(paymentData);
            navigate(`/orders/${orderId}`);
        } catch (err) {
            setError(err.message || 'Failed to create payment');
        } finally {
            setLoading(false);
        }
    };

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value,
        });
    };

    const formatPrice = (price) => {
        return new Intl.NumberFormat('id-ID', {
            style: 'currency',
            currency: 'IDR',
            minimumFractionDigits: 0,
        }).format(price);
    };

    if (!order) return <div className="loading">Loading...</div>;

    return (
        <div className="create-payment-page">
            <h1>Create Payment</h1>

            <div className="payment-summary">
                <h3>Order Summary</h3>
                <p><strong>Order Number:</strong> {order.order_number}</p>
                <p><strong>Total Amount:</strong> {formatPrice(order.total_amount)}</p>
            </div>

            {error && <div className="error-message">{error}</div>}

            <form onSubmit={handleSubmit} className="payment-form">
                <div className="form-group">
                    <label htmlFor="payment_method">Payment Method</label>
                    <select
                        id="payment_method"
                        name="payment_method"
                        value={formData.payment_method}
                        onChange={handleChange}
                        required
                    >
                        <option value="bank_transfer">Bank Transfer</option>
                        <option value="e-wallet">E-Wallet</option>
                        <option value="credit_card">Credit Card</option>
                    </select>
                </div>

                <div className="form-group">
                    <label htmlFor="payment_proof_url">Payment Proof URL</label>
                    <input
                        type="url"
                        id="payment_proof_url"
                        name="payment_proof_url"
                        value={formData.payment_proof_url}
                        onChange={handleChange}
                        placeholder="https://example.com/proof.jpg"
                        required
                    />
                    <small>Upload your payment proof to an image hosting service and paste the URL here</small>
                </div>

                <div className="form-group">
                    <label htmlFor="notes">Notes (Optional)</label>
                    <textarea
                        id="notes"
                        name="notes"
                        value={formData.notes}
                        onChange={handleChange}
                        placeholder="Transfer from BCA a/n John Doe"
                        rows="3"
                    />
                </div>

                <div className="form-actions">
                    <button type="button" onClick={() => navigate(`/orders/${orderId}`)} className="btn-cancel">
                        Cancel
                    </button>
                    <button type="submit" className="btn-submit" disabled={loading}>
                        {loading ? 'Submitting...' : 'Submit Payment'}
                    </button>
                </div>
            </form>
        </div>
    );
};

export default CreatePaymentPage;
