import { showAuth } from './app/auth.js';
import { showFeed } from './app/feed.js';
import { setupLayout } from './app/layout.js';
import { showErrorPage } from './app/utils.js';
import { setupWebSocket } from './app/websocket.js';

document.addEventListener('DOMContentLoaded', () => {
    if (location.pathname !== "/") {
        showErrorPage(404);
        return;
    }
    handleLoading();
});

export const handleLoading = () => {
    const user = JSON.parse(localStorage.getItem("user") || "{}");
    const token = localStorage.getItem("token");

    if (user && token) {
        setupWebSocket();
        setupLayout();
        showFeed(user);
    } else {
        localStorage.removeItem("user");
        localStorage.removeItem("token");
        showAuth();
    }
}
