import { handleLogout } from "./auth.js";
import { showCreatePost } from "./create_post.js";
import { showFeed } from "./feed.js";

export const setupLayout = () => {
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
            </main>
        </div>
    `;
    // Initialize components
    setupSidebar();
    setupEventListeners();
}

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
    newPostBtn.addEventListener('click', showCreatePost);

    const searchInput = document.querySelector('#sidebar input[type="search"]');
    // searchInput.addEventListener('input', handleSearch);
};
