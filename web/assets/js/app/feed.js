import { getPosts } from './api.js';
import { handleLogout } from './auth.js';
import { showErrorPage, formatTime } from './utils.js';
import { setupWebSocket } from './websocket.js';

export const showFeed = () => {
    document.body.innerHTML = `
        <div id="header-container">
            <header>
                <button id="sidebar-toggle" class="sidebar-toggle">
                    <i class="fa-solid fa-bars"></i>
                </button>
                <div class="title">Forum</div>
                <div class="header-btns">
                    <div class="new-post-btn">
                        <i class="fa-sharp fa-solid fa-plus"></i>
                        <span>new post</span>
                    </div>
                    <div class="logout-btn">
                        <i class="fa-solid fa-power-off"></i>
                        <span>log out</span>
                    </div>
                </div>
            </header>
        </div>

        <div id="body-container">
            <aside id="sidebar" class="sidebar">
                <input type="search" placeholder="search for user...">
                <div class="members-list"></div>
            </aside>

            <main>
                <div class="post-container"></div>
            </main>
        </div>
    `;

    // Initialize components
    setupSidebar();
    setupEventListeners();
    setupWebSocket();
    loadPosts();
};

const setupSidebar = () => {
    const sidebarToggle = document.getElementById('sidebar-toggle');
    const sidebar = document.getElementById('sidebar');
    const body = document.body;

    // Create and append overlay
    const overlay = document.createElement('div');
    overlay.className = 'sidebar-overlay';
    document.getElementById('body-container').appendChild(overlay);

    const toggleSidebar = () => {
        sidebar.classList.toggle('active');
        overlay.classList.toggle('active');
        body.classList.toggle('sidebar-open');
    };

    sidebarToggle.addEventListener('click', toggleSidebar);
    overlay.addEventListener('click', toggleSidebar);

    // Handle window resize
    window.addEventListener('resize', () => {
        if (window.innerWidth > 900 && sidebar.classList.contains('active')) {
            toggleSidebar();
        }
    });
};

const setupEventListeners = () => {
    const homeBtn = document.querySelector('.title')
    homeBtn.addEventListener('click', showFeed);

    const logoutBtn = document.querySelector('.logout-btn');
    logoutBtn.addEventListener('click', handleLogout);

    const newPostBtn = document.querySelector('.new-post-btn');
    newPostBtn.addEventListener('click', () => showCreatePost);

    const searchInput = document.querySelector('#sidebar input[type="search"]');
    // searchInput.addEventListener('input', handleSearch);
};

const loadPosts = async () => {
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
            <h3 onclick="openPost(${post.id})">${post.title}</h3>
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
        postContainer.append(postDiv)
    });
};
