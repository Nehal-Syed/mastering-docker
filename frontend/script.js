// Use relative URL - will work with whatever port Nginx is on
const API_URL = '/api';

async function fetchProducts() {
    try {
        const response = await fetch(`${API_URL}/products`);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        
        const instanceId = response.headers.get('X-Instance-ID');
        if (instanceId) {
            document.getElementById('instanceId').textContent = instanceId;
        }
        const products = await response.json();
        displayProducts(products);
    } catch (error) {
        console.error('Error fetching products:', error);
        document.getElementById('productsList').innerHTML = '<div class="loading">⚠️ Cannot connect to backend. Make sure Docker is running.</div>';
    }
}

async function createProduct(product) {
    try {
        const response = await fetch(`${API_URL}/products`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(product)
        });
        
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return await response.json();
    } catch (error) {
        console.error('Error creating product:', error);
        throw error;
    }
}

// Also fix the price parsing - remove commas
document.getElementById('productForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const priceInput = document.getElementById('price').value;
    // Remove commas and any non-numeric characters except decimal point
    const cleanPrice = priceInput.replace(/,/g, '').replace(/[^0-9.]/g, '');
    
    const product = {
        name: document.getElementById('name').value,
        description: document.getElementById('description').value,
        price: parseFloat(cleanPrice),
        quantity: parseInt(document.getElementById('quantity').value) || 0
    };
    
    // Validate
    if (!product.name) {
        alert('Product name is required');
        return;
    }
    if (isNaN(product.price) || product.price <= 0) {
        alert('Valid price is required');
        return;
    }
    
    try {
        const response = await fetch(`${API_URL}/products`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(product)
        });
        
        if (response.ok) {
            document.getElementById('productForm').reset();
            fetchProducts();
        } else {
            const error = await response.text();
            alert(`Error creating product: ${error}`);
        }
    } catch (error) {
        console.error('Error creating product:', error);
        alert('Error creating product. Check if backend is running.');
    }
});

// Update the displayProducts function to handle numbers properly
function displayProducts(products) {
    const container = document.getElementById('productsList');
    if (!products || products.length === 0) {
        container.innerHTML = '<div class="loading">No products found. Add your first product!</div>';
        return;
    }

    container.innerHTML = products.map(product => `
        <div class="product-card">
            <h3>${escapeHtml(product.name)}</h3>
            <p>${escapeHtml(product.description) || 'No description'}</p>
            <div class="price">$${product.price.toFixed(2)}</div>
            <div class="quantity">📦 Stock: ${product.quantity}</div>
            <button class="edit" onclick="editProduct(${product.id})">✏️ Edit</button>
            <button onclick="deleteProduct(${product.id})">🗑️ Delete</button>
        </div>
    `).join('');
}