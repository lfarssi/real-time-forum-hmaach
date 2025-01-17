import { handleLogout } from "./auth.js"
import { getUsers, getPosts } from "./api.js"
import { setupWebSocket, sendMessage } from "./websocket.js"
import { updateUserStatus } from "./utils.js"

export const showFeed = (user) => {
    document.body.innerHTML = ``;
    const feedContainer = document.createElement('div');
    feedContainer.id = 'feed-container';

    feedContainer.innerHTML = `
    <div id="feed">
    <div id="user-display"></div>
    <div class="chat-container">
        <div class="message-form">
            <label for="ws-receiver" class="form-label">To:</label>
            <input type="text" id="ws-receiver" class="form-input" placeholder="Receiver ID">
            <label for="ws-message" class="form-label">Message:</label>
            <input type="text" id="ws-message" class="form-input" placeholder="Type your message">
            <button id="ws-send-message" class="send-button">Send</button>
        </div>
    </div>
    <span id="ws-result"></span>
            <div id="user-list-container">
                <h3>Users</h3>
                <div id="user-list"></div>
            </div>
            <div id="posts-container"></div>
        </div>
    `;
    document.body.appendChild(feedContainer);

    // Display user info if logged in
    if (user) {
        const userDisplayContainer = document.getElementById('user-display');
        const userDisplay = document.createElement('div');
        userDisplay.id = "user-info";
        userDisplay.innerHTML = `
            <h2>Welcome, ${user.first_name} ${user.last_name}!</h2>
            <p><strong>Nickname:</strong> ${user.nickname}</p>
            <p><strong>Email:</strong> ${user.email}</p>
            <p><strong>Age:</strong> ${user.age}</p>
            <p><strong>Gender:</strong> ${user.gender}</p>
            <button id="logout-submit">Log out</button>
        `;
        userDisplayContainer.appendChild(userDisplay);

        // Setup logout functionality
        handleLogout();
    }

    loadUsers();
    loadPosts();
    
    // Setup WebSocket and load data
    setupWebSocket();

    // Add event listener for WebSocket messages
    handleChatMessages()
};

async function loadUsers() {
    try {
        const token = localStorage.getItem("token");
        const response = await getUsers(token);
        const userListContainer = document.getElementById("user-list");
        userListContainer.innerHTML = "";

        if (!response.users || response.users.length === 0) {
            userListContainer.innerText = "No response.users to display";
            return;
        }

        response.users.forEach(user => {
            const userElement = document.createElement("div");

            userElement.classList.add("user");
            userElement.setAttribute("data-user-id", user.id);

            userElement.innerHTML = `
                <p>
                    <strong>${user.first_name} ${user.last_name}</strong>
                    (@${user.nickname})
                    <span class="user-status"></span>
                </p>
            `;

            userListContainer.appendChild(userElement);
        });

        if (response.connected && Array.isArray(response.connected) && response.connected.length > 0) {
            updateUserStatus(response.connected);
        }
    } catch (error) {
        console.error("Error loading users:", error);
    }
}


async function loadPosts() {
    try {
        const token = localStorage.getItem("token");
        const posts = await getPosts(1, token);
        const postContainer = document.getElementById("posts-container");

        postContainer.innerHTML = ""; // Clear any existing content

        if (!posts || posts.length === 0) {
            postContainer.innerText = "No posts available";
            return;
        }

        posts.forEach(post => {
            const postElement = document.createElement("div");

            postElement.classList.add("post")

            postElement.innerHTML = `
                <h3>${post.title}</h3>
                <p><strong>Author:</strong> ${post.user_first_name} ${post.user_last_name} (@${post.user_nickname})</p>
                <p><strong>Category:</strong> ${post.categories.map(cat => cat.label).join(", ")}</p>
                <p><strong>Comments:</strong> ${post.comments_count}</p>
                <p><strong>Created at:</strong> ${new Date(post.created_at).toLocaleString()}</p>
                <p>${post.content}</p>
                <p><strong>Comments:</strong> ${post.comments_count}</p>
                <p><strong>Likes:</strong> ${post.likes_count}</p>
                <p><strong>Dislikes:</strong> ${post.dislikes_count}</p>
                <p><strong>Is reacted by you:</strong> ${post.is_reacted}</p>
            `;

            postContainer.appendChild(postElement);
        });
    } catch (error) {
        console.error("Error loading posts:", error);
    }
}

const handleChatMessages = () => {
    const send = document.getElementById('ws-send-message');
    if (send) {
        send.addEventListener('click', () => {
            const receiver = document.getElementById('ws-receiver');
            const message = document.getElementById('ws-message');
            if (receiver && message) {
                sendMessage(parseInt(receiver.value), message.value);
                message.value = '';
                receiver.value = '';
            }
        });
    }
}
