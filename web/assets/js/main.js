import { showAuth } from './app/auth.js';
import { showFeed } from './app/feed.js';

addEventListener('DOMContentLoaded', () => {
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
