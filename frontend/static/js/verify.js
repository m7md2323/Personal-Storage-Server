// 1. Ask the server to print the code
async function requestCode() {
    try {
        const response = await fetch('/api/request-code', { method: 'POST' });
        
        if (response.ok) {
            // Hide the request button, show the input box
            document.getElementById('request-section').style.display = 'none';
            document.getElementById('submit-section').style.display = 'block';
        } else {
            alert("Failed to request code from server.");
        }
    } catch (error) {
        console.error("Error:", error);
    }
}

// 2. Send the typed code back to the server
async function submitCode() {
    const codeInput = document.getElementById('access-code').value;
    
    try {
        const response = await fetch('/api/submit-code', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ code: codeInput })
        });

        if (response.ok) {
            // Success! The server gave us the cookie. Go to the Netflix screen!
            window.location.href = "/";
        } else {
            // Wrong code
            document.getElementById('error-msg').style.display = 'block';
        }
    } catch (error) {
        console.error("Error:", error);
    }
}