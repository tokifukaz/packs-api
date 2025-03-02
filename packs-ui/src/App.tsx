import { useState, useEffect } from 'react'
import './App.css'

interface Order {
  id: string;
  items: number;
  packSizes: number[];
  packQuantity: Record<string, number>;
  createdAt: string;
  updatedAt: string;
}

interface OrderRequest {
  items: number;
  packSizes: number[];
}

interface OrdersResponse {
  data: Order[];
}

function App() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string>('');
  
  // Form state
  const [itemCount, setItemCount] = useState<number>(0);
  const [packSizesInput, setPackSizesInput] = useState<string>('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string>('');


  // Define your API base URL
  const API_BASE_URL = import.meta.env.REACT_APP_API_URL || 'http://localhost:8001';  // Replace with your actual API URL

  useEffect(() => {
    fetchOrders();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setSubmitError('');

    try {
      // Convert comma-separated string to number array and sort
      const packSizes = packSizesInput
        .split(',')
        .map(size => parseInt(size.trim()))
        .filter(size => !isNaN(size))
        .sort((a, b) => a - b);

      const orderData: OrderRequest = {
        items: itemCount,
        packSizes: packSizes
      };

      const response = await fetch(`${API_BASE_URL}/api/orders`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(orderData),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      // Reset form
      setItemCount(0);
      setPackSizesInput('');
      
      // Refresh orders list
      await fetchOrders();
    } catch (error) {
      setSubmitError(error instanceof Error ? error.message : 'Failed to submit order');
    } finally {
      setIsSubmitting(false);
    }
  };

  const fetchOrders = async () => {
    try {
      setIsLoading(true);
      const response = await fetch(`${API_BASE_URL}/api/orders`);
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const result: OrdersResponse = await response.json();
      console.log('API Response:', result); // Add this log
      setOrders(result.data);
    } catch (error) {
      console.error('Fetch error:', error);
      setError(error instanceof Error ? error.message : 'Failed to fetch orders');
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }


  return (
    <div className="container">
      <h1>Orders List</h1>

      <div className="order-form">
        <h2>Add New Order</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="items">Items Count:</label>
            <input
              type="number"
              id="items"
              value={itemCount}
              onChange={(e) => setItemCount(parseInt(e.target.value))}
              min="1"
              required
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="packSizes">Pack Sizes (comma-separated):</label>
            <input
              type="text"
              id="packSizes"
              value={packSizesInput}
              onChange={(e) => setPackSizesInput(e.target.value)}
              placeholder="e.g., 250, 500, 1000, 2000"
              required
            />
          </div>

          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? 'Submitting...' : 'Add Order'}
          </button>

          {submitError && <div className="error">{submitError}</div>}
        </form>
      </div>

      {orders.length > 0 ? (
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Items</th>
              <th>Pack Sizes</th>
              <th>Pack Quantity</th>
              <th>Created At</th>
              <th>Updated At</th>
            </tr>
          </thead>
          <tbody>
            {orders.map((order) => (
              <tr key={order.id}>
                <td>{order.id}</td>
                <td>{order.items}</td>
                <td>{order.packSizes.join(', ')}</td>
                <td>
                  {Object.entries(order.packQuantity).map(([size, quantity]) => (
                    `${size}: ${quantity}`
                  )).join(', ')}
                </td>
                <td>{new Date(order.createdAt).toLocaleString()}</td>
                <td>{new Date(order.updatedAt).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      ) : (
        <div>No orders found</div>
      )}
    </div>
  )
}

export default App