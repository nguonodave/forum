async function handleVote(type, postId = null, commentId = null) {
    try {
        const response = await fetch('/api/vote', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                postId,
                commentId,
                type
            })
        });
        
        const data = await response.json();
        if (data.success) {
            // Update UI elements
            const container = postId ? 
                document.querySelector(`#post-${postId}`) : 
                document.querySelector(`#comment-${commentId}`);
                
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