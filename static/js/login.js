


/*
grab elements from html by their id
*/
const switchBtn = document.getElementById('switchBtn');
const switchText = document.getElementById('switchText');
const loginForm = document.getElementById('loginForm');
const signupForm = document.getElementById('signupForm');



/*
add an event listener on switch button, when button is clicked, it toggles betwen login and signup form
e.preventDefault() is to prevent page from loading
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
    } else {
        // switch to login form
        switchText.textContent = "Don't have an account?";
        switchBtn.textContent = 'Sign up';
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
            const password = document.getElementById('signupPassword').value;
            const confirmPassword = document.getElementById('signupConfirmPassword').value;
            

            // send password to backend together with email/username and see if they match
            if (password !== confirmPassword) {
                alert('Passwords do not match!');
                return;
            }
        }
        
        // Here you would typically send data to a server
        alert('Form submitted successfully!');

        // clear all fields after submission
        form.reset();
    });
});

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