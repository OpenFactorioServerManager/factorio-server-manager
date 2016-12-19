import React from 'react';
import {IndexLink} from 'react-router';
import Socket from '../../socket.js';

class ConsoleContent extends React.Component {
    constructor(props) {
        super(props);
        this.componentDidMount = this.componentDidMount.bind(this);
        this.connectWebsocket = this.connectWebsocket.bind(this);
    }

    componentDidMount() {
        this.connectWebsocket()
    }

    connectWebsocket() {
        var ws_scheme = window.location.protocol == "https:" ? "wss" : "ws";
        let ws = new WebSocket(ws_scheme + "://" + window.location.host + "/ws");
        let socket = this.socket = new Socket(ws);

        
    }
}
