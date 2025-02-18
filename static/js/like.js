// Handle vote (like/dislike) functionality
async function handleVote(type, postId = null, commentId = null) {
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
            // Find the container for the post or comment
            const container = postId ?
                document.querySelector(`#post-${postId}`) :
                document.querySelector(`#comment-${commentId}`);

            if (!container) {
                console.error("container not found");
                return;
            }

            // Update like and dislike counts
            const likeCount = container.querySelector('.like-count');
            const dislikeCount = container.querySelector('.dislike-count');
            const likeBtn = container.querySelector('.like-btn');
            const dislikeBtn = container.querySelector('.dislike-btn');

            if (likeCount) likeCount.textContent = data.likeCount;
            if (dislikeCount) dislikeCount.textContent = data.dislikeCount;

            // Toggle active class based on user's vote
            if (data.userVote === 'like') {
                likeBtn.classList.add('active');
                dislikeBtn.classList.remove('active');
            } else if (data.userVote === 'dislike') {
                dislikeBtn.classList.add('active');
                likeBtn.classList.remove('active');
            } else {
                likeBtn.classList.remove('active');
                dislikeBtn.classList.remove('active');
            }
        } else {
            console.error("vote failed", data.error);
        }
    } catch (error) {
        console.log("error voting", error);
    }
}

// Toggle comment section visibility
function toggleCommentSection(commentBtn) {
    const commentSection = commentBtn.closest('.post-a').querySelector('.comments-section');
    if (commentSection) {
        commentSection.style.display = commentSection.style.display === 'none' ? 'block' : 'none';
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
            handleVote('like', postId);
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
            handleVote('dislike', postId);
        });
    });

    // Comment buttons
    document.querySelectorAll('.comment-btn').forEach(button => {
        button.addEventListener('click', () => {
            toggleCommentSection(button);
        });
    });
});