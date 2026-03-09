document.addEventListener('DOMContentLoaded', () => {
    setupInputs();
    
    // Wire the verify form
    const verifyForm = document.getElementById('step-verify');
    if(verifyForm) {
        verifyForm.addEventListener('submit', (e) => {
            e.preventDefault();
            verifyCode();
        });
    }
});

function sendCode() {
    const btn = document.getElementById('btn-send');
    const subtitle = document.getElementById('verify-subtitle');
    
    btn.innerHTML = '<i class="fas fa-circle-notch fa-spin"></i> Sending...';
    btn.disabled = true;

    // Call backend to generate and print code
    fetch('/api/request-code', { method: 'POST' })
        .then(res => res.json())
        .then(data => {
            document.getElementById('step-send').style.display = 'none';
            document.getElementById('step-verify').style.display = 'block';
            subtitle.innerText = "Code sent successfully.";
            subtitle.style.color = "var(--text-main)";
            
            // Focus first input
            document.getElementById('code-1').focus();
        })
        .catch(err => {
            alert("Failed to send code.");
            btn.innerHTML = 'Send Code';
            btn.disabled = false;
        });
}

function verifyCode() {
    const p1 = document.getElementById('code-1').value;
    const p2 = document.getElementById('code-2').value;
    const p3 = document.getElementById('code-3').value;
    
    // Format code with dashes to match what backend generated
    const fullCode = p1 + "-" + p2 + "-" + p3;
    
    if (fullCode.length !== 12) {
        alert("Please enter a valid 10-digit code.");
        return;
    }

    const btn = document.querySelector('#step-verify button[type="submit"]');
    btn.innerHTML = '<i class="fas fa-circle-notch fa-spin"></i> Verifying...';
    btn.disabled = true;

    // Send code to backend to verify
    fetch('/api/submit-code', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code: fullCode })
    })
    .then(async res => {
        if (res.ok) {
            window.location.href = '/my_files';
        } else {
            const data = await res.json();
            alert(data.error || "Invalid verification code.");
            btn.innerHTML = 'Verify Device';
            btn.disabled = false;
        }
    })
    .catch(err => {
        alert("Verification connection error.");
        btn.innerHTML = 'Verify Device';
        btn.disabled = false;
    });
}

// Auto-advance cursor logic between the three inputs
function setupInputs() {
    const c1 = document.getElementById('code-1');
    const c2 = document.getElementById('code-2');
    const c3 = document.getElementById('code-3');
    
    if(!c1 || !c2 || !c3) return;

    // Restrict inputs to numbers only
    [c1, c2, c3].forEach(input => {
        input.addEventListener('input', function(e) {
            this.value = this.value.replace(/[^0-9]/g, '');
        });
    });

    c1.addEventListener('input', () => {
        if(c1.value.length === 3) c2.focus();
    });

    c2.addEventListener('input', () => {
        if(c2.value.length === 3) c3.focus();
    });
    
    // Backspace logic to jump back
    c2.addEventListener('keydown', (e) => {
        if(e.key === 'Backspace' && c2.value.length === 0) {
            c1.focus();
        }
    });

    c3.addEventListener('keydown', (e) => {
        if(e.key === 'Backspace' && c3.value.length === 0) {
            c2.focus();
        }
    });
}
