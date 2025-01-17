
export const setupWebSocket = () => {
    const token = localStorage.getItem('token');
    let ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

    ws.onopen = function () {
        console.log('WebSocket is open');
    };

    ws.onmessage = function (event) {
        const result = document.getElementById('ws-result');
        result.innerText = 'Received: ' + event.data;

        try {
            const data = JSON.parse(event.data);
            // console.log(data);
            if (Array.isArray(data)) {
                updateUserStatus(data);
            }
        } catch (e) {
            console.error('Error parsing WebSocket message:', e);
        }
    };

    return ws;
};
