import { showNotification, updateUserStatus } from "./utils.js";

let ws

export const setupWebSocket = () => {
    const token = localStorage.getItem('token');
    ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

    ws.onopen = function () {
        console.log('WebSocket is open');
    };

    ws.onmessage = function (event) {
        try {
            const data = JSON.parse(event.data);

            if (data.type === 'users-status' && Array.isArray(data.users)) {
                updateUserStatus(data.users);
            } else if (data.type === 'message') {
                const notification = `New message from ${data.sender}: ${data.content}`
                showNotification("message", notification)
            } else if (data.type === 'error') {
                showNotification("error", data.message)
            }
        } catch (error) {
            showNotification("error", error)
            console.error('Error parsing WebSocket message:', error);
        }
    };
};

export const sendMessage = (receiver, message) => {
    const data = {
        receiver_id: receiver,
        type: 'message',
        content: message
    }
    ws.send(JSON.stringify(data));
}

export const closeWebsocket = () => {
    ws.close();
}