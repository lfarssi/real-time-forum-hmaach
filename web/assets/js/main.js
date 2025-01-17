import { showAuth } from './app/auth.js';
import { showFeed } from './app/feed.js';
import { getComments, createPost, createComment } from './app/api.js';

addEventListener('DOMContentLoaded', () => {
    // Check if user and token exist
    const user = JSON.parse(localStorage.getItem("user") || "{}");
    const token = localStorage.getItem("token");

    if (user && token) {
        showFeed(user);
    } else {
        localStorage.removeItem("user");
        localStorage.removeItem("token");
        showAuth();
    }
});

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