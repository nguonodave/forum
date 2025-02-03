// Generate random posts
const postsContainer = document.querySelector('.posts-container');
for (let i = 0; i < 20; i++) {
    const post = document.createElement('div');
    post.className = 'post';
    post.innerHTML = `
                <h3>Post Title ${i + 1}</h3>
                <p>
                Post content will go here <br>
                Post content will go here <br>
                Post content will go here <br>
                </p>
            `;
    postsContainer.appendChild(post);
}
