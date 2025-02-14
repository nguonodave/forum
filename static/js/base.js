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
var isUserLoggedIn = document.body.getAttribute('data-is-logged-in') === 'true';

function handleCreatePost() {
  if (isUserLoggedIn) {
    openCreatePostDiv();
  } else {
    document.getElementById('loginPromptOverlay').style.display = 'flex';
  }
}
// Close login prompt overlay if the user clicks outside
document.getElementById('loginPromptOverlay').addEventListener('click', function(event) {
  if (event.target === this) {
    closeLoginPromptOverlay();
  }
});

function closeLoginPromptOverlay() {
  document.getElementById('loginPromptOverlay').style.display = 'none';
}

function openCreatePostDiv() {
  document.getElementById('createPostOverlay').style.display = 'flex';
}
