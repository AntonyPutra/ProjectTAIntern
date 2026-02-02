# Mini OMS Frontend

React + Vite frontend for the Mini Order Management System.

## Tech Stack

- **React 18.2** - UI library
- **Vite 5** - Build tool & dev server
- **React Router DOM 6** - Client-side routing
- **Axios** - HTTP client with interceptors
- **Context API** - State management for authentication

## Project Structure

```
src/
├── api/
│   ├── http.js              # Axios instance with JWT interceptor
│   └── index.js             # API service methods
├── components/
│   ├── common/
│   │   └── Navbar.jsx       # Navigation bar
│   ├── products/
│   │   └── ProductCard.jsx  # Product card component
│   ├── orders/              # (Empty, can add order components)
│   └── ...
├── context/
│   └── AuthContext.jsx      # Authentication context provider
├── pages/
│   ├── auth/
│   │   ├── LoginPage.jsx
│   │   └── RegisterPage.jsx
│   ├── products/
│   │   ├── ProductListPage.jsx
│   │   └── ProductDetailPage.jsx
│   ├── orders/
│   │   ├── OrderListPage.jsx
│   │   ├── OrderDetailPage.jsx
│   │   ├── CreateOrderPage.jsx
│   │   └── CreatePaymentPage.jsx
│   └── admin/
│       ├── AdminProductsPage.jsx
│       └── AdminOrdersPage.jsx
├── routes/
│   ├── AppRouter.jsx        # Main router configuration
│   ├── PrivateRoute.jsx     # Protected route guard
│   └── AdminRoute.jsx       # Admin-only route guard
├── styles/
│   ├── global.css           # Global styles & variables
│   ├── navbar.css
│   ├── auth.css
│   ├── products.css
│   ├── productcard.css
│   ├── orders.css
│   └── admin.css
├── utils/                   # (Empty, for utility functions)
├── App.jsx                  # Root app component
└── main.jsx                 # Entry point
```

## Setup & Installation

### 1. Install Dependencies

```bash
cd mini-oms-frontend
npm install
```

### 2. Configure Environment

Create `.env` file from template:
```bash
copy .env.example .env
```

Edit `.env`:
```
VITE_API_BASE_URL=http://localhost:8080
```

### 3. Run Development Server

```bash
npm run dev
```

The app will be available at `http://localhost:5173`

### 4. Build for Production

```bash
npm run build
```

## Features

### Authentication
- ✅ User Registration with validation
- ✅ User Login with JWT token storage
- ✅ Auto-redirect after authentication
- ✅ JWT token injection via axios interceptor
- ✅ Logout functionality

### Products
- ✅ Browse all products (public)
- ✅ View product details with stock info
- ✅ Add products to cart
- ✅ Checkout from cart

### Orders
- ✅ Create order with multiple items
- ✅ View order list (user sees own orders)
- ✅ View order details with items
- ✅ Manage cart (update quantities, remove items)

### Payments
- ✅ Create payment for order
- ✅ Select payment method
- ✅ Upload payment proof URL
- ✅ View payment status on order detail

### Admin Features (Admin-only)
- ✅ Full CRUD for products (Create, Read, Update, Delete)
- ✅ View all orders from all users
- ✅ See order details with customer info
- ✅ See payment status for all orders

## State Management Decision

**Chosen: Context API** instead of Zustand

**Rationale:**
- ✅ Built-in to React (no extra dependencies)
- ✅ Perfect for simple authentication state
- ✅ Easy for interns to understand and maintain
- ✅ Sufficient for this project's scope
- ⚠️ For complex state with multiple slices, Zustand would be better

## Route Protection

### Public Routes
- `/login` - Login page
- `/register` - Registration page
- `/products` - Product list (anyone can view)
- `/products/:id` - Product detail

### Protected Routes (JWT required)
- `/orders` - User's order list
- `/orders/create` - Create new order
- `/orders/:id` - Order detail
- `/orders/:orderId/payment` - Create payment

### Admin Routes (JWT + admin role)
- `/admin/products` - Manage products (CRUD)
- `/admin/orders` - View all orders

## API Integration

All API calls use the axios instance (`src/api/http.js`) which:
1. Automatically adds JWT token to Authorization header
2. Handles 401 responses by redirecting to login
3. Extracts response data automatically
4. Provides consistent error messages

## User Flows

### 1. New User Registration
```
Register → Auto-login → Redirect to Products
```

### 2. Existing User Login
```
Login → Redirect to Products
```

### 3. Shopping & Ordering
```
Products → Add to Cart → Checkout → Create Order → View Order → Create Payment
```

### 4. Admin Product Management
```
Admin Products → Create/Edit/Delete Products
```

### 5. Admin Order Management
```
Admin Orders → View All Orders → See Customer & Payment Info
```

## Development Notes

### Adding New Pages
1. Create page component in `src/pages/`
2. Add route in `src/routes/AppRouter.jsx`
3. Use `PrivateRoute` or `AdminRoute` if needed
4. Create CSS file in `src/styles/` if needed

### Adding New API Endpoints
1. Add method to `src/api/index.js`
2. Use the `apiClient` instance for automatic token injection

### Styling
- Uses CSS variables for consistent theming
- Responsive design (mobile-friendly)
- Professional but minimal UI focus

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_BASE_URL` | Backend API URL | `http://localhost:8080` |

## Browser Support

- Modern browsers (Chrome, Firefox, Safari, Edge)
- ES6+ features used

## Next Steps for Students

1. **Test all user flows** - Register, login, browse, order, payment
2. **Test admin flows** - Create/edit/delete products, view all orders
3. **Add features**:
   - Search & filter products
   - Pagination for products/orders
   - User profile page
   - Order status update (admin)
   - Payment verification (admin)
4. **Improve UX**:
   - Add loading skeletons
   - Toast notifications
   - Form validation messages
   - Image upload for products

## Troubleshooting

### API Connection Issues
- Ensure backend is running on `http://localhost:8080`
- Check `.env` file has correct `VITE_API_BASE_URL`
- Check browser console for CORS errors

### Authentication Issues
- Clear localStorage and try logging in again
- Check if JWT token is in localStorage: `localStorage.getItem('access_token')`
- Verify token is being sent in requests (Network tab)

## License

MIT
