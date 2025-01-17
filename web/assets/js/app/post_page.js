// post_detail.js
import { getComments, createComment } from './api.js';
import { showErrorPage, formatTime } from './utils.js';

export const showPostDetail = async (post) => {
    const mainContainer = document.querySelector('main');
    mainContainer.innerHTML = `
        <div class="post-detail-container">
            <div class="post main-post">
                <div class="user-info">
                    <img src="https://ui-avatars.com/api/?name=${post.nickname}" alt="profile">
                    <div>
                        <div class="username">${post.nickname}</div>
                        <div class="timestamp">${formatTime(post.created_at)}</div>
                    </div>
                </div>
                <div class="post-content-detailed">
                    <h3>${post.title}</h3>
                    <p>${post.content}</p>
                </div>
                <div class="tags-reactions">
                    <div class="tags">
                        ${post.categories.map(category => `<span>${category.label}</span>`).join('')}
                    </div>
                    <div class="reactions">
                        <div>
                            <i class="fa-solid fa-thumbs-up ${post.is_reacted === 1 ? 'like' : ''}" data-action="like"></i>
                            <span>${post.likes_count}</span>
                        </div>
                        <div>
                            <i class="fa-solid fa-thumbs-down ${post.is_reacted === -1 ? 'dislike' : ''}" data-action="dislike"></i>
                            <span>${post.dislikes_count}</span>
                        </div>
                        <div>
                            <i class="fa-solid fa-comment-dots"></i>
                            <span>${post.comments_count}</span>
                        </div>
                    </div>
                </div>
            </div>

            <div class="comment-form">
                <textarea id="comment-input" placeholder="leave a comment..." required></textarea>
                <div class="form-actions">
                    <button class="comment-btn">Comment</button>
                </div>
            </div>

            <div class="comments-list"></div>
        </div>
    `;

    // Load comments
    try {
        const token = localStorage.getItem('token');
        const comments = await getComments(post.id, 1, token);
        console.log("hello");
        
        renderComments(comments);
        setupCommentForm(post.id);
    } catch (error) {
        showErrorPage(error);
    }
};

const renderComments = (comments) => {
    const commentsContainer = document.querySelector('.comments-list');
    if (!comments) return
    comments.forEach(comment => {
        const commentElement = document.createElement('div')
        commentElement.className = 'comment'
        commentElement.innerHTML = `
        <div class="user-info">
            <img src="https://ui-avatars.com/api/?name=${comment.nickname}" alt="profile">
            <div>
                <div class="username">${comment.nickname}</div>
                <div class="timestamp">${formatTime(comment.created_at)}</div>
            </div>
        </div>
        <p>${comment.content}</p>
        `
        commentsContainer.append(commentElement)
    });
};

const setupCommentForm = (postId) => {
    const commentInput = document.getElementById('comment-input');
    const commentBtn = document.querySelector('.comment-btn');

    commentBtn.addEventListener('click', async () => {
        const content = commentInput.value.trim();
        if (!content) return;

        try {
            const token = localStorage.getItem('token');
            const commentData = {post_id: postId, content: content};

            const response = await createComment(commentData, token);
            if (response.status === 200) {
                const comments = await getComments(postId, 1, token);
                document.querySelector('.comments-list').innerHTML = '';
                renderComments(comments);
                commentInput.value = '';
            } else {
                throw response;
            }
        } catch (error) {
            showErrorPage(error);
        }
    });
};