import { Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import '../../styles/navbar.css';

const Navbar = () => {
    const { user, logout, isAdmin } = useAuth();

    return (
        <nav className="navbar">
            <div className="nav-container">
                <Link to="/products" className="nav-brand">
                    Mini OMS
                </Link>

                <div className="nav-menu">
                    <Link to="/products" className="nav-link">
                        Products
                    </Link>

                    {user && (
                        <>
                            <Link to="/orders" className="nav-link">
                                My Orders
                            </Link>

                            {isAdmin() && (
                                <>
                                    <Link to="/admin/products" className="nav-link">
                                        Manage Products
                                    </Link>
                                    <Link to="/admin/orders" className="nav-link">
                                        Manage Orders
                                    </Link>
                                </>
                            )}
                        </>
                    )}
                </div>

                <div className="nav-auth">
                    {user ? (
                        <>
                            <span className="user-name">
                                {user.name} {isAdmin() && <span className="badge-admin">Admin</span>}
                            </span>
                            <button onClick={logout} className="btn-logout">
                                Logout
                            </button>
                        </>
                    ) : (
                        <>
                            <Link to="/login" className="btn-link">
                                Login
                            </Link>
                            <Link to="/register" className="btn-primary">
                                Register
                            </Link>
                        </>
                    )}
                </div>
            </div>
        </nav>
    );
};

export default Navbar;
