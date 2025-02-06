const switchBtn = document.getElementById('switchBtn');
const switchText = document.getElementById('switchText');
const loginForm = document.getElementById('loginForm');
const signupForm = document.getElementById('signupForm');
const formTitle = document.getElementById('form-title');
const BASE_URL = window.location.origin;
console.log(BASE_URL)

function toggleForm(isSignup) {
    if (isSignup) {
        loginForm.classList.add('hidden');
        signupForm.classList.remove('hidden');
        switchText.textContent = 'Already have an account?';
        switchBtn.textContent = 'Log in';
        formTitle.textContent = 'Sign up to forum';
    } else {
        loginForm.classList.remove('hidden');
        signupForm.classList.add('hidden');
        switchText.textContent = "Don't have an account?";
        switchBtn.textContent = 'Sign up';
        formTitle.textContent = 'Log in to forum';
    }
}

document.addEventListener('DOMContentLoaded', () => {
    toggleForm(window.location.pathname === '/register');
});

window.addEventListener('popstate', () => {
    toggleForm(window.location.pathname === '/register');
});

switchBtn.addEventListener('click', (e) => {
    e.preventDefault();
    const isSignup = switchBtn.textContent === 'Sign up';
    history.pushState(null, '', isSignup ? '/register' : '/login');
    toggleForm(isSignup);
});

/*
get all forms using querySelectorAll()
for each form, add a submit event listener
e.preventDefault() prevents form from submitting
*/
document.querySelectorAll('.form').forEach(form => {
    form.addEventListener('submit', (e) => {
        e.preventDefault();

        if (form.id === 'signupForm'){
            const email = document.getElementById('signupEmail').value;
            const password = document.getElementById('signupPassword').value;
            const username = document.getElementById('signupName').value;
            const confirmPassword = document.getElementById('signupConfirmPassword').value;
            // make sure username does not contain @ symbol
            if (username.includes('@')){
                alert("username cannot contain '@' symbol")
                return
            }
            if(password !== confirmPassword){
                alert('passwords do not match');
                return;
            }

            // send signup data to the backend
            signup(email,password,username).then(response => {
                if (response.ok) {
                    alert('signup success');
                    window.location.href = '/login';
                }else{
                    alert('signup failed' + response.message);
                }
            });

        }else if (form.id === 'loginForm'){
            const usernameOrEmail = document.getElementById('loginEmailOrUsername').value;
            const password = document.getElementById('loginPassword').value;

            // username should not contain email
            const isEmail = usernameOrEmail.includes('@');

            login(usernameOrEmail,password,isEmail).then(response => {
                if(response.ok){
                    alert('login success');
                    window.location.href = '/'; // redirect to homepage
                }else{
                    alert('login failed'+response.message);
                }
            });
        }
        // form.reset();
    });
})


async function login(usernameOrEmail, password, isEmail){
    const body = {}
    if (isEmail){
        body['email'] = usernameOrEmail;
    }else{
        body['username'] = usernameOrEmail;
    }
    body['password'] = password;
    console.log("login body",body);

    try {
        const response = await fetch('/login', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(body),
        });
        const data = await response.text();
        return { ok : response.ok , message: data.message }
    }catch(error){
        console.log(error);
        return { ok: false, message: error };
    }
}

async function signup(email, password, username) {
    const body = { email, password, username };

    try {
        const response = await fetch('/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body),
        });

        if (response.ok) {
            const data = await response.json();
            return { ok: true, message: data.message };
        } else {
            const errorText = await response.text();
            return { ok: false, message: errorText };
        }
    } catch (error) {
        console.log(">>>", error);
        return { ok: false, message: error.toString() };
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

document.querySelectorAll('.input-group input').forEach(input => {
    input.addEventListener('oninput', () => {
        if (input.value.length >= 1) {
            input.nextElementSibling.style.display = 'none'; // Hide the label
        } else {
            input.nextElementSibling.style.display = 'block'; // Show label if empty
        }
    });
});




// add password meter here
// remember to change url endpoint depending on user if login in or sign up
// add remember me checkbox
// add forgot password functionality
// connect to backend