import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { productAPI } from '../../api';
import ProductCard from '../../components/products/ProductCard';
import '../../styles/products.css';

const ProductListPage = () => {
    const navigate = useNavigate();
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [cart, setCart] = useState([]);

    useEffect(() => {
        fetchProducts();
    }, []);

    const fetchProducts = async () => {
        try {
            const response = await productAPI.getAll();
            setProducts(response.data || []);
        } catch (err) {
            setError(err.message || 'Failed to fetch products');
        } finally {
            setLoading(false);
        }
    };

    const handleAddToCart = (product) => {
        const existingItem = cart.find(item => item.product_id === product.id);

        if (existingItem) {
            setCart(cart.map(item =>
                item.product_id === product.id
                    ? { ...item, quantity: item.quantity + 1 }
                    : item
            ));
        } else {
            setCart([...cart, { product_id: product.id, quantity: 1, product }]);
        }
    };

    const handleCheckout = () => {
        if (cart.length === 0) {
            alert('Your cart is empty');
            return;
        }
        // Navigate to create order page with cart data
        navigate('/orders/create', { state: { cart } });
    };

    if (loading) return <div className="loading">Loading products...</div>;
    if (error) return <div className="error">{error}</div>;

    return (
        <div className="products-page">
            <div className="page-header">
                <h1>Products</h1>
                {cart.length > 0 && (
                    <div className="cart-summary">
                        <span>{cart.length} items in cart</span>
                        <button onClick={handleCheckout} className="btn-checkout">
                            Checkout
                        </button>
                    </div>
                )}
            </div>

            {products.length === 0 ? (
                <div className="empty-state">No products available</div>
            ) : (
                <div className="products-grid">
                    {products.map((product) => (
                        <ProductCard
                            key={product.id}
                            product={product}
                            onAddToCart={handleAddToCart}
                        />
                    ))}
                </div>
            )}
        </div>
    );
};

export default ProductListPage;
