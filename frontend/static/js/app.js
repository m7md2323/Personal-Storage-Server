/**
 * MyCloud Pro - Unified Frontend Logic
 * Version 3.1: Final Stabilized Build
 */

document.addEventListener('DOMContentLoaded', () => {
    // 1. Determine which page we are on and load specific data
    if (document.getElementById('user-list')) fetchUsers();
    if (document.getElementById('file-list')) loadFiles();
    if (document.getElementById('photo-grid')) loadPhotos();

    // 2. Initialize dynamic event listeners
    setupEventListeners();
    setupUploadZone();
});

// --- AUTHENTICATION & USER SELECTION ---

let selectedUser = "";

function setupEventListeners() {
    // Handle Login Form Submission
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', (e) => {
            e.preventDefault();
            handleLogin();
        });
    }

    // Handle Create User Form Submission
    const createForm = document.getElementById('create-user-form');
    if (createForm) {
        createForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const username = document.getElementById('new-username').value;
            const password = document.getElementById('new-password').value;
            const btn = createForm.querySelector('button');

            btn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Creating...';
            btn.disabled = true;

            try {
                const res = await fetch('/api/create_user', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username, password })
                });

                if (res.ok) {
                    alert(`Welcome ${username}! Your private storage is ready.`);
                    window.location.href = '/'; 
                } else {
                    const data = await res.json();
                    alert("Error: " + (data.error || "Could not create user"));
                }
            } catch (err) {
                alert("Connection to USB server lost.");
            } finally {
                btn.innerHTML = 'Create User';
                btn.disabled = false;
            }
        });
    }
}

async function fetchUsers() {
    const list = document.getElementById('user-list');
    if (!list) return;

    try {
        const response = await fetch('/api/users');
        const users = await response.json(); // If SQLite returns objects, this is [{username: '...'}, ...]

        list.innerHTML = users.map(user => {
            // Check if user is an object or a string to be safe
            const name = typeof user === 'object' ? user.username : user;
            
            return `
                <div class="user-avatar" onclick="showPasswordInput('${name}')">
                    <i class="fas fa-user-circle avatar-icon"></i>
                    <div class="user-name">${name}</div>
                </div>`;
        }).join('');
    } catch (err) {
        list.innerHTML = '<p>Error loading users.</p>';
    }
}

function showPasswordInput(username) {
    selectedUser = username;
    const userGrid = document.getElementById('user-list');
    const loginForm = document.getElementById('login-form');
    
    userGrid.style.opacity = '0';
    setTimeout(() => {
        userGrid.style.display = 'none';
        loginForm.style.display = 'block';
        document.getElementById('selected-user-text').innerText = "Logging in as " + username;
        document.getElementById('password').focus();
    }, 300);
}

async function handleLogin(e) {
    if (e) e.preventDefault(); // Stop the form from reloading the page
    
    const password = document.getElementById('password').value;

    const res = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
            username: selectedUser, // This was set by showPasswordInput
            password: password 
        })
    });

    if (res.status === 200) {
        window.location.href = '/my_files';
    } else {
        alert("Invalid password");
    }
}

// --- FILE EXPLORER (USER ISOLATED) ---

async function loadFiles() {
    const list = document.getElementById('file-list');
    try {
        const response = await fetch('/api/explorer'); 
        const items = await response.json();
        list.innerHTML = '';

        if (!items || items.length === 0) {
            list.innerHTML = '<div class="loading-state"><p>No files in your storage yet.</p></div>';
            return;
        }

        items.forEach(item => {
            const rawUrl = `/api/raw?name=${encodeURIComponent(item.name)}`;
            const isPhoto = /\.(jpg|jpeg|png|gif|webp)$/i.test(item.name);
            const isDoc = /\.(pdf|txt|log)$/i.test(item.name);

            let actions = "";
            if (!item.isDir) {
                if (isPhoto) actions += `<button class="action-btn" onclick="openViewer('${rawUrl}')"><i class="fas fa-eye"></i></button>`;
                if (isDoc) actions += `<a href="${rawUrl}" target="_blank" class="action-btn"><i class="fas fa-external-link-alt"></i></a>`;
                actions += `<a href="${rawUrl}" download class="action-btn"><i class="fas fa-download"></i></a>`;
            }

            const icon = item.isDir ? '<i class="fas fa-folder"></i>' : (isPhoto ? '<i class="fas fa-image"></i>' : '<i class="fas fa-file-alt"></i>');
            list.innerHTML += `
                <div class="file-row">
                    <div class="file-info"><span class="file-icon">${icon}</span> ${item.name}</div>
                    <div class="file-actions">${actions}</div>
                </div>`;
        });
    } catch (e) {
        list.innerHTML = '<div class="loading-state"><p>Error accessing your USB folder.</p></div>';
    }
}

// --- PHOTO GALLERY ---

let galleryImages = [];
let currentPhotoIndex = 0;

async function loadPhotos() {
    const gallery = document.getElementById('gallery');
    const popup = document.getElementById('popup');
    const selectedImage = document.getElementById('selectedImage');
    if (!gallery) return;

    try {
        const response = await fetch('/api/explorer');
        const items = await response.json();

        // Filter for images only
        galleryImages = items
            .filter(i => /\.(jpg|jpeg|png|gif|webp|heic)$/i.test(i.name))
            .map(i => ({
                name: i.name,
                url: `/api/raw?name=${encodeURIComponent(i.name)}`
            }));

        // Update count badge
        const badge = document.getElementById('photo-count');
        if (badge) badge.textContent = `${galleryImages.length} photo${galleryImages.length !== 1 ? 's' : ''}`;

        if (galleryImages.length === 0) {
            grid.innerHTML = `
                <div class="gallery-empty">
                    <i class="fas fa-images"></i>
                    <p>No photos found in your storage.</p>
                    <span>Upload your first photo to get started.</span>
                </div>`;
            return;
        }

        grid.innerHTML = galleryImages.map((img, index) => `
            <div class="photo-item" onclick="openGalleryPhoto(${index})">
                <img src="${img.url}" alt="${img.name}" loading="lazy">
                <div class="photo-overlay">
                    <div class="photo-name">${img.name}</div>
                    <div class="photo-item-actions">
                        <span class="photo-item-btn" title="View"><i class="fas fa-expand-alt"></i></span>
                        <a class="photo-item-btn" href="${img.url}" download="${img.name}" title="Download" onclick="event.stopPropagation()"><i class="fas fa-download"></i></a>
                    </div>
                </div>
            </div>`).join('');

    } catch (e) {
        grid.innerHTML = `
            <div class="gallery-empty">
                <i class="fas fa-exclamation-circle"></i>
                <p>Error connecting to your photo vault.</p>
            </div>`;
    }
}

function openGalleryPhoto(index) {
    const viewer = document.getElementById('gallery-viewer');
    if (!viewer || galleryImages.length === 0) return;

    currentPhotoIndex = index;
    _renderGalleryViewer();
    viewer.classList.add('active');
}

function _renderGalleryViewer() {
    const img    = galleryImages[currentPhotoIndex];
    const total  = galleryImages.length;

    const el = document.getElementById('viewer-img');
    el.classList.add('loading');
    el.onload = () => el.classList.remove('loading');
    el.src = img.url;
    el.alt = img.name;

    document.getElementById('viewer-filename').textContent = img.name;
    document.getElementById('viewer-counter').textContent  = `${currentPhotoIndex + 1} / ${total}`;
    document.getElementById('viewer-download').href        = img.url;
    document.getElementById('viewer-download').download    = img.name;

    document.getElementById('btn-prev').disabled = currentPhotoIndex === 0;
    document.getElementById('btn-next').disabled = currentPhotoIndex === total - 1;
}

function navigatePhoto(direction) {
    const next = currentPhotoIndex + direction;
    if (next < 0 || next >= galleryImages.length) return;
    currentPhotoIndex = next;
    _renderGalleryViewer();
}

function closeGalleryViewer() {
    const viewer = document.getElementById('gallery-viewer');
    if (viewer) {
        viewer.classList.remove('active');
        document.getElementById('viewer-img').src = '';
    }
}

// Keyboard navigation for gallery
document.addEventListener('keydown', (e) => {
    const viewer = document.getElementById('gallery-viewer');
    if (!viewer || !viewer.classList.contains('active')) return;
    if (e.key === 'Escape')      closeGalleryViewer();
    if (e.key === 'ArrowLeft')   navigatePhoto(-1);
    if (e.key === 'ArrowRight')  navigatePhoto(1);
});

function handleProUpload(files) {
    if (files.length === 0) return;
    const file = files[0];
    const container = document.getElementById('upload-progress-container');
    const fill = document.getElementById('progress-fill');
    const percentText = document.getElementById('progress-percent');

    container.style.display = 'block';
    document.getElementById('filename-display').innerText = file.name;

    const formData = new FormData();
    formData.append('file', file);

    const xhr = new XMLHttpRequest();
    xhr.open('POST', '/api/upload', true);

    xhr.upload.onprogress = (e) => {
        if (e.lengthComputable) {
            const percent = Math.round((e.loaded / e.total) * 100);
            fill.style.width = percent + '%';
            percentText.innerText = percent + '%';
        }
    };

    xhr.onload = () => {
        if (xhr.status === 200) {
            setTimeout(() => {
                container.style.display = 'none';
                loadFiles();
                if (document.getElementById('photo-grid')) loadPhotos();
            }, 1000);
        } else {
            alert("Upload failed. Check permissions.");
        }
    };
    xhr.send(formData);
}

// --- UTILS ---

function openViewer(url) {
    const viewer = document.getElementById('viewer');
    document.getElementById('viewer-img').src = url;
    viewer.classList.add('active');
}

function closeViewer() {
    document.getElementById('viewer').classList.remove('active');
    document.getElementById('viewer-img').src = '';
}

async function logout() {
    await fetch('/api/logout', { method: 'POST' });
    window.location.href = '/';
}