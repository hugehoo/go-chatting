import ws from 'k6/ws';
import { check } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import { Counter } from 'k6/metrics';

export const options = {
    stages: [
        { duration: '30', target: 100 },  // 2초 동안 20 VU로 증가
        { duration: '30', target: 200 },  // 2초 동안 20 VU로 증가
        { duration: '30', target: 300 },  // 2초 동안 20 VU로 증가
        { duration: '30', target: 400 }  // 2초 동안 20 VU로 증가
    ],
};


export default function () {
    const url = 'ws://localhost:1010/room-chat';
    const params = {
        headers: {
            'Cookie': 'auth=test'
        }
    };


    const res = ws.connect(url, params, function (socket) {
        socket.on('open', () => {
            console.log('WebSocket 연결됨');
            // 메시지 데이터 준비
            const messageData = JSON.stringify({
                "name": `user_${randomString(5)}`,
                "message": `테스트 메시지 ${randomString(10)}`,
                "room": "Load Test Room",
                "when": new Date().toISOString(),
            });

            // 메시지 전송
            socket.send(messageData);
        });
        socket.on('message', (data) => {
            console.log(`메시지 수신:`, data);
        });

        socket.on('error', (e) => {
            if (e.error() != 'websocket: close sent') {
                console.log('예상치 못한 오류 발생:', e.error());
            }
        });

        socket.on('close', () => {
            console.log('WebSocket 연결 닫힘');
        });

        // 1초 후에 연결 종료
        setTimeout(() => {
            socket.close();
        }, 1000);
    });

    check(res, { 'status is 101': (r) => r && r.status === 101 });
}