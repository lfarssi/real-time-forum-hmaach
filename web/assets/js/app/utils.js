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
    document.body.innerHTML = `
        <div style="display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100vh; gap: 20px;">
            <h1 style="color: var(--main-red);">Error</h1>
            <p style="color: var(--text-primary);">${error.message}</p>
            <button onclick="showFeed()" 
                    style="padding: 10px 20px; background: var(--main-green); border: none; border-radius: 5px; color: white; cursor: pointer;">
                Retry
            </button>
        </div>
    `;
}

export const formatTime = (time) => {
    const date = new Date(time);
    const diff = Date.now() - date;

    // Convert milliseconds to minutes/hours/days
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);

    if (minutes < 60) {
        return `${minutes} minute${minutes !== 1 ? 's' : ''} ago`;
    } else if (hours < 24) {
        return `${hours} hour${hours !== 1 ? 's' : ''} ago`;
    } else if (days < 7) {
        return `${days} day${days !== 1 ? 's' : ''} ago`;
    } else {
        return date.toLocaleDateString();
    }
}
