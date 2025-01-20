import { showAuth } from "./auth.js";
import { closeWebsocket } from "./websocket.js";
import { handleLoading } from "../main.js";

// get FormData and convert it to an object
export const getFormData = (form) => {
    const formData = new FormData(form);
    const formObject = {};
    formData.forEach((value, key) => {
        if (key === 'age') {
            formObject[key] = Number(value);
        } else {
            formObject[key] = value;
        }
    });
    return formObject
}

export const formatTime = (time) => {
    if (!time) return ""
    const date = new Date(time);
    const diff = Date.now() - date;

    // Convert milliseconds to minutes/hours/days
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);

    if (minutes < 1) {
        return `seconds ago`;
    } else if (minutes < 60) {
        return `${minutes} minute${minutes !== 1 ? 's' : ''} ago`;
    } else if (hours < 24) {
        return `${hours} hour${hours !== 1 ? 's' : ''} ago`;
    } else if (days < 7) {
        return `${days} day${days !== 1 ? 's' : ''} ago`;
    } else {
        return date.toLocaleDateString();
    }
}

export const trimString = (string, width) => {
    if (!string || string.length <= width) return string
    return string.slice(0, width) + "..."
}

export const updateUserStatus = (connectedUsers) => {
    if (!Array.isArray(connectedUsers)) return
    const userListContainer = document.getElementById("user-list");
    const allUserElements = userListContainer.querySelectorAll('.user');

    allUserElements.forEach(userElement => {
        const userID = userElement.getAttribute('data-user-id');
        if (userID) {
            if (connectedUsers.includes(parseInt(userID))) {
                userElement.querySelector('div .user-status').classList.add('online')
            } else {
                userElement.querySelector('div .user-status').classList.remove('online')
            }
        }
    });
};

export const showErrorPage = (status, message) => {
    if (status === 404 || status === 500) {
        const mainContainer = document.querySelector('main') || document.body;
        mainContainer.innerHTML = /*html*/ `
            <div class="error-page">
                <div class="error-content">
                    <i class="fa-solid ${status === 404 ? 'fa-circle-question' : 'fa-triangle-exclamation'} error-icon"></i>
                    <h1>${status}</h1>
                    <h2>${status === 404 ? 'Page Not Found' : 'Internal Server Error'}</h2>
                    <p>${status === 404
                ? 'The page you are looking for is unavailable.'
                : 'Something went wrong. Please try again later.'}</p>
                    ${status === 404
                ? `<button class="error-btn" id="go-home">
                                <i class="fa-solid fa-home"></i> Go to Home
                           </button>`
                : `<button class="error-btn" id="reload">
                                <i class="fa-solid fa-rotate-right"></i> Refresh Page
                           </button>`}
                </div>
            </div>
        `;

        const goHome = document.getElementById('go-home');
        if (goHome) {
            goHome.addEventListener('click', () => {
                history.pushState({}, '', '/');
                handleLoading();
            });
        }

        const reload = document.getElementById('reload');
        if (reload) {
            reload.addEventListener('click', () => {
                handleLoading();
            });
        }
    } else {
        showNotification('error', message || 'An error occurred. Please try again.');
    }
};

export const showNotification = (type, message) => {
    const existingNotifications = document.querySelectorAll(`.notification.${type}`);
    if (existingNotifications) existingNotifications.forEach(notification => notification.remove());

    const notification = document.createElement('div');
    notification.className = `notification ${type}`;

    const icon = type === 'error' ? 'fa-circle-exclamation' :
        type === 'success' ? 'fa-circle-check' :
            'fa-circle-info';

    notification.innerHTML = `
        <div class="notification-content">
            <i class="fa-solid ${icon}"></i>
            <span class="notification-message">${message}</span>
        </div>
        <button class="notification-close">
            <i class="fa-solid fa-xmark"></i>
        </button>
    `;

    document.body.appendChild(notification);

    // Add show class after a small delay (for animation)
    setTimeout(() => notification.classList.add('show'), 10);

    const closeButton = notification.querySelector('.notification-close');
    const closeNotification = () => {
        notification.classList.remove('show');
        setTimeout(() => notification.remove(), 300);
    };

    closeButton.addEventListener('click', closeNotification);

    setTimeout(closeNotification, 5000);
};

export const handleUnauthorized = (response) => {
    if (response.status === 401) {
        localStorage.removeItem('user');
        localStorage.removeItem('token');

        // Create and show popup
        const popup = document.createElement('div');
        popup.className = 'popup-container';

        const message = response.message === "Session expired"
            ? "Your session has expired"
            : "You've been logged out";

        popup.innerHTML = `
            <div class="popup">
                <i class="fa-solid fa-circle-exclamation popup-icon"></i>
                <span class="popup-message">${message}</span>
                <span class="popup-subtext">Redirecting to login page...</span>
            </div>
        `;

        document.body.appendChild(popup);
        document.body.style.overflow = 'hidden';

        setTimeout(() => {
            document.body.style.overflow = '';
            closeWebsocket();
            showAuth();
        }, 3000);
        return true;
    }
    return false;
};

export const debounce = (func, delay) => {
    let timer
    return function (...args) {
        clearTimeout(timer)
        timer = setTimeout(() => func(...args), delay)
    }
}