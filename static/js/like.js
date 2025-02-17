async function handleVote(type, postId = null, commentId = null) {
    console.log("Voting:", { type, postId, commentId }); // Debugging log

    try {
        const response = await fetch('/api/vote', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include', // Ensures cookies (session) are sent
            body: JSON.stringify({
                postId,
                commentId,
                type
            })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const data = await response.json();
        console.log("Vote response:", data); // Debugging log

        if (data.success) {
            // Locate the post or comment container
            const container = postId ? 
                document.querySelector(`#post-${postId}`) : 
                document.querySelector(`#comment-${commentId}`);

            if (!container) {
                console.error("Container not found!");
                return;
            }

            const likeCount = container.querySelector('.like-count');
            const dislikeCount = container.querySelector('.dislike-count');
            const likeBtn = container.querySelector('.like-btn');
            const dislikeBtn = container.querySelector('.dislike-btn');

            // Update counts
            likeCount.textContent = data.likeCount;
            dislikeCount.textContent = data.dislikeCount;

            // Update button states
            likeBtn.classList.toggle('active', data.userVote === 'like');
            dislikeBtn.classList.toggle('active', data.userVote === 'dislike');
        } else {
            console.error('Vote failed:', data.error);
        }
    } catch (error) {
        console.error('Error voting:', error);
    }
}
