import { getComments, createComment } from './api.js';
import { showErrorPage, formatTime, showNotification, debounce } from './utils.js';
import { handleReaction } from './feed.js'

let currentCommentPage = 1;
let isLoadingComments = false;
let hasMoreComments = true;
let postId

export const showPostDetail = async (post) => {
    currentCommentPage = 1;
    isLoadingComments = false;
    hasMoreComments = true;
    postId = post.id;    

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
                    <pre>${post.content}</pre>
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

            <div class="comments-list">
                <div class="loading-indicator">
                    <i class="fa-solid fa-spinner fa-spin"></i> Loading more comments...
                </div>
            </div>
        </div>
    `;
    const likeIcon = document.querySelector('.fa-thumbs-up');
    const dislikeIcon = document.querySelector('.fa-thumbs-down');
    likeIcon.addEventListener('click', () => handleReaction(post.id, 'like', likeIcon));
    dislikeIcon.addEventListener('click', () => handleReaction(post.id, 'dislike', dislikeIcon));

    setupCommentForm();
    await loadComments();
    setupCommentScroll();
};

const loadComments = async () => {
    try {
        const token = localStorage.getItem('token');
        const response = await getComments(postId, currentCommentPage, token);
        if (response.status !== 200) throw response;
        renderComments(response.comments);
    } catch (error) {
        showErrorPage(error.status, error.message);
    }
};

const renderComments = (comments) => {
    const commentsContainer = document.querySelector('.comments-list');
    const loadingIndicator = commentsContainer.querySelector('.loading-indicator');

    if (!comments || comments.length === 0) {
        loadingIndicator.style.display = 'none';
        hasMoreComments = false;
        return;
    }

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
        <pre>${comment.content}</pre>
        `
        commentsContainer.insertBefore(commentElement, loadingIndicator);
    });

    if (comments.length < 10) {
        loadingIndicator.style.display = 'none';
        hasMoreComments = false;
        return
    }
    loadingIndicator.style.display = 'block';

};

const setupCommentScroll = () => {
    const main = document.querySelector('main')
    const scrollCommentFunc = debounce(async () => {
        if (isLoadingComments || !hasMoreComments) return;
        console.log(main.scrollHeight, main.clientHeight, main.scrollTop);
        
        // Load more when scrolled near bottom (100px threshold)
        if (main.scrollHeight - (main.clientHeight + main.scrollTop) <= 100) {
            isLoadingComments = true;

            try {
                currentCommentPage++;
                const token = localStorage.getItem('token');
                const response = await getComments(postId, currentCommentPage, token);

                if (response.status !== 200) throw response;
                renderComments(response.comments);
            } catch (error) {
                showErrorPage(error.status, error.message);
            } finally {
                isLoadingComments = false;
            }
        }
    }, 800);

    main.removeEventListener('scroll', scrollCommentFunc);
    main.addEventListener('scroll', scrollCommentFunc);
};

const setupCommentForm = () => {
    const commentInput = document.getElementById('comment-input');
    const commentBtn = document.querySelector('.comment-btn');

    commentBtn.addEventListener('click', async () => {
        const content = commentInput.value.trim();
        if (!content) return showNotification('error', 'please write a comment');

        try {
            const token = localStorage.getItem('token');
            const commentData = { post_id: postId, content: content };

            const response = await createComment(commentData, token);
            if (response.status !== 200) throw response;

            document.querySelector('.comments-list').innerHTML = '';

            // Reset comments list and pagination
            const commentsContainer = document.querySelector('.comments-list');
            commentsContainer.innerHTML = '<div class="loading-indicator" style="display: none"><i class="fa-solid fa-spinner fa-spin"></i> Loading more comments...</div>';

            currentCommentPage = 1;
            hasMoreComments = true;
            await loadComments(postId);
            commentInput.value = '';

            const totalComments = document.querySelector('.main-post').querySelector('.fa-comment-dots').nextElementSibling;
            totalComments.textContent = parseInt(totalComments.textContent) + 1;
        } catch (error) {
            showErrorPage(error.status, error.message);
        }
    });
};