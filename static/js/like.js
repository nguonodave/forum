// Handle vote (like/dislike) functionality
async function handleVote(type, postId = null, commentId = null,content = null) {
    console.log("voting", { type, postId, commentId });
    if (!postId && !commentId) {
        console.log("both post id and comment id are missing");
        return;
    }

    try {
        const response = await fetch("/api/vote", {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ postId: postId, commentId: commentId, type: type }),
        });

        if (!response.ok) {
            throw new Error(`HTTP ERROR ${response.status}`);
        }

        const data = await response.json();
        console.log("vote response", data);

        if (data.success) {
            console.log(data)
            function changeLikeCounts(className, newCount,span) {
                const targets = document.querySelectorAll(`.${className}`);
                targets.forEach((target) => {
                    if (target.getAttribute('data-post-id') === postId) {
                        console.log("Found element:", target);
                        const targetA = target.querySelector(`.${span}`);
                        console.log("another A",targetA)
                        targetA.innerText = String(newCount);
                    }
                });
            }

            changeLikeCounts('like-btn', data.likeCount,'span-x')
            changeLikeCounts('dislike-btn', data.dislikeCount,'span-x')
            changeLikeCounts('comment-btn', data.likeCount,'comment-like-count');
            changeLikeCounts('dislike-btn', data.dislikeCount,'comment-dislike-count');
            // Toggle active class based on user's vote
            if (data.userVote === 'like') {
                likeBtn.classList.add('active');
                dislikeBtn.classList.remove('active');
                likeBtn.display.color = 'green';
            } else if (data.userVote === 'dislike') {
                dislikeBtn.classList.add('active');
                likeBtn.classList.remove('active');
                dislikeBtn.display.color = 'red';
            } else {
                likeBtn.classList.remove('active');
                dislikeBtn.classList.remove('active');
            }
        } else {
            showNotification('vote failed');
            console.error("vote failed", data.error);
        }
    } catch (error) {
        console.log("error voting", error);
    }
}

function toggleCommentSection(button) {
    const commentSection = button.closest('#post-a').querySelector('.comments-section');
    console.log("commentSection", commentSection);
    if (commentSection) {
        // Toggle visibility
        commentSection.style.display = (commentSection.style.display === 'none' || commentSection.style.display === '') ? 'block' : 'none';

        // Toggle button text
        const isExpanded = button.getAttribute('aria-expanded') === 'true';
        console.log("is expanded",isExpanded)
        button.setAttribute('aria-expanded', !isExpanded);
    }
}

// Attach event listeners to like, dislike, and comment buttons
document.addEventListener('DOMContentLoaded', () => {
    // Like buttons
    document.querySelectorAll('.like-btn').forEach(button => {
        button.addEventListener('click', () => {
            const postId = button.getAttribute('data-post-id');
            if (!postId) {
                console.error("post id for like is null or empty");
                return;
            }
            handleVote('like', postId,null);
        });
    });

    // Dislike buttons
    document.querySelectorAll('.dislike-btn').forEach(button => {
        button.addEventListener('click', () => {
            const postId = button.getAttribute('data-post-id');
            if (!postId) {
                console.error("post id for dislike is null or empty");
                return;
            }
            handleVote('dislike', postId,null);
        });
    });

    // like buttons for comments
    document.querySelectorAll('.like-btn[data-comment-id]').forEach(button => {
        button.addEventListener('click', () => {
            const commentId = button.getAttribute('data-comment-id');
            console.log("comment id ???",commentId);
            handleVote('like', null, commentId);
        });
    });

    // dislike buttons for comments
    document.querySelectorAll('.dislike-btn[data-comment-id]').forEach(button => {
        button.addEventListener('click', () => {
            const commentId = button.getAttribute('data-comment-id');
            handleVote('dislike', null, commentId);
        });
    });

    //  comment
    document.querySelectorAll('#add-comment-btn').forEach(button => {
        button.addEventListener('click', () => {
            const postId = button.getAttribute('data-post-id');
            const inputField =  button.previousElementSibling; // input field
            const commentMessage = inputField.value.trim(); // remove white spaces
            if(commentMessage){
                handleComment(postId,commentMessage,button,inputField);
            }
        })
    })
    // Comment buttons
    document.querySelectorAll('.comment-btn').forEach(button => {
        button.addEventListener('click', () => {
            console.log("comment btn clicked <<")
            toggleCommentSection(button);
        });
    });
});

async function handleComment(postId, comment, button) {
    console.log('handleComment called:', postId, comment);

    const inputField = button.previousElementSibling;

    try {
        const response = await fetch('/api/comment', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ post_id: postId, content: comment }),
        });

        const data = await response.json();
        console.log("Response from backend:", data);

        if (response.ok) {
            // Find the correct comment section for this post
            const postContainer = button.closest('.post-content');
            const commentSection = postContainer.querySelector('.comment-section');

            // Create a new comment element
            const newComment = document.createElement('div');
            newComment.classList.add('comment');
            newComment.id = `comment-${data.id}`;
            newComment.innerHTML = `
                <p>${data.content}</p>
                <div class="comment-actions">
                    <button class="like-btn" data-comment-id="${data.id}">
                        👍 <span class="comment-like-count">0</span>
                    </button>
                    <button class="dislike-btn" data-comment-id="${data.id}">
                        👎 <span class="comment-dislike-count">0</span>
                    </button>
                </div>
            `;

            // Append new comment to the section
            commentSection.appendChild(newComment);
            commentSection.style.display = 'block'; // Show comment section

            // Clear input field
            inputField.value = "";

            // Update comment count
            const commentCountSpan = postContainer.querySelector('.comment-btn .comment-count');
            if (commentCountSpan) {
                commentCountSpan.textContent = parseInt(commentCountSpan.textContent, 10) + 1;
            }

        } else {
            console.error("Failed to add comment:", data.error);
        }
    } catch (error) {
        console.error("Error:", error);
    }
}