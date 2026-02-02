import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { productAPI } from '../../api';
import '../../styles/products.css';

const ProductDetailPage = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [product, setProduct] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [quantity, setQuantity] = useState(1);

    useEffect(() => {
        fetchProduct();
    }, [id]);

    const fetchProduct = async () => {
        try {
            const response = await productAPI.getById(id);
            setProduct(response.data);
        } catch (err) {
            setError(err.message || 'Failed to fetch product');
        } finally {
            setLoading(false);
        }
    };

    const handleAddToCart = () => {
        const cart = [{ product_id: product.id, quantity, product }];
        navigate('/orders/create', { state: { cart } });
    };

    const formatPrice = (price) => {
        return new Intl.NumberFormat('id-ID', {
            style: 'currency',
            currency: 'IDR',
            minimumFractionDigits: 0,
        }).format(price);
    };

    if (loading) return <div className="loading">Loading...</div>;
    if (error) return <div className="error">{error}</div>;
    if (!product) return <div className="error">Product not found</div>;

    return (
        <div className="product-detail-page">
            <button onClick={() => navigate(-1)} className="btn-back">
                ← Back
            </button>

            <div className="product-detail">
                {product.image_url && (
                    <div className="product-image-large">
                        <img src={product.image_url} alt={product.name} />
                    </div>
                )}

                <div className="product-info-large">
                    <h1>{product.name}</h1>
                    <p className="product-price-large">{formatPrice(product.price)}</p>
                    <p className="product-stock-info">
                        Stock: <strong>{product.stock}</strong>
                    </p>
                    <p className="product-description-large">{product.description}</p>

                    {product.stock > 0 && (
                        <div className="product-actions-large">
                            <div className="quantity-selector">
                                <label>Quantity:</label>
                                <input
                                    type="number"
                                    min="1"
                                    max={product.stock}
                                    value={quantity}
                                    onChange={(e) => setQuantity(parseInt(e.target.value))}
                                />
                            </div>
                            <button onClick={handleAddToCart} className="btn-order">
                                Order Now
                            </button>
                        </div>
                    )}

                    {product.stock === 0 && (
                        <div className="out-of-stock-notice">Out of Stock</div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default ProductDetailPage;
