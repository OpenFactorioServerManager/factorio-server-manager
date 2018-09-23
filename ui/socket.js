import EventEmitter from 'events';

class Socket {
    constructor(ws, ee = new EventEmitter()){
        this.ws = ws;
        this.ee = ee;
        ws.onmessage = this.message.bind(this);
        ws.onopen = this.open.bind(this);
        ws.onclose = this.close.bind(this);

        this.opened = false;
    }
    on(name, fn){
        this.ee.on(name, fn);
    }
    off(name, fn){
        this.ee.removeListener(name, fn);
    }
    emit(name, data){
        if(this.ws.readyState == WebSocket.OPEN) {
            const message = JSON.stringify({name, data});
            this.ws.send(message);
        }

        return this.ws.readyState;
    }
    message(e){
        try{
            let message = JSON.parse(e.data);
            // console.log(message.name, message.data);
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
