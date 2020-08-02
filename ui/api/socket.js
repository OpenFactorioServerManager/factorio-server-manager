import EventEmitter from 'events';

class Socket {
    constructor(){
        let ws_scheme = window.location.protocol === "https:" ? "wss" : "ws";
        this.ws = new WebSocket(ws_scheme + "://" + window.location.host + "/ws");

        this.ee = new EventEmitter();
        this.ws.onmessage = this.message.bind(this);
        this.ws.onopen = this.open.bind(this);
        this.ws.onclose = this.close.bind(this);

        this.opened = false;
    }
    on(name, fn){
        this.ee.on(name, fn);
    }
    off(name, fn){
        this.ee.removeListener(name, fn);
    }
    emit(name, data){
        if(this.ws.readyState === WebSocket.OPEN) {
            const message = JSON.stringify({name, data});
            this.ws.send(message);
        }

        return this.ws.readyState;
    }
    message(e){
        try{
            let message = JSON.parse(e.data);
            this.ee.emit(message.name, message.data);
        }
        catch(err){
            this.ee.emit('error', err);
        }
    }
    open(){
        this.ee.emit('connect');
    }
    close(){
        this.ee.emit('disconnect');
    }
}

export default Socket;
