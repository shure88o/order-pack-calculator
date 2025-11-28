// API base URL
const API_BASE = window.location.origin;

// State
let currentPackSizes = [];

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
    loadPackSizes();

    document.getElementById('addPackBtn').addEventListener('click', addPackSizeInput);
    document.getElementById('submitPacksBtn').addEventListener('click', updatePackSizes);
    document.getElementById('calculateBtn').addEventListener('click', calculatePacks);

    // Allow Enter key to trigger calculation
    document.getElementById('orderQty').addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            calculatePacks();
        }
    });
});

// Load current pack sizes from API
async function loadPackSizes() {
    try {
        const response = await fetch(`${API_BASE}/api/packs`);
        const data = await response.json();

        if (data.pack_sizes) {
            currentPackSizes = data.pack_sizes;
            renderPackSizes();
        }
    } catch (error) {
        console.error('Failed to load pack sizes:', error);
        showError('Failed to load pack sizes from server');
    }
}

// Render pack size inputs
function renderPackSizes() {
    const container = document.getElementById('packSizesContainer');
    container.innerHTML = '';

    currentPackSizes.forEach((size, index) => {
        const div = document.createElement('div');
        div.className = 'pack-size-item';

        const input = document.createElement('input');
        input.type = 'number';
        input.min = '1';
        input.value = size;
        input.dataset.index = index;

        const removeBtn = document.createElement('button');
        removeBtn.textContent = 'Remove';
        removeBtn.className = 'btn btn-danger';
        removeBtn.onclick = () => removePackSize(index);

        div.appendChild(input);
        div.appendChild(removeBtn);
        container.appendChild(div);
    });
}

// Add new pack size input
function addPackSizeInput() {
    currentPackSizes.push(250);
    renderPackSizes();
}

// Remove pack size
function removePackSize(index) {
    currentPackSizes.splice(index, 1);
    renderPackSizes();
}

// Update pack sizes via API
async function updatePackSizes() {
    const inputs = document.querySelectorAll('#packSizesContainer input');
    const newSizes = [];

    inputs.forEach(input => {
        const value = parseInt(input.value);
        if (value && value > 0) {
            newSizes.push(value);
        }
    });

    if (newSizes.length === 0) {
        showPackMessage('Please add at least one valid pack size', false);
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/api/packs`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ pack_sizes: newSizes }),
        });

        const data = await response.json();

        if (response.ok) {
            currentPackSizes = data.pack_sizes;
            renderPackSizes();
            showPackMessage(data.message || 'Pack sizes updated successfully', true);
        } else {
            showPackMessage(data.error || 'Failed to update pack sizes', false);
        }
    } catch (error) {
        console.error('Failed to update pack sizes:', error);
        showPackMessage('Failed to update pack sizes', false);
    }
}

// Calculate optimal packs
async function calculatePacks() {
    const orderQty = parseInt(document.getElementById('orderQty').value);

    if (isNaN(orderQty) || orderQty < 0) {
        showError('Please enter a valid order quantity');
        return;
    }

    hideError();

    try {
        const response = await fetch(`${API_BASE}/api/calculate`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ order_quantity: orderQty }),
        });

        const data = await response.json();

        if (response.ok) {
            displayResults(data);
        } else {
            showError(data.error || 'Failed to calculate packs');
        }
    } catch (error) {
        console.error('Failed to calculate packs:', error);
        showError('Failed to calculate packs');
    }
}

// Display calculation results
function displayResults(data) {
    const resultsContainer = document.getElementById('resultsContainer');
    const resultsBody = document.getElementById('resultsBody');
    const totalSummary = document.getElementById('totalSummary');

    resultsBody.innerHTML = '';

    if (data.packs && data.packs.length > 0) {
        data.packs.forEach(pack => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${pack.size}</td>
                <td>${pack.quantity}</td>
            `;
            resultsBody.appendChild(row);
        });

        totalSummary.textContent = `Total: ${data.total_items} items in ${data.total_packs} pack${data.total_packs !== 1 ? 's' : ''}`;
        resultsContainer.classList.remove('hidden');
    } else {
        totalSummary.textContent = 'No packs needed';
        resultsContainer.classList.remove('hidden');
    }
}

// Show pack configuration message
function showPackMessage(message, success) {
    const messageEl = document.getElementById('packMessage');
    messageEl.textContent = message;
    messageEl.className = success ? 'message success' : 'message error';
    messageEl.style.display = 'block';

    setTimeout(() => {
        messageEl.style.display = 'none';
    }, 3000);
}

// Show error message
function showError(message) {
    const errorEl = document.getElementById('errorMessage');
    errorEl.textContent = message;
    errorEl.classList.remove('hidden');
    document.getElementById('resultsContainer').classList.add('hidden');
}

// Hide error message
function hideError() {
    document.getElementById('errorMessage').classList.add('hidden');
}
