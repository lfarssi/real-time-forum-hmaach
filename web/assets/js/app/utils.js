export const writeMessage = (elementID, message) => {
    const element = document.getElementById(elementID);
    if (element) {
        element.textContent = message;
    }
    setTimeout(() => {
        element.textContent = ''
    }, 5000);
}

export const updateUserStatus = (connectedUsers) => {
    const userListContainer = document.getElementById("user-list");
    const allUserElements = userListContainer.querySelectorAll('.user');

    allUserElements.forEach(userElement => {
        const userID = userElement.getAttribute('data-user-id');
        if (userID) {
            if (connectedUsers.includes(parseInt(userID))) {
                console.log(parseInt(userID))
                userElement.querySelector('span').classList.add('online')
            } else {
                userElement.querySelector('span').classList.remove('online')
            }
        }
    });
};
