import { getUsers } from "./api.js";
import { handleLogout } from "./auth.js";
import { showDirectMessages } from "./chat.js";
import { showCreatePost } from "./create_post.js";
import { showFeed } from "./feed.js";
import { updateUserStatus, formatTime } from "./utils.js";

export const setupLayout = () => {
    const user = JSON.parse(localStorage.getItem("user"));
    document.body.innerHTML = /*html*/`
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
                        <span>log out ${user.nickname}</span>
                    </div>
                </div>
            </header>
        </div>

        <div id="body-container">
            <aside id="sidebar" class="sidebar">
                <div class="members-list" id="user-list"></div>
            </aside>

            <main>
            </main>
        </div>
    `;
    // Initialize components
    setupSidebar();
    loadUsers();
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
};

export const loadUsers = async () => {
    try {
        const token = localStorage.getItem("token");
        const userID = JSON.parse(localStorage.getItem('user')).id
        const response = await getUsers(token);
        const userListContainer = document.querySelector(".members-list");
        userListContainer.innerHTML = "";

        if (!response.users || response.users.length === 0) {
            userListContainer.innerText = "No response.users to display";
            return;
        }

        response.users.forEach(user => {
            const userElement = document.createElement("div");
            userElement.classList.add("user");
            userElement.setAttribute("data-user-id", user.id);
            var last_message = (user.last_message.sender_id == userID ? "You: " : "") + user.last_message.message
            if (last_message.length > 25) last_message = last_message.slice(0, 25) + "..."
            userElement.innerHTML = `
                <div>
                    <img src="https://ui-avatars.com/api/?name=${user.first_name + user.last_name}" alt="profile">
                    <div class="user-status"></div>
                </div>
                <div class="user-list-info">
                    <span>${user.nickname}</span>
                    <div class="last-message">
                        <span>${last_message}</span>
                        <span>${formatTime(user.last_message.sent_at)}</span>
                    </div>
                </div>
            `;

            userListContainer.appendChild(userElement);

            userElement.addEventListener('click', () => showDirectMessages(user.id))
        });

        if (response.connected && Array.isArray(response.connected) && response.connected.length > 0) {
            updateUserStatus(response.connected);
        }
    } catch (error) {
        console.error("Error loading users:", error);
    }
}
