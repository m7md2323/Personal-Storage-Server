document.addEventListener('DOMContentLoaded', loadUsers);

// Fetch and display users
async function loadUsers() {
    const grid = document.getElementById('user-grid');
    const res = await fetch('/api/users');
    const users = await res.json();

    grid.innerHTML = users.length ? '' : '<p>No profiles yet.</p>';
    
    users.forEach(user => {
        const div = document.createElement('div');
        div.className = 'user-card';
        div.innerHTML = `<strong>${user.username}</strong>`;
        div.onclick = () => window.location.href = `/dashboard?user=${user.id}`;
        grid.appendChild(div);
    });
}

// Add a new user
async function addUser() {
    const username = document.getElementById('username-input').value;
    if (!username) return alert("Enter a name");

    const res = await fetch('/api/users', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: username })
    });

    if (res.ok) {
        document.getElementById('username-input').value = '';
        loadUsers(); // Refresh the list
    }
}