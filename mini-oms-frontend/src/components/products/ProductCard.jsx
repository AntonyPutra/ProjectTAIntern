import { Link } from 'react-router-dom';
import '../../styles/productcard.css';

const ProductCard = ({ product, onAddToCart }) => {
    const formatPrice = (price) => {
        return new Intl.NumberFormat('id-ID', {
            style: 'currency',
            currency: 'IDR',
            minimumFractionDigits: 0,
        }).format(price);
    };

    return (
        <div className="product-card">
            {product.image_url && (
                <img src={product.image_url} alt={product.name} className="product-image" />
            )}
            <div className="product-info">
                <h3 className="product-name">{product.name}</h3>
                <p className="product-description">{product.description}</p>
                <div className="product-footer">
                    <span className="product-price">{formatPrice(product.price)}</span>
                    <span className={`product-stock ${product.stock === 0 ? 'out-of-stock' : ''}`}>
                        Stock: {product.stock}
                    </span>
                </div>
                <div className="product-actions">
                    <Link to={`/products/${product.id}`} className="btn-view">
                        View Details
                    </Link>
                    {product.stock > 0 && onAddToCart && (
                        <button onClick={() => onAddToCart(product)} className="btn-add-cart">
                            Add to Cart
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
};

export default ProductCard;
