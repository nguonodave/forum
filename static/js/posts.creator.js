const postsContainer = document.querySelector('.posts-container');
let page = 1; // Start with page 1
const limit = 10; // Number of posts per request
let isLoading = false; // Prevent multiple requests at the same time

// Function to fetch posts from the backend
async function fetchPosts() {
    if (isLoading) return; // Prevent duplicate requests
    isLoading = true;

    try {
        const response = await fetch(`/api/posts?page=${page}&limit=${limit}`);
        const posts = await response.json();

        if (posts.length === 0) {
            window.removeEventListener("scroll", handleScroll); // Stop loading if no more posts
            return;
        }

        posts.forEach(post => {
            const postElement = document.createElement("div");
            postElement.className = "post";
            postElement.innerHTML = `
                <h3>${post.title}</h3>
                <p>${post.content}</p>
            `;
            postsContainer.appendChild(postElement);
        });

        page++; // Move to the next page for future requests
    } catch (error) {
        console.error("Error loading posts:", error);
    } finally {
        isLoading = false;
    }
}

// Function to trigger fetching when user scrolls to bottom
function handleScroll() {
    if (window.innerHeight + window.scrollY >= document.body.offsetHeight - 100) {
        fetchPosts();
    }
}

// Load initial posts
fetchPosts();

// Attach scroll event listener
window.addEventListener("scroll", handleScroll);
