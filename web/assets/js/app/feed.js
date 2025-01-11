import { handleLogout } from "./auth.js"
import { getPosts } from "./api.js"

export const showFeed = (user) => {
    const authContainer = document.getElementById("auth-container");
    const feed = document.getElementById("feed");

    authContainer.style.display = "none";
    feed.style.display = "block";

    if (user) {
        const userDisplay = document.createElement("div");
        userDisplay.id = "user-display";
        userDisplay.innerHTML = `
        <h2>Welcome, ${user.first_name} ${user.last_name}!</h2>
        <p><strong>Nickname:</strong> ${user.nickname}</p>
        <p><strong>Email:</strong> ${user.email}</p>
        <p><strong>Age:</strong> ${user.age}</p>
        <p><strong>Gender:</strong> ${user.gender}</p>
        <button id="logout-submit">Log out</button>
        `;
        feed.innerHTML = "";
        feed.appendChild(userDisplay);
        handleLogout();
    }
    loadPosts()
};

async function loadPosts() {
    try {
        const token = localStorage.getItem("token");
        const posts = await getPosts(1, token);
        console.log("posts: ", posts);
    } catch (error) {
        console.log(error);
    }
}