import { showNotification, updateUserStatus, trimString } from "./utils.js";
import { appendMessage, chatID, showTypingIndicator } from './chat.js';
import { loadUsers } from './layout.js';

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
            switch (data.type) {
                case 'message': {
                    // Check if we're in the chat with this sender
                    const currentChat = document.querySelector('.chat-main');

                    if (currentChat && data.sender_id === chatID) {
                        appendMessage(data);
                    } else {
                        const notification = `New message from ${trimString(data.sender, 10)}: ${trimString(data.content, 7)}`;
                        showNotification("message", notification);
                    }
                    break;
                }

                case 'typing-start': {
                    handleTypingStatus(data)
                    break;
                }

                case 'typing-stop': {
                    handleTypingStatus(data)
                    break;
                }

                case 'users-status':
                    updateUserStatus(data.users);
                    break;

                case 'refresh-users':
                    loadUsers();
                    break;

                case 'error':
                    showNotification("error", data.message);
                    break;

                default:
                    showNotification("error", data.type);
                    console.log(data.type);
                    break;
            }
        } catch (error) {
            showNotification("error", error);
            console.error('Error parsing WebSocket message:', error);
        }
    };
};

const handleTypingStatus = (data) => {
    const currentChat = document.querySelector('.chat-main');
    if (currentChat && data.sender_id === chatID) {
        const isTyping = data.type === 'typing-start';
        showTypingIndicator(isTyping, data.sender);
    }
};

export const sendMessage = (receiver, type, message) => {
    const data = {
        receiver_id: receiver,
        type: type,
        content: message
    }
    ws.send(JSON.stringify(data));
}

export const getOnlineUsers = () => {
    ws.send("status")
}

export const closeWebsocket = () => {
    ws.close();
}
