import { getPosts, reactToPost } from './api.js';
import { showErrorPage, formatTime } from './utils.js';
import { showPostDetail } from './post_page.js';

export const showFeed = async () => {
    document.querySelector('main').innerHTML = ''
    const postContainer = document.createElement('div')
    postContainer.className = 'post-container'
    document.querySelector('main').append(postContainer)

    try {
        const token = localStorage.getItem('token');
        const response = await getPosts(1, token);
        if (response.status === 200) renderPosts(response.posts);
        else throw response
    } catch (error) {
        console.log(error);
        showErrorPage(error.status, error.response);
    }
};

const renderPosts = (posts) => {
    if (!posts) return
    const postContainer = document.querySelector('.post-container');
    posts.forEach(post => {
        const postDiv = document.createElement('div');
        postDiv.className = 'post';
        postDiv.innerHTML =/*html*/`
        <div class="user-info">
            <img src="https://ui-avatars.com/api/?name=${post.nickname}" alt="profile">
            <div>
                <div class="username">${post.nickname}</div>
                <div class="timestamp">${formatTime(post.created_at)}</div>
            </div>
        </div>
        <div class="post-content">
            <h3>${post.title}</h3>
            <p><pre>${post.content}</pre></p>
        </div>
        <div class="tags-reactions">
            <div class="tags">
                ${post.categories.map(category => `<span>${category.label}</span>`).join('')}
            </div>
            <div class="reactions">
                <div>
                    <i class="fa-solid fa-thumbs-up ${post.is_reacted === 1 ? 'like' : ''}"></i>
                    <span>${post.likes_count}</span>
                </div>
                <div>
                    <i class="fa-solid fa-thumbs-down ${post.is_reacted === -1 ? 'dislike' : ''}"></i>
                    <span>${post.dislikes_count}</span>
                </div>
                <div>
                    <i class="fa-solid fa-comment-dots"></i>
                    <span>${post.comments_count}</span>
                </div>
            </div>
        </div>
        `
        // add reaction events
        const likeIcon = postDiv.querySelector('.fa-thumbs-up');
        const dislikeIcon = postDiv.querySelector('.fa-thumbs-down');
        likeIcon.addEventListener('click', () => handleReaction(post.id, 'like', likeIcon));
        dislikeIcon.addEventListener('click', () => handleReaction(post.id, 'dislike', dislikeIcon));

        // add events to show post detail
        const title = postDiv.querySelector('.post-content h3');
        const commentIcon = postDiv.querySelector('.fa-comment-dots');
        title.addEventListener('click', () => showPostDetail(post));
        commentIcon.addEventListener('click', () => showPostDetail(post));

        postContainer.append(postDiv)
    });
};

export const handleReaction = async (postId, type, icon) => {
    try {
        const token = localStorage.getItem('token');
        const response = await reactToPost({ post_id: postId, reaction: type }, token);

        if (response.status === 200) {
            // Update reaction counts and styles
            const isLike = type === 'like';
            const otherIcon = icon.closest('.reactions').querySelector(isLike ? '.fa-thumbs-down' : '.fa-thumbs-up');

            // Toggle current reaction
            if (icon.classList.contains(type)) {
                icon.classList.remove(type);
                icon.nextElementSibling.textContent = parseInt(icon.nextElementSibling.textContent) - 1;
            } else {
                icon.classList.add(type);
                icon.nextElementSibling.textContent = parseInt(icon.nextElementSibling.textContent) + 1;

                // Remove other reaction if exists
                if (otherIcon.classList.contains(isLike ? 'dislike' : 'like')) {
                    otherIcon.classList.remove(isLike ? 'dislike' : 'like');
                    otherIcon.nextElementSibling.textContent = parseInt(otherIcon.nextElementSibling.textContent) - 1;
                }
            }
        } else {
            throw response;
        }
    } catch (error) {
        showErrorPage(error.status, error.message);
    }
};
