
//post actions
document.addEventListener('DOMContentLoaded', () => {
    const likeBtn = document.querySelector('.like-btn');
    const dislikeBtn = document.querySelector('.dislike-btn');
    const commentBtn = document.querySelector('.comment-btn');
    const commentsSection = document.querySelector('.comments-section');
    const commentCountSpan = document.querySelector('.comment-count');
  
    let commentCount = document.querySelectorAll('.comment').length;
    let likeCount = 1000; // fetch from server
    let dislikeCount = 2;
  
    const likeCountSpan = document.querySelector('.like-count');
    const dislikeCountSpan = document.querySelector('.dislike-count');
  
    // Initialize counts on load
    likeCountSpan.textContent = likeCount;
    dislikeCountSpan.textContent = dislikeCount;
    commentCountSpan.textContent = commentCount;
  
    let liked = false;
    let disliked = false;
  
    // Function to handle button activation
    function handleButtonActivation(activeButton, otherButtons) {
      activeButton.classList.toggle('active');
      otherButtons.forEach(button => button.classList.remove('active'));
    }
  
    function updateCounts() {
      likeCountSpan.textContent = likeCount;
      dislikeCountSpan.textContent = dislikeCount;
      commentCountSpan.textContent = commentCount;
    }
  
    function sendDataToServer() {
      const data = {
        likeCount: likeCount,
        dislikeCount: dislikeCount,
        commentCount: commentCount,
        liked: liked,
        disliked: disliked
      };
  
      fetch('/update_data', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(data)
        })
        .then(response => {
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }
          console.log('Data sent successfully to server!');
        })
        .catch(error => {
          console.error('There was an error sending data to the server:', error);
        });
    }
  
  
    likeBtn.addEventListener('click', () => {
        if (!isUserLoggedIn) {
            return;
        }
      if (liked) {
        likeBtn.classList.remove('active');
        likeCount--;
        liked = false;
      } else {
        if (disliked) {
          dislikeBtn.classList.remove('active');
          dislikeCount--;
          disliked = false;
        }
        likeBtn.classList.add('active');
        likeCount++;
        liked = true;
      }
      dislikeBtn.classList.remove('active');
      commentsSection.style.display = 'none'; // Hide comments
      updateCounts();
      sendDataToServer(); // Send data to the server
    });
  
    dislikeBtn.addEventListener('click', () => {
      if (disliked) {
        dislikeBtn.classList.remove('active');
        dislikeCount--;
        disliked = false;
      } else {
        if (liked) {
          likeBtn.classList.remove('active');
          likeCount--;
          liked = false;
        }
        dislikeBtn.classList.add('active');
        dislikeCount++;
        disliked = true;
      }
      likeBtn.classList.remove('active');
      commentsSection.style.display = 'none'; // Hide comments
      updateCounts();
      sendDataToServer(); // Send data to the server
    });
  
    commentBtn.addEventListener('click', () => {
      handleButtonActivation(commentBtn, [likeBtn, dislikeBtn]);
      commentsSection.style.display = commentBtn.classList.contains('active') ? 'block' : 'none'; //Show comments
    });
  
    // Function to handle comment like/dislike
    function setupCommentActions(commentElement) {
      const likeBtn = commentElement.querySelector('.comment-like-btn');
      const dislikeBtn = commentElement.querySelector('.comment-dislike-btn');
      const likeCountSpan = commentElement.querySelector('.comment-like-count');
      const dislikeCountSpan = commentElement.querySelector('.comment-dislike-count');
  
      let likeCount = 0;
      let dislikeCount = 0;
      let liked = false;
      let disliked = false;
  
      likeCountSpan.textContent = likeCount;
      dislikeCountSpan.textContent = dislikeCount;
  
      likeBtn.addEventListener('click', () => {
        if (liked) {
          likeBtn.classList.remove('active');
          likeCount--;
          liked = false;
        } else {
          if (disliked) {
            dislikeBtn.classList.remove('active');
            dislikeCount--;
            disliked = false;
          }
          likeBtn.classList.add('active');
          likeCount++;
          liked = true;
        }
        likeCountSpan.textContent = likeCount;
        dislikeCountSpan.textContent = dislikeCount;
        sendDataToServer();
      });
  
      dislikeBtn.addEventListener('click', () => {
        if (disliked) {
          dislikeBtn.classList.remove('active');
          dislikeCount--;
          disliked = false;
        } else {
          if (liked) {
            likeBtn.classList.remove('active');
            likeCount--;
            liked = false;
          }
          dislikeBtn.classList.add('active');
          dislikeCount++;
          disliked = true;
        }
        likeCountSpan.textContent = likeCount;
        dislikeCountSpan.textContent = dislikeCount;
        sendDataToServer();
      });
    }
  
    // Initialize comment like/dislike buttons for existing comments
    document.querySelectorAll('.comment').forEach(comment => {
      setupCommentActions(comment);
    });
  
    // Function to add a new comment
    function addComment(commentText) {
      const newComment = document.createElement('div');
      newComment.classList.add('comment');
      newComment.innerHTML = `
        <p>${commentText}</p>
        <div class="comment-actions">
          <button class="comment-like-btn">
            <svg class="icon" viewBox="0 0 24 24">
              <path d="M0 0h24v24H0z" fill="none"/>
              <path d="M1 21h4V9H1v12zm22-11c0-1.1-.9-2-2-2h-6.31l.95-4.57.03-.32c0-.41-.17-.79-.44-1.06L14.17 1 7.59 7.59C7.22 7.92 7 8.42 7 9v10c0 1.1.9 2 2 2h9c.83 0 1.54-.5 1.84-1.22l3.02-7.05c.09-.23.14-.47.14-.73v-1.91l-.01-.01L23 10z"/>
            </svg>
            (<span class="comment-like-count">0</span>)
          </button>
          <button class="comment-dislike-btn">
            <svg class="icon" viewBox="0 0 24 24">
              <path d="M0 0h24v24H0z" fill="none"/>
              <path d="M15 3H6c-.83 0-1.54.5-1.84 1.22l-3.02 7.05c-.09.23-.14.47-.14.73v1.91l.01.01L1 14c0 1.1.9 2 2 2h6.31l-.95 4.57-.03.32c0 .41.17.79.44 1.06L9.83 23l6.59-6.59c.36-.36.58-.86.58-1.41V5c0-1.1-.9-2-2-2zm4 0v12h4V3h-4z"/>
            </svg>
            (<span class="comment-dislike-count">0</span>)
          </button>
        </div>
      `;
      commentsSection.appendChild(newComment);
  
      setupCommentActions(newComment);
  
      commentCount++;
      updateCounts();
      sendDataToServer(); //send data to server
    }
  
    const addCommentBtn = document.querySelector('#add-comment-btn');
    const newCommentText = document.querySelector('#new-comment-text');
  
    addCommentBtn.addEventListener('click', () => {
      const commentText = newCommentText.value.trim();
      if (commentText !== '') {
        addComment(commentText);
        newCommentText.value = ''; // Clear the textarea
      }
    });
  
    // Example: Adding a new comment after 3 seconds
    setTimeout(() => {
      addComment("Another new comment!");
    }, 3000);
  
    // Initialize comment count on load
    updateCounts();
  });