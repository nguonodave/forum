// Setup responsive design
function _setup() {
    // Handle mobile navigation
    const appIcon = document.querySelector('.app-icon');
    const leftSidebar = document.querySelector('.left-sidebar');
    const overlay = document.querySelector('.overlay');

    appIcon.addEventListener('click', () => {
        leftSidebar.classList.toggle('open');
        overlay.style.display = leftSidebar.classList.contains('open') ? 'block' : 'none';
    });

    overlay.addEventListener('click', () => {
        leftSidebar.classList.remove('open');
        overlay.style.display = 'none';
    });

    // Handle header scroll behavior on mobile
    let lastScroll = 0;
    const header = document.querySelector('.header');

    window.addEventListener('scroll', () => {
        if (window.innerWidth <= 575) {
            const currentScroll = window.pageYOffset;

            if (currentScroll > lastScroll && currentScroll > 60) {
                header.classList.add('hidden');
            } else {
                header.classList.remove('hidden');
            }

            lastScroll = currentScroll;
        }
    });

    // Handle window resize
    window.addEventListener('resize', () => {
        if (window.innerWidth > 575) {
            overlay.style.display = 'none';
            header.classList.remove('hidden');
        }
    });

}

document.addEventListener('DOMContentLoaded', _setup);
// Access the body tag where the data-is-logged-in attribute is set
const isUserLoggedIn = document.body.getAttribute('data-is-logged-in') === 'true';
console.log("is user logged in ?",isUserLoggedIn)


//Prevent user from performing post actions if not logged in
const actionButtons = document.querySelectorAll(".like-btn, .dislike-btn, #new-comment-text, #add-comment-btn, .comment-like-btn, .comment-dislike-btn");

actionButtons.forEach((button) => {
  button.addEventListener("click", function (event) {
    if(!isUserLoggedIn){
      event.preventDefault();
      document.getElementById('loginPromptOverlay').style.display = 'flex';
    }
  })
})


function handleCreatePost() {
  if (isUserLoggedIn) {
    openCreatePostDiv();
  } else {
    document.getElementById('loginPromptOverlay').style.display = 'flex';
  }
}

function openCreatePostDiv() {
  document.getElementById('createPostOverlay').style.display = 'flex';
}

function closeCreatePostOverlay() {
  document.getElementById('createPostOverlay').style.display = 'none';
}

const loginPromptOverlay = document.getElementById('loginPromptOverlay');
const createPostOverlay = document.getElementById('createPostOverlay');

function closeLoginPromptOverlay() {
  document.getElementById('loginPromptOverlay').style.display = 'none';
}

if (createPostOverlay) {
  createPostOverlay.addEventListener('click', function (event) {
    if (event.target === this) {
      closeCreatePostOverlay();
    }
  });
}

if (loginPromptOverlay) {
  loginPromptOverlay.addEventListener('click', function (event) {
    if (event.target === this) {
      closeLoginPromptOverlay();
    }
  });
}
const logoutBtn = document.getElementById('logout-btn');
if (logoutBtn) {
    console.log(logoutBtn)
    logoutBtn.addEventListener('click',logout);
}
async function logout() {
    try{
        const response = await fetch('/logout', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
        });
        if (response.ok){
            const data = await response.json();
            window.location.href = '/';//redirect to home
        }else{
            const errData = await response.json();
        }
    }catch(error){
        console.log(error);
    }
}

const likeBtn  = document.getElementById('like-btn');
likeBtn.addEventListener('click', async () => {
    try {
        const response = await fetch('/api/vote', {
            method:'POST',
            headers: {'Content-Type': 'application/json'},
        });
        if (response.ok){
        }else{
            const errData = await response.json();
        }
    }catch(error){
        console.log(error);
    }
})

const dislikeBtn = document.getElementById('dislike-btn');
dislikeBtn.addEventListener('click', async () => {
    try {
        const response = await fetch('/api/vote', {
            method:'POST',
            headers: {'Content-Type': 'application/json'},
        });
        if (response.ok){
        }else{
            const errData = await response.json();
        }
    }catch(error){
        console.log(error);
    }
})