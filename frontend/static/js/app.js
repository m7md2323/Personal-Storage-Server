document.addEventListener('DOMContentLoaded', () => {
    fetchUsers();
});

// 1. Get all users from the Go Backend
async function fetchUsers() {
    try {
        const response = await fetch('/api/users');
        const users = await response.json();
        
        const listContainer = document.getElementById('user-list');
        listContainer.innerHTML = ''; // Clear loading text

        if (users.length === 0) {
            listContainer.innerHTML = '<p>No users found. Create one below!</p>';
            return;
        }

        users.forEach(user => {
            const div = document.createElement('div');
            div.className = 'user-item';
            div.innerText = user.username;
            div.onclick = () => loginAs(user.id);
            listContainer.appendChild(div);
        });
    } catch (err) {
        console.error("Failed to load users", err);
    }
}

// 2. Request a 10-digit code for a new user
async function requestAccess() {
    const username = document.getElementById('new-username').value;
    if (!username) return alert("Enter a username first!");

    const response = await fetch('/api/request-access', {
        method: 'POST',
        body: JSON.stringify({ username: username })
    });

    if (response.ok) {
        // Redirect to the verify screen we discussed
        window.location.href = "/verify.html";
    }
}

function loginAs(userId) {
    // This will lead to the authentication step
    console.log("Logging in as user ID:", userId);
}