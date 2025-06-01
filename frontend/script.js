// Configuration
const API_BASE_URL = 'https://linkreducer-production.up.railway.app'; // Update this to match your backend URL

// DOM Elements
const shortenForm = document.getElementById('shortenForm');
const resultSection = document.getElementById('result');
const toast = document.getElementById('toast');

// Initialize the app
document.addEventListener('DOMContentLoaded', function() {
    setupEventListeners();
    setMinDate();
    initializeTheme();
});

// Set minimum date to current date
function setMinDate() {
    const now = new Date();
    const localDate = new Date(now.getTime() - now.getTimezoneOffset() * 60000);
    document.getElementById('expirationDate').min = localDate.toISOString().slice(0, 16);
}

// Setup event listeners
function setupEventListeners() {
    shortenForm.addEventListener('submit', handleShortenURL);
    document.getElementById('copyBtn').addEventListener('click', copyToClipboard);
    document.getElementById('themeToggle').addEventListener('click', toggleTheme);
}

// Handle URL shortening form submission
async function handleShortenURL(e) {
    e.preventDefault();
    
    const originalUrl = document.getElementById('originalUrl').value;
    const customCode = document.getElementById('customCode').value;
    const expirationDate = document.getElementById('expirationDate').value;
    
    const payload = {
        original_url: originalUrl
    };
    
    if (customCode) {
        payload.short_code = customCode;
    }
    
    if (expirationDate) {
        payload.expiration_date = new Date(expirationDate).toISOString();
    }
    
    try {
        showLoading(true);
        const response = await fetch(`${API_BASE_URL}/urls`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(payload)
        });
        
        const data = await response.json();
          if (response.ok) {
            displayResult(data);
            shortenForm.reset();
            showToast('URL shortened successfully!', 'success');
        } else {
            throw new Error(data.error || 'Failed to shorten URL');
        }
    } catch (error) {
        console.error('Error:', error);
        showToast(error.message, 'error');
    } finally {
        showLoading(false);
    }
}

// Display the shortened URL result
function displayResult(data) {
    document.getElementById('shortUrlDisplay').value = data.short_url;
    document.getElementById('originalUrlDisplay').textContent = data.original_url;
    document.getElementById('createdAtDisplay').textContent = formatDate(data.created_at);
    document.getElementById('expirationDisplay').textContent = formatDate(data.expiration_date);
    document.getElementById('hitCountDisplay').textContent = data.hit_count || 0;
    
    resultSection.style.display = 'block';
    resultSection.scrollIntoView({ behavior: 'smooth' });
}

// Copy shortened URL to clipboard
async function copyToClipboard() {
    const shortUrl = document.getElementById('shortUrlDisplay').value;
    
    try {
        await navigator.clipboard.writeText(shortUrl);
        showToast('URL copied to clipboard!', 'success');
    } catch (error) {
        // Fallback for older browsers
        const textArea = document.createElement('textarea');
        textArea.value = shortUrl;
        document.body.appendChild(textArea);
        textArea.select();
        document.execCommand('copy');
        document.body.removeChild(textArea);
        showToast('URL copied to clipboard!', 'success');
    }
}

// Get full short URL
function getShortUrl(shortCode) {
    // Extract domain from API_BASE_URL or use localhost
    const domain = API_BASE_URL.replace('http://', '').replace('https://', '').replace(':8080', '');    return `${API_BASE_URL.includes('https') ? 'https' : 'http'}://${domain}/${shortCode}`;
}

// Show toast notification
function showToast(message, type = 'success') {
    toast.textContent = message;
    toast.className = `toast ${type}`;
    toast.classList.add('show');
    
    setTimeout(() => {
        toast.classList.remove('show');
    }, 3000);
}

// Utility functions
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
}

function truncateText(text, maxLength) {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
}

// Theme Management
function initializeTheme() {
    const savedTheme = localStorage.getItem('theme');
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    
    if (savedTheme === 'dark' || (!savedTheme && prefersDark)) {
        document.body.classList.add('dark-mode');
        updateThemeIcon(true);
    }
}

function toggleTheme() {
    const isDark = document.body.classList.toggle('dark-mode');
    localStorage.setItem('theme', isDark ? 'dark' : 'light');
    updateThemeIcon(isDark);
}

function updateThemeIcon(isDark) {
    const themeToggle = document.getElementById('themeToggle');
    const icon = themeToggle.querySelector('i');
    
    if (isDark) {
        icon.className = 'fas fa-sun';
        themeToggle.title = 'Switch to Light Mode';
    } else {
        icon.className = 'fas fa-moon';
        themeToggle.title = 'Switch to Dark Mode';
    }
}
