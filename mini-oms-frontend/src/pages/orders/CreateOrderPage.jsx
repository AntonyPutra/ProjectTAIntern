import { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { orderAPI } from '../../api';
import '../../styles/orders.css';

const CreateOrderPage = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const [cart, setCart] = useState(location.state?.cart || []);
    const [notes, setNotes] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    useEffect(() => {
        if (!cart || cart.length === 0) {
            navigate('/products');
        }
    }, [cart, navigate]);

    const calculateTotal = () => {
        return cart.reduce((total, item) => {
            return total + (item.product.price * item.quantity);
        }, 0);
    };

    const formatPrice = (price) => {
        return new Intl.NumberFormat('id-ID', {
            style: 'currency',
            currency: 'IDR',
            minimumFractionDigits: 0,
        }).format(price);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        const orderData = {
            items: cart.map(item => ({
                product_id: item.product_id,
                quantity: item.quantity,
            })),
            notes,
        };

        try {
            const response = await orderAPI.create(orderData);
            navigate(`/orders/${response.data.id}`);
        } catch (err) {
            setError(err.message || 'Failed to create order');
        } finally {
            setLoading(false);
        }
    };

    const updateQuantity = (productId, newQuantity) => {
        if (newQuantity < 1) return;
        setCart(cart.map(item =>
            item.product_id === productId
                ? { ...item, quantity: newQuantity }
                : item
        ));
    };

    const removeItem = (productId) => {
        setCart(cart.filter(item => item.product_id !== productId));
    };

    return (
        <div className="create-order-page">
            <h1>Create Order</h1>

            {error && <div className="error-message">{error}</div>}

            <div className="order-summary">
                <h2>Order Items</h2>
                <div className="cart-items">
                    {cart.map(item => (
                        <div key={item.product_id} className="cart-item">
                            <div className="item-info">
                                <h3>{item.product.name}</h3>
                                <p className="item-price">{formatPrice(item.product.price)}</p>
                            </div>
                            <div className="item-controls">
                                <button onClick={() => updateQuantity(item.product_id, item.quantity - 1)}>-</button>
                                <span className="item-quantity">{item.quantity}</span>
                                <button onClick={() => updateQuantity(item.product_id, item.quantity + 1)}>+</button>
                                <button onClick={() => removeItem(item.product_id)} className="btn-remove">Remove</button>
                            </div>
                            <div className="item-subtotal">
                                {formatPrice(item.product.price * item.quantity)}
                            </div>
                        </div>
                    ))}
                </div>

                <div className="order-total">
                    <strong>Total:</strong>
                    <span>{formatPrice(calculateTotal())}</span>
                </div>
            </div>

            <form onSubmit={handleSubmit} className="order-form">
                <div className="form-group">
                    <label htmlFor="notes">Order Notes (Optional)</label>
                    <textarea
                        id="notes"
                        value={notes}
                        onChange={(e) => setNotes(e.target.value)}
                        placeholder="Special instructions or delivery notes"
                        rows="4"
                    />
                </div>

                <div className="form-actions">
                    <button type="button" onClick={() => navigate('/products')} className="btn-cancel">
                        Cancel
                    </button>
                    <button type="submit" className="btn-submit" disabled={loading}>
                        {loading ? 'Creating Order...' : 'Place Order'}
                    </button>
                </div>
            </form>
        </div>
    );
};

export default CreateOrderPage;
