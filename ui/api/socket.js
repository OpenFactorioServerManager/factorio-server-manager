import EventEmitter from "events";

const bus = new EventEmitter();

const ws_scheme = window.location.protocol === "https:" ? "wss" : "ws";

function connect() {
    const socket = new WebSocket(ws_scheme + "://" + window.location.host + "/ws");

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

    socket.onerror = e => {
        socket.close();
    }

    socket.onclose = e => {
        // reconnect after 5 seconds
        setTimeout(connect, 5000);
    }
}

connect();

export default bus;