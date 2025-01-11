export const writeMessage = (elementID, message) => {
    const element = document.getElementById(elementID);
    if (element) {
        element.textContent = message;
    }
    setTimeout(() => {
        element.textContent = ''
    }, 2000);
}

