import EventEmitter from "events";

const ws_scheme = window.location.protocol === "https:" ? "wss" : "ws";
const socket = new WebSocket(ws_scheme + "://" + window.location.host + "/ws");

const bus = new EventEmitter();

bus.on('log subscribe', () => {
    socket.send(
        JSON.stringify(
            {
                room_name: "",
                controls: {
                    type: "subscribe",
                    value: "gamelog"
                }
            }
        )
    );
});

bus.on('log unsubscribe', () => {
    socket.send(
        JSON.stringify(
            {
                room_name: "",
                controls: {
                    type: "unsubscribe",
                    value: "gamelog"
                }
            }
        )
    );
})

bus.on('server status subscribe', () => {
    socket.send(
        JSON.stringify(
            {
                room_name: "",
                controls: {
                    type: "subscribe",
                    value: "server_status"
                }
            }
        )
    );
});

bus.on('command send', command => {
    socket.send(
        JSON.stringify(
            {
                room_name: "",
                controls: {
                    type: "command",
                    value: command
                }
            }
        )
    );
});

socket.onmessage = e => {
    const {room_name, message} = JSON.parse(e.data);
    bus.emit(room_name, message);
}

export default bus;