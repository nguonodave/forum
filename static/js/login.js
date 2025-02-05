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
        if (form.id === 'signupForm') {
            const email = document.getElementById('signupEmail').value;
            const password = document.getElementById('signupPassword').value;
            const username = document.getElementById('signupName').value;
            const confirmPassword = document.getElementById('signupConfirmPassword').value;
            
            // send password to backend together with email/username and see if they match
            if (password !== confirmPassword) {
                alert('Passwords do not match!');
                return;
            }
            signup(email,username,password).then();

        }else if(form.id === 'loginForm'){
            const usernameOrEmail = document.getElementById('loginEmailOrUsername').value;
            const password = document.getElementById('loginPassword').value;

            // if loginEmailOrUsername contains @ symbol means it is an email value
            // if it does not contain @ then it is the username
            if (usernameOrEmail.includes('@')) {
                login(usernameOrEmail,password,true).then();
            }else{
                login(usernameOrEmail,password,false).then();
            }
        }
        
        // Here you would typically send data to a server
        alert('Form submitted successfully!');

        // clear all fields after submission
        form.reset();
    });
});

async function login(emailOrUsername,password,isEmail){
    const body = isEmail
        ? {email:emailOrUsername, password:password}
        : {username: '', password: '', isEmail: false};

    try {
        let response = await fetch(`${BASE_URL}/login`, {
            method: 'POST',
            body: JSON.stringify(body),
            headers: {'Content-Type': 'application/json'}
        });

        let data = await response.json();
        console.log(data);
        if (response.ok){
            alert("login success");
        }else{
            alert('login error'+data.message);
        }

    }catch(error){
        alert("net error"+error.message);
    }
}

async function signup(email,username,password){
    const body = {email,username,password};
    try{
        let response = await fetch(`${BASE_URL}/signup`, {
            method: 'POST',
            body: JSON.stringify(body),
            headers: {'Content-Type': 'application/json'}
        });
        let data = await response.json();
        console.log("signup", data);
        if (response.ok){
            alert("signup success");
        }else{
            alert("error signing up"+data.message);
        }
    }catch(error){
        alert("net error"+error.message);
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