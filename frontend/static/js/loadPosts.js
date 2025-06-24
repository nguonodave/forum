import { notify } from "./pageContext/notification.js";
export let posts;
export const loadPosts = async (ws) => {
    const token = localStorage.getItem("sessionToken")
    try {
        const postsRes = await fetch(`/home?token=${token}`);

        
        const postData = await postsRes.json();
        posts = postData.posts;
        
        if (!postsRes.ok) {
            notify(postData.message, "var(--danger-color)");
            setTimeout(() => {
                localStorage.clear();
                window.location.href = "/welcome";
            }, 3000);
            return;
            
        }
        const postsContainer = document.querySelector(".posts-section");
        if (!posts || posts.length === 0) {
            // console.log("No posts available");
            postsContainer.innerHTML = `
            <h4 style="color:white; text-align:center;">
            No posts available
            </h4>`;
            return;
        }

        let pageData = "";

        posts.forEach((post) => {
            const commentsSection = Array.isArray(post.comments)
                ? post.comments.map(comment => `
                    <div class="comment">
                        <div class="avatar">
                        <i class="fa-regular fa-user"></i>
                        </div>
                        <div class="comment-content">
                            <div class="comment-header">
                                <h4>${comment.user_name}</h4>
                                <span class="comment-time">${comment.created_at}</span>
                            </div>
                            <p class="comment-text">${comment.content}</p>
                         <!---   <form class="commentLikeForm">
                                <input type="hidden" name="commentid" value="${comment.id}">
                                <input type="hidden" name="reaction" value="like">
                                <button class="like-btn"><i class="fas fa-thumbs-up"></i> ${comment.likes}</button>
                            </form>
                            <form class="commentDislikeForm">
                                <input type="hidden" name="commentid" value="${comment.id}">
                                <input type="hidden" name="reaction" value="dislike">
                                <button class="like-btn"><i class="fas fa-thumbs-down"></i> ${comment.dislikes}</button>
                            </form> --->
                        </div>
                    </div>
                `).join("")
                : "";

            pageData += `
                <div class="post" postid=${post.id}>
                    <div class="post-header">
                        <div class="avatar">
                        <i class="fa-regular fa-user"></i>
                        </div>
                        <div class="post-info">
                            <h4>${post.name}</h4>
                            <span class="post-time">${post.created_at}</span>
                            <span class="post-category"><span style="font-weight:bold;">Category(ies):</span> ${post.category}</span>
                        </div>
                    </div>
                    <div class="post-content">
                        <h6>${post.title}</h6>
                        <p>${post.content}</p>
                    </div>
                    <div class="post-actions">
                        <form id="postLikesForm">
                            <input type="hidden" name="postid" value="${post.id}">
                            <input type="hidden" name="reaction" value="like">
                            <button class="like-btn"><i class="fas fa-thumbs-up"></i><span id="likes-count">${post.likes}</span></button>
                        </form>
                        <form id="postDislikeBtn">
                            <input type="hidden" name="postid" value="${post.id}">
                            <input type="hidden" name="reaction" value="dislike">
                            <button class="like-btn"><i class="fas fa-thumbs-down"></i> <span id="dislike-count">${post.dislikes}</span></button>
                        </form>
                        <button class="comment-btn"><i class="fas fa-comment"></i>${post.comments_num}</button>
                    </div>
                    <div class="comments-section">
                        ${commentsSection}
                        <div class="add-comment">
                            <div class="avatar">
                                <img src="https://avatars.dicebear.com/api/adventurer/johndoe.svg" alt="image">
                            </div>
                            <form class="addCommentForm">
                                <input type="hidden" name="identity" value="comment">
                                <input type="hidden" name="postId" value="${post.id}">
                                <input type="text" placeholder="Write a comment..." name="comment" required>
                                <button type="submit"><i class="fas fa-paper-plane"></i></button>
                            </form>
                        </div>
                    </div>
                </div>
            `;
        });

        postsContainer.innerHTML = pageData;

    } catch (error) {
        console.error("Error loading posts:", error);
    }

    const postsSection = document.querySelector(".posts-section");

    postsSection.replaceWith(postsSection.cloneNode(true));
    const newPostsSection = document.querySelector(".posts-section");

    newPostsSection.addEventListener("submit", async (e) => {
        e.preventDefault();

        const form = e.target;

        // Handle Likes and Dislikes for posts
        if (form.matches("#postLikesForm, #postDislikeBtn")) {
            const id = form.querySelector("input[name='postid']").value;
            const reaction = form.querySelector("input[name='reaction']").value;
            const endpoint = reaction === "like" ? "/likePost" : "/dislikePost";

            try {
                const res = await fetch(endpoint, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ id, reaction }),
                });

                if (res.ok) {
                    const resBody = await res.json();
                    const reactions = {
                        type: "reaction",
                        data: { ...resBody, id, reaction },
                    };

                    if (ws && ws.readyState === WebSocket.OPEN) {
                        ws.send(JSON.stringify(reactions));
                    }
                } else {
                    console.error(`Failed to ${reaction} post`, res.status);
                }
            } catch (error) {
                console.error(`Error on ${reaction}:`, error);
            }
        }

        if (form.matches(".commentLikeForm, .commentDislikeForm")) {
            const commentId = form.querySelector('input[name="commentid"]').value;
            const reaction = form.querySelector('input[name="reaction"]').value;
            const endpoint = reaction === "like" ? "/likeComment" : "/dislikeComment";

            try {
                const res = await fetch(endpoint, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ commentId, reaction }),
                });

                if (res.ok) {
                    const resBody = await res.json();
                    console.log(`Comment ${reaction} successful:`, resBody);
                } else {
                    console.error(`Failed to ${reaction} comment`, res.status);
                }
            } catch (error) {
                console.error(`Error on comment ${reaction}:`, error);
            }
        }
    });
};
