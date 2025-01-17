import { getPosts } from './api.js';
import { showErrorPage, formatTime } from './utils.js';
import { showPostDetail } from './post_page.js';

export const showFeed = async () => {
    document.querySelector('main').innerHTML = ''
    const postContainer = document.createElement('div')
    postContainer.className = 'post-container'
    document.querySelector('main').append(postContainer)

    try {
        const token = localStorage.getItem('token');
        const posts = await getPosts(1, token);
        renderPosts(posts);
    } catch (error) {
        showErrorPage(error);
    }
};

const renderPosts = (posts) => {
    const postContainer = document.querySelector('.post-container');
    posts.forEach(post => {
        const postDiv = document.createElement('div');
        postDiv.className = 'post';
        postDiv.innerHTML = `
        <div class="user-info">
            <img src="https://ui-avatars.com/api/?name=${post.nickname}" alt="profile">
            <div>
                <div class="username">${post.nickname}</div>
                <div class="timestamp">${formatTime(post.created_at)}</div>
            </div>
        </div>
        <div class="post-content">
            <h3>${post.title}</h3>
            <p>${post.content}</p>
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
        const title = postDiv.querySelector('.post-content h3');
        const commentIcon = postDiv.querySelector('.fa-comment-dots');

        title.addEventListener('click', () => showPostDetail(post));
        commentIcon.addEventListener('click', () => showPostDetail(post));

        postContainer.append(postDiv)
    });
};
