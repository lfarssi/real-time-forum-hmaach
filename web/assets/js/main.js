import {
    getPosts,
    getComments,
    createPost,
    createComment
} from './app/api.js';
import {
    submitRegistration,
    submitLogin,
    submitLogout
} from './app/auth.js';

const token = localStorage.getItem('token')

async function loadPosts() {
    try {
        const posts = await getPosts(1);
        console.log("posts: ", posts);
    } catch (error) {
        console.log(error);
    }
}

async function loadComments() {
    try {
        const comments = await getComments(1, 1);
        console.log("comments: ", comments);
    } catch (error) {
        console.log(error);
    }
}

async function submitPost() {
    const post = {
        title: "New Post",
        content: "This is a new post created by the user.",
        categories: [2, 4]
    }
    try {
        const message = await createPost(post, token);
        console.log("creation post message: ", message);
    } catch (error) {
        console.log(error);
    }
}

async function submitComment() {
    const comment = {
        post_id: 1,
        content: "This is a new post created by the user."
    }
    try {
        const message = await createComment(comment, token);
        console.log("creation post message: ", message);
    } catch (error) {
        console.log(error);
    }
}

if (!localStorage.getItem('token') || !localStorage.getItem('user')) {
    const container = document.getElementById('container');
    const loginForm = document.createElement("div");
    loginForm.innerHTML = `
        <div id="login">
            <input type="text" id="login-identifier" placeholder="Nickname or Email">
            <input type="password" id="login-password" placeholder="Password">
            <button id="login-submit">Log in</button>
        </div>
    `;
    container.appendChild(loginForm);
    localStorage.removeItem("user");
    localStorage.removeItem("token");
} else {
    const user = JSON.parse(localStorage.getItem('user'));

    const container = document.getElementById('container');
    const userDisplay = document.createElement("div");
    userDisplay.innerHTML = `
        <h2>Welcome, ${user.first_name} ${user.last_name}!</h2>
        <p><strong>Nickname:</strong> ${user.nickname}</p>
        <p><strong>Email:</strong> ${user.email}</p>
        <p><strong>Age:</strong> ${user.age}</p>
        <p><strong>Gender:</strong> ${user.gender}</p>
        <button id="logout-submit">Log out</button>
    `;
    container.appendChild(userDisplay);
}

const loginForm = document.getElementById("login-submit")
if (loginForm) {
    loginForm.addEventListener("click", (e) => {
        const credentials = {
            "nickname": document.getElementById("login-identifier").value,
            "password": document.getElementById("login-password").value
        }
        submitLogin(credentials)
    })
}

const logoutBtn = document.getElementById("logout-submit")
if (logoutBtn) {
    logoutBtn.addEventListener("click", (e) => {
        submitLogout()
    })
}

loadPosts()
loadComments()
// submitPost()
// submitComment()
// submitRegistration()
// submitLogin()
// submitLogout()