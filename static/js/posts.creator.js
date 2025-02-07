// Select the posts container in your HTML
const postsContainer = document.querySelector('.posts-container');

// Fetch the paginated posts data from your server
function fetchPosts(page = 1, limit = 10) {
  fetch(`/api/posts?page=${page}&limit=${limit}`)
    .then((response) => response.json())
    .then((postsData) => {
      // Render the posts after receiving the data
      postsContainer.innerHTML = ''; // Clear the container before rendering new posts
      postsData.forEach((postData) => {
        console.log(postData)
        const post = document.createElement('div');
        post.className = 'post';

        post.innerHTML = `
          <h3>${postData.Title}</h3>
          <p>${postData.Content}</p>
          <div class="post-footer">
            <span>Category: ${postData.Category.Name}</span>
            <span>Votes: ${postData.Votes}</span>
            <span>Posted by: ${postData.User.Username}</span>
            <span>Posted on: ${new Date(postData.CreatedAt).toLocaleDateString()}</span>
          </div>
        `;

        // Append the post to the posts container
        postsContainer.appendChild(post);
      });
    })
    .catch((error) => {
      console.error('Error fetching posts:', error);
    });
}

// Call the fetchPosts function to load posts on page load
fetchPosts();
