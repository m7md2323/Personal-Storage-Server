/**
 * MyCloud Pro - Unified Frontend Logic
 * Version 4.0: Modernized & Simplified Build
 */

document.addEventListener('DOMContentLoaded', () => {
    // 1. Identify current page context
    const hasUserList = !!document.getElementById('user-list');
    const hasFileList = !!document.getElementById('file-list');
    const hasGallery = !!document.getElementById('gallery');

    // 2. Execute necessary data fetchers
    if (hasUserList) fetchUsers();
    if (hasFileList) loadFiles();
    if (hasGallery) loadPhotos();

    // 3. Setup core event listeners
    setupAuthListeners();
    setupUploadZone();
});

// ==========================================
// AUTHENTICATION & USER SELECTION
// ==========================================
let selectedUser = "";

function setupAuthListeners() {
    // Login Form
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }

    // Create User Form
    const createForm = document.getElementById('create-user-form');
    if (createForm) {
        createForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const username = document.getElementById('new-username').value;
            const password = document.getElementById('new-password').value;
            const btn = createForm.querySelector('button');
            const originalText = btn.innerHTML;

            btn.innerHTML = '<i class="fas fa-circle-notch fa-spin"></i> Creating...';
            btn.disabled = true;

            try {
                const res = await fetch('/api/create_user', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username, password })
                });

                if (res.ok) {
                    window.location.href = '/'; 
                } else {
                    const data = await res.json();
                    alert("Error: " + (data.error || "Could not create vault"));
                }
            } catch (err) {
                alert("Connection failed.");
            } finally {
                btn.innerHTML = originalText;
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
        const users = await response.json(); 

        if (users.length === 0) {
            list.innerHTML = `
                <div style="grid-column: 1 / -1; text-align: center; color: var(--text-light); padding: 20px;">
                    No vaults found. Create one below!
                </div>`;
            return;
        }

        list.innerHTML = users.map(user => {
            const name = typeof user === 'object' ? user.username : user;
            return `
                <div class="user-avatar" onclick="showPasswordInput('${name}')">
                    <i class="fas fa-fingerprint avatar-icon"></i>
                    <div class="user-name">${name}</div>
                </div>`;
        }).join('');
    } catch (err) {
        list.innerHTML = '<p style="color:red;">Error loading vaults.</p>';
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
        loginForm.style.animation = 'fadeUp 0.4s forwards';
        document.getElementById('selected-user-text').innerText = "Vault: " + username;
        document.getElementById('password').focus();
    }, 200);
}

async function handleLogin(e) {
    e.preventDefault(); 
    const password = document.getElementById('password').value;
    const btn = e.target.querySelector('button[type="submit"]');
    const originalText = btn.innerHTML;
    
    btn.innerHTML = '<i class="fas fa-circle-notch fa-spin"></i> Unlocking...';
    btn.disabled = true;

    try {
        const res = await fetch('/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username: selectedUser, password: password })
        });

        if (res.ok) {
            window.location.href = '/my_files';
        } else {
            alert("Invalid password");
            document.getElementById('password').value = '';
            document.getElementById('password').focus();
        }
    } catch (e) {
        alert("Connection error.");
    } finally {
        btn.innerHTML = originalText;
        btn.disabled = false;
    }
}

async function logout() {
    try {
        await fetch('/api/logout', { method: 'POST' });
    } catch(e) {}
    window.location.href = '/';
}

// ==========================================
// FILE EXPLORER
// ==========================================

async function loadFiles() {
    const list = document.getElementById('file-list');
    try {
        const response = await fetch('/api/explorer'); 
        if (!response.ok) throw new Error("Fetch failed");
        
        const items = await response.json();
        list.innerHTML = '';

        if (!items || items.length === 0) {
            list.innerHTML = `
                <div class="loading-state" style="padding: 40px;">
                    <i class="fas fa-folder-open" style="margin-bottom: 12px; color: var(--border-solid);"></i>
                    <p>Vault is empty</p>
                </div>`;
            return;
        }

        items.forEach(item => {
            const rawUrl = `/api/raw?name=${encodeURIComponent(item.name)}`;
            const ext = item.name.split('.').pop().toLowerCase();
            const isPhoto = ['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(ext);
            const isDoc = ['pdf', 'txt', 'log'].includes(ext);

            let actions = "";
            if (!item.isDir) {
                if (isPhoto) {
                    actions += `<button class="action-btn" onclick="openSingleViewer('${rawUrl}', '${item.name}')" title="View"><i class="fas fa-eye"></i></button>`;
                }
                if (isDoc) {
                    actions += `<a href="${rawUrl}" target="_blank" class="action-btn" title="Open"><i class="fas fa-external-link-alt"></i></a>`;
                }
                actions += `<a href="${rawUrl}" download class="action-btn" title="Download"><i class="fas fa-download"></i></a>`;
                actions += `<button class="action-btn" style="color: var(--danger);" onclick="deleteFile('${item.name}')" title="Delete"><i class="fas fa-trash-alt"></i></button>`;
            }

            const iconClass = item.isDir ? 'fa-folder' : (isPhoto ? 'fa-image' : 'fa-file-alt');
            
            list.innerHTML += `
                <div class="file-row">
                    <div class="file-info">
                        <span class="file-icon"><i class="fas ${iconClass}"></i></span> 
                        ${item.name}
                    </div>
                    <div class="file-actions" style="display:flex; gap:8px;">${actions}</div>
                </div>`;
        });
    } catch (e) {
        list.innerHTML = '<div class="loading-state"><p style="color:var(--danger)">Error accessing secure vault.</p></div>';
    }
}

// ==========================================
// PHOTO GALLERY
// ==========================================

let galleryImages = [];
let currentPhotoIndex = 0;

async function loadPhotos() {
    const gallery = document.getElementById('gallery');
    if (!gallery) return;

    try {
        const response = await fetch('/api/explorer');
        if (!response.ok) throw new Error("Fetch failed");
        
        const items = await response.json();

        // Filter valid images
        galleryImages = items
            .filter(i => /\.(jpg|jpeg|png|gif|webp|heic)$/i.test(i.name))
            .map(i => ({
                name: i.name,
                url: `/api/raw?name=${encodeURIComponent(i.name)}`
            }));

        const badge = document.getElementById('photo-count');
        if (badge) badge.textContent = `${galleryImages.length} items`;

        if (galleryImages.length === 0) {
            gallery.innerHTML = `
                <div class="gallery-empty">
                    <i class="fas fa-image" style="opacity: 0.5;"></i>
                    <p>No photos securely stored.</p>
                    <span>Upload media in the Files tab to see them here.</span>
                </div>`;
            return;
        }

        gallery.innerHTML = galleryImages.map((img, index) => `
            <div class="photo-item" onclick="openGalleryPhoto(${index})">
                <img src="${img.url}" alt="${img.name}" loading="lazy">
                <div class="photo-overlay">
                    <div class="photo-name">${img.name}</div>
                    <div class="photo-item-actions">
                        <span class="photo-item-btn" title="Expand"><i class="fas fa-expand"></i></span>
                        <a class="photo-item-btn" href="${img.url}" download="${img.name}" title="Download" onclick="event.stopPropagation()"><i class="fas fa-download"></i></a>
                        <button class="photo-item-btn" style="color: var(--danger);" title="Delete" onclick="event.stopPropagation(); deleteFile('${img.name}')"><i class="fas fa-trash-alt"></i></button>
                    </div>
                </div>
            </div>`).join('');

    } catch (e) {
        gallery.innerHTML = `
            <div class="gallery-empty">
                <i class="fas fa-exclamation-triangle" style="color: var(--danger)"></i>
                <p>Failed to decrypt media vault.</p>
            </div>`;
    }
}

// Lightbox Logic for Gallery
function openGalleryPhoto(index) {
    const viewer = document.getElementById('gallery-viewer');
    if (!viewer || galleryImages.length === 0) return;

    currentPhotoIndex = index;
    _renderGalleryViewer();
    viewer.classList.add('active');
}

function _renderGalleryViewer() {
    const img    = galleryImages[currentPhotoIndex];
    const el     = document.getElementById('viewer-img');
    
    el.style.opacity = '0.3';
    el.onload = () => el.style.opacity = '1';
    el.src = img.url;
    el.alt = img.name;

    document.getElementById('viewer-filename').textContent = img.name;
    document.getElementById('viewer-counter').textContent  = `${currentPhotoIndex + 1} / ${galleryImages.length}`;
    document.getElementById('viewer-download').href        = img.url;
    document.getElementById('viewer-download').download    = img.name;

    document.getElementById('btn-prev').disabled = currentPhotoIndex === 0;
    document.getElementById('btn-next').disabled = currentPhotoIndex === galleryImages.length - 1;
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

// Universal Keyboard Nav
document.addEventListener('keydown', (e) => {
    const galViewer = document.getElementById('gallery-viewer');
    if (galViewer && galViewer.classList.contains('active')) {
        if (e.key === 'Escape')      closeGalleryViewer();
        if (e.key === 'ArrowLeft')   navigatePhoto(-1);
        if (e.key === 'ArrowRight')  navigatePhoto(1);
    }
    
    const singleViewer = document.getElementById('viewer');
    if (singleViewer && singleViewer.classList.contains('active')) {
        if (e.key === 'Escape') closeViewer();
    }
});

// Single File Viewer (from file explorer)
function openSingleViewer(url, name) {
    const viewer = document.getElementById('viewer');
    const img = document.getElementById('viewer-img');
    const filenameLabel = document.getElementById('viewer-filename-single');
    const dlLink = document.getElementById('viewer-download-single');
    
    if(filenameLabel) filenameLabel.textContent = name || 'View File';
    if(dlLink) {
        dlLink.href = url;
        if(name) dlLink.download = name;
    }
    
    if(img) img.src = url;
    if(viewer) viewer.classList.add('active');
}

function closeViewer() {
    const viewer = document.getElementById('viewer');
    if(viewer) {
        viewer.classList.remove('active');
        const img = document.getElementById('viewer-img');
        if(img) img.src = '';
    }
}

// ==========================================
// UPLOADER
// ==========================================

function setupUploadZone() {
    const zone = document.getElementById('upload-zone');
    const input = document.getElementById('file-input');
    if (!zone || !input) return;

    // Click to select
    zone.addEventListener('click', () => input.click());

    // File selected via dialog
    input.addEventListener('change', (e) => {
        if (e.target.files.length) handleUpload(e.target.files[0]);
        input.value = ''; // Reset
    });

    // Drag and Drop
    zone.addEventListener('dragover', (e) => {
        e.preventDefault();
        zone.classList.add('dragover');
    });

    zone.addEventListener('dragleave', () => {
        zone.classList.remove('dragover');
    });

    zone.addEventListener('drop', (e) => {
        e.preventDefault();
        zone.classList.remove('dragover');
        if (e.dataTransfer.files.length) handleUpload(e.dataTransfer.files[0]);
    });
}

function handleUpload(file) {
    if (!file) return;
    
    const container = document.getElementById('upload-progress-container');
    const fill = document.getElementById('progress-fill');
    const percentText = document.getElementById('progress-percent');
    
    if (!container || !fill || !percentText) {
        // Fallback upload without visual progress bar
        submitFile(file);
        return;
    }

    container.style.display = 'block';
    document.getElementById('filename-display').innerText = file.name;

    submitFile(file, (percent) => {
        fill.style.width = percent + '%';
        percentText.innerText = percent + '%';
    });
}

function submitFile(file, onProgress) {
    const formData = new FormData();
    formData.append('file', file);

    const xhr = new XMLHttpRequest();
    xhr.open('POST', '/api/upload', true);

    xhr.upload.onprogress = (e) => {
        if (e.lengthComputable && onProgress) {
            onProgress(Math.round((e.loaded / e.total) * 100));
        }
    };

    xhr.onload = () => {
        if (xhr.status === 200) {
            setTimeout(() => {
                const container = document.getElementById('upload-progress-container');
                if(container) container.style.display = 'none';
                
                // Refresh views
                if (document.getElementById('file-list')) loadFiles();
                if (document.getElementById('gallery')) loadPhotos();
            }, 800);
        } else {
            alert("Upload failed. Disconnected or vault locked.");
            const container = document.getElementById('upload-progress-container');
            if(container) container.style.display = 'none';
        }
    };
    
    xhr.onerror = () => {
        alert("Upload connection error.");
        const container = document.getElementById('upload-progress-container');
        if(container) container.style.display = 'none';
    }

    xhr.send(formData);
}

// ==========================================
// FILE DELETION
// ==========================================

async function deleteFile(filename) {
    if (!confirm(`Are you sure you want to completely delete "${filename}"? This cannot be undone.`)) {
        return;
    }

    try {
        const res = await fetch('/api/delete', {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: filename })
        });

        if (res.ok) {
            // Refresh the current view
            if (window.location.pathname.includes('/photos')) {
                loadPhotos();
            } else {
                loadFiles();
            }
        } else {
            const err = await res.json();
            alert(err.error || "Failed to delete file.");
        }
    } catch (e) {
        alert("Connection error occurred while deleting.");
    }
}