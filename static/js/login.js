const switchBtn = document.getElementById('switchBtn');
const switchText = document.getElementById('switchText');
const loginForm = document.getElementById('loginForm');
const signupForm = document.getElementById('signupForm');
const formTitle = document.getElementById('form-title');
const BASE_URL = window.location.origin;

// Toggle between login and signup forms
function toggleForm(isSignup) {
    if (isSignup) {
        loginForm.classList.add('hidden');
        signupForm.classList.remove('hidden');
        switchText.textContent = 'Already have an account?';
        switchBtn.textContent = 'Log in';
        formTitle.textContent = 'Sign up';
    } else {
        loginForm.classList.remove('hidden');
        signupForm.classList.add('hidden');
        switchText.textContent = "Don't have an account?";
        switchBtn.textContent = 'Sign up';
        formTitle.textContent = 'Log in';
    }
}

// Initialize form based on the current URL path
document.addEventListener('DOMContentLoaded', () => {
    toggleForm(window.location.pathname === '/register');
});

// Handle browser back/forward buttons
window.addEventListener('popstate', () => {
    toggleForm(window.location.pathname === '/register');
});

// Switch between login and signup forms
switchBtn.addEventListener('click', (e) => {
    e.preventDefault();
    const isSignup = switchBtn.textContent === 'Sign up';
    history.pushState(null, '', isSignup ? '/register' : '/login');
    toggleForm(isSignup);
});

// Handle form submissions
document.querySelectorAll('.form').forEach(form => {
    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        if (form.id === 'signupForm') {
            const email = document.getElementById('signupEmail').value;
            const password = document.getElementById('signupPassword').value;
            const username = document.getElementById('signupName').value;
            const confirmPassword = document.getElementById('signupConfirmPassword').value;

            // Validate username and password
            if (username.includes('@')) {
                showNotification("Username cannot contain '@' symbol",'error');
                return;
            }
            if (password !== confirmPassword) {
                showNotification('Passwords do not match',"error");
                return;
            }

            // Send signup data to the backend
            const response = await signup(email, password, username);
            if (response.ok) {
                showNotification('Signup successful','success');
                window.location.href = '/login';
            } else {
                console.log(response);
                showNotification('Signup failed: ' + response.message, 'error');
            }

        } else if (form.id === 'loginForm') {
            const usernameOrEmail = document.getElementById('loginEmailOrUsername').value;
            const password = document.getElementById('loginPassword').value;
            const isEmail = usernameOrEmail.includes('@');

            // Send login data to the backend
            const response = await login(usernameOrEmail, password, isEmail);
            if (response.ok) {
                showNotification('Login successful');
                setTimeout(() => {
                    window.location.href = '/';
                },500)
            } else {
                showNotification('Login failed: ' + response.message,'error');
            }
        }
    });
});

// Login function
async function login(usernameOrEmail, password, isEmail) {
    const body = isEmail ? { email: usernameOrEmail, password } : { username: usernameOrEmail, password };

    try {
        const response = await fetch('/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body),
        });
        const data = await response.json(); // Parse JSON response
        return { ok: response.ok, message: data.message };
    } catch (error) {
        console.error('Login error:', error);
        return { ok: false, message: 'An error occurred during login' };
    }
}

// Signup function
async function signup(email, password, username) {
    const body = { email, password, username };

    try {
        const response = await fetch('/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body),
        });
        const data = await response.json(); // Parse JSON response
        return { ok: response.ok, message: data.message };
    } catch (error) {
        console.error('Signup error:', error);
        return { ok: false, message: 'An error occurred during signup' };
    }
}

// Add input validation styling
document.querySelectorAll('input').forEach(input => {
    input.addEventListener('input', () => {
        if (input.value) {
            input.classList.add('filled');
        } else {
            input.classList.remove('filled');
        }
    });
});

function showNotification(message, type = 'success') {
    const notification = document.getElementById('notification');
    notification.textContent = message;
    notification.classList.remove('success', 'error', 'warning', 'info', 'show'); // Remove all previous classes
    notification.classList.add(type, 'show');

    setTimeout(() => {
        notification.classList.remove(type, 'show');
    }, 300);
}

const loginPassword = document.getElementById('loginPassword');
const signupPassword = document.getElementById('signupPassword');
const passwordMessage1 = document.querySelector("#passwordMessage1");
const passwordMessage2 = document.querySelector("#passwordMessage2");


// if (loginPassword) loginPassword.addEventListener('input', passwordVerifier);
// if (signupPassword) signupPassword.addEventListener('input', passwordVerifier);


if (loginPassword) loginPassword.addEventListener('input', (e) => {
    passwordVerifier(e, passwordMessage1);
});
if (signupPassword) signupPassword.addEventListener('input', (e) => {
    passwordVerifier(e, passwordMessage2);
});

function passwordVerifier(event, pass) {
    console.log(">>>>")
    let password = event.target.value;
    const output = [];
    password = password.trim()
    if (password.length > 0) {
        pass.style.display = 'block';
    }else{
        pass.style.display = 'none';
    }

    if (password.length < 8) {
        output.push("password must be at least 8 characters long");
    }
    if (password.includes('@')) {
        output.push("password cannot contain '@' symbol");
    }
    if (!/[A-Z]/.test(password)) {
        output.push("password must contain at least one UPPERCASE letter");
    }
    if (!/[a-z]/.test(password)) {
        output.push("password must contain at least one lowercase letter");
    }
    if (!/[0-9]/.test(password)) {
        output.push("password must contain at least one number");
    }
    if (!/[\W_]/.test(password)) {
        output.push("password must contain at least one special character (!, #, $)");
    }
    if (output.length > 0) {
        pass.innerHTML = output.join('<br>');
        pass.style.color = "red";
    } else {
        pass.innerHTML = "strong password!";
        pass.style.color = "green";
    }
}