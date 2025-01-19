import { showAuth } from "./auth.js";

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

export const showErrorPage = (error) => {
    showNotification('error', error)
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

export const showNotification = (type, message) => {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;

    notification.innerHTML = `
        <span class="notification-message">${message}</span>
        <button class="notification-close">&times;</button>
    `;

    document.body.appendChild(notification);

    const closeButton = notification.querySelector('.notification-close');
    closeButton.addEventListener('click', () => {
        notification.remove();
    });
}

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
            showAuth();
        }, 3000);
        return true;
    }
    return false;
};