import EventEmitter from "events";

const ws_scheme = window.location.protocol === "https:" ? "wss" : "ws";
const socket = new WebSocket(ws_scheme + "://" + window.location.host + "/ws");

const bus = new EventEmitter();

bus.on('log subscribe', () => {
    socket.send(JSON.stringify({name: 'log subscribe'}));
});

bus.on('server status subscribe', () => {
    socket.send(JSON.stringify({name: 'server status subscribe'}));
});

bus.on('command send', command => {
    socket.send(JSON.stringify({name: 'command send', data: command}));
});

socket.onmessage = e => {
    const {name, data} = JSON.parse(e.data)
    bus.emit(name, data);
}

export default bus;