import { handleLogout } from "./auth.js"
import { getPosts } from "./api.js"

export const showFeed = (user) => {
    const authContainer = document.getElementById("auth-container");
    const feed = document.getElementById("feed");

    authContainer.style.display = "none";
    feed.style.display = "block";

    if (user) {
        const userDisplay = document.getElementById("user-display");
        userDisplay.innerHTML = "";
        userDisplay.innerHTML = `
        <h2>Welcome, ${user.first_name} ${user.last_name}!</h2>
        <p><strong>Nickname:</strong> ${user.nickname}</p>
        <p><strong>Email:</strong> ${user.email}</p>
        <p><strong>Age:</strong> ${user.age}</p>
        <p><strong>Gender:</strong> ${user.gender}</p>
        <button id="logout-submit">Log out</button>
        `;
        handleLogout();
    }
    loadPosts()
};

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
                <p>Comments :${post.comments_count}</p>
            `;

            postContainer.appendChild(postElement);
        });
    } catch (error) {
        console.error("Error loading posts:", error);
    }
}
