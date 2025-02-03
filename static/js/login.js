/*
grab elements from html by their id
*/
const switchBtn = document.getElementById('switchBtn');
const switchText = document.getElementById('switchText');
const loginForm = document.getElementById('loginForm');
const signupForm = document.getElementById('signupForm');
const formTitle = document.getElementById('form-title');


/*
add an event listener on switch button, when button is clicked, it toggles between login and signup form
e.preventDefault() is to prevent page from re-loading
classList.toggle() switches/flips between the signup and login, displaying the correct info
*/
switchBtn.addEventListener('click', (e) => {
    e.preventDefault();
    loginForm.classList.toggle('hidden');
    signupForm.classList.toggle('hidden');
    
    // switch to signup form
    if (switchBtn.textContent === 'Sign up') {
        switchText.textContent = 'Already have an account?';
        switchBtn.textContent = 'Log in';
        formTitle.textContent = 'Sign up to forum';
    } else {
        // switch to log in form
        switchText.textContent = "Don't have an account?";
        switchBtn.textContent = 'Sign up';
        formTitle.textContent = 'Log in to forum'
    }
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
    const body = {};
    if (isEmail){
        body["email"] = emailOrUsername;
    }else{
        body["username"] = emailOrUsername;
    }
    body[password] = password;
    let data =  await fetch('http://localhost:8080/login',
        {
            method: 'POST',
            body: JSON.stringify(body),
            headers: {'Content-Type': 'application/json; charset=UTF-8'}
        }
    )
    data = await data.json();
    if (data.ok){
        alert("login successfully!");
    }else{
        alert("error login in"+data.message);
    }
}

async function signup(email,username,password){
    const payload = {
        "email":email,
        "password":password,
        "username":username,
        headers: {'Content-Type': 'application/json; charset=UTF-8'},
        method:'POST'
    }
    let data =  await fetch('http://localhost:8080/signup', payload)
    if (data.ok){
        alert("signup successfully!");
    }else{
        alert("error login in"+data.message);
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

// add password meter here
// remember to change url endpoint depending on user if loggin in or sign up
// add remember me checkbox
// add forgot password functionality
// connect to backend