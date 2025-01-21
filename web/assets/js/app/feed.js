import { getPosts, reactToPost } from './api.js';
import { showErrorPage, formatTime, debounce } from './utils.js';
import { showPostDetail } from './post_page.js';

let currentPage = 1;
let isLoading = false;
let hasMorePosts = true;

export const showFeed = async () => {
    currentPage = 1;
    isLoading = false;
    hasMorePosts = true;

    document.querySelector('main').innerHTML = ''
    const postContainer = document.createElement('div')
    postContainer.className = 'post-container'
    document.querySelector('main').append(postContainer)

    const loadingIndicator = document.createElement('div');
    loadingIndicator.className = 'loading-indicator';
    loadingIndicator.style.display = 'none';
    loadingIndicator.innerHTML = /*html*/`<i class="fa-solid fa-spinner fa-spin"></i> Loading more posts...`;
    postContainer.append(loadingIndicator);

    try {
        const token = localStorage.getItem('token');
        const response = await getPosts(currentPage, token);
        if (response.status !== 200) throw response;
        renderPosts(response.posts);
        setupInfiniteScroll();
    } catch (error) {
        console.log(error);
        showErrorPage(error.status, error.response);
    }
};

const renderPosts = (posts) => {
    const loadingIndicator = document.querySelector('.loading-indicator')
    const postContainer = document.querySelector('.post-container');
    if (!postContainer) return;
    if (!posts || posts.length === 0) {
        loadingIndicator.style.display = 'none';
        hasMorePosts = false;
        return;
    }

    posts.forEach(post => {
        const postDiv = document.createElement('div');
        postDiv.className = 'post';
        postDiv.innerHTML = /*html*/`
        <div class="user-info">
            <img src="https://ui-avatars.com/api/?name=${post.nickname}" alt="profile">
            <div>
                <div class="username">${post.nickname}</div>
                <div class="timestamp">${formatTime(post.created_at)}</div>
            </div>
        </div>
        <div class="post-content">
            <h3>${post.title}</h3>
            <pre>${post.content}</pre>
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
        title.addEventListener('click', () => showPostDetail(post.id));
        commentIcon.addEventListener('click', () => showPostDetail(post.id));

        postContainer.insertBefore(postDiv, loadingIndicator)
    });

    if (posts.length < 10) {
        loadingIndicator.style.display = 'none';
        hasMorePosts = false;
        return
    }
    loadingIndicator.style.display = 'block';
};

export const handleReaction = async (postId, type, icon) => {
    try {
        const token = localStorage.getItem('token');
        const response = await reactToPost({ post_id: postId, reaction: type }, token);

        if (response.status !== 200) throw response

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
    } catch (error) {
        console.error(error);
        showErrorPage(error.status, error.message);
    }
};

const setupInfiniteScroll = () => {
    const main = document.querySelector('main');

    // debounce scroll handler to prevent excessive calls
    const scrollFunc = debounce(async () => {
        if (isLoading || !hasMorePosts) return;

        // Load more when user scrolls to bottom with 100px threshold
        if (main.scrollHeight - (main.clientHeight + main.scrollTop) <= 100) {
            isLoading = true;

            try {
                currentPage++;
                const token = localStorage.getItem('token');
                const response = await getPosts(currentPage, token);

                if (response.status !== 200) throw response;
                renderPosts(response.posts);
            } catch (error) {
                showErrorPage(error.status, error.message);
            } finally {
                isLoading = false;
            }
        }
    }, 800);

    main.removeEventListener('scroll', scrollFunc);
    main.addEventListener('scroll', scrollFunc);
}
