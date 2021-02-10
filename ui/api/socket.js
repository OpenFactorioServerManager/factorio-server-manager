import EventEmitter from "events";

const bus = new EventEmitter();

const ws_scheme = window.location.protocol === "https:" ? "wss" : "ws";

function connect() {
    const socket = new WebSocket(ws_scheme + "://" + window.location.host + "/ws");

    function logSubscribeEvent() {
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
    }

    function logUnsubscribeEvent() {
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
    }

    function serverStatusSubscribeEvent() {
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
    }

    function commandSendEvent(command) {
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
    }

    function registerEventEmitter() {
        bus.on('log subscribe', logSubscribeEvent);
        bus.on('log unsubscribe', logUnsubscribeEvent);
        bus.on('server status subscribe', serverStatusSubscribeEvent);
        bus.on('command send', commandSendEvent);
    }

    function unregisterEventEmitter() {
        bus.off('log subscribe', logSubscribeEvent);
        bus.off('log unsubscribe', logUnsubscribeEvent);
        bus.off('server status subscribe', serverStatusSubscribeEvent);
        bus.off('command send', commandSendEvent);
    }

    socket.onmessage = e => {
        const {room_name, message} = JSON.parse(e.data);
        bus.emit(room_name, message);
    }

    socket.onerror = e => {
        socket.close();
    }

    socket.onclose = e => {
        unregisterEventEmitter()
        // reconnect after 5 seconds
        setTimeout(connect, 5000);
    }

    socket.onopen = e => {
        registerEventEmitter(socket)
    }
}

connect();

export default bus;