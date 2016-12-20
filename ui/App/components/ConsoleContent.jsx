import React from 'react';
import {IndexLink} from 'react-router';
import Console from 'react-console-component'
import Socket from '../../socket.js';

class ConsoleContent extends React.Component {
    constructor(props) {
        super(props);
        this.componentDidMount = this.componentDidMount.bind(this);
        this.connectWebsocket = this.connectWebsocket.bind(this);
        this.handleCommand = this.handleCommand.bind(this);
        this.onConnect = this.onConnect.bind(this);
        this.onNewLogLine = this.onNewLogLine.bind(this);
        this.state = {}
    }

    componentDidMount() {
        this.connectWebsocket();
    }

    connectWebsocket() {
        var ws_scheme = window.location.protocol == "https:" ? "wss" : "ws";
        let ws = new WebSocket(ws_scheme + "://" + window.location.host + "/ws");
        let socket = this.socket = new Socket(ws);
        socket.on('connect', this.onConnect.bind(this));
        socket.on('log update', this.onNewLogLine.bind(this));
    }

    handleCommand(command) {
        this.refs.console.log(command);
        this.refs.console.return();
    }

    onConnect() {
        this.setState({connected: true});
        this.socket.emit("log subscribe");
        this.refs.console.log("connected to Factorio Server")
    }

    onNewLogLine(logline) {
        console.log(logline);
        console.log(this.refs.console);
        this.refs.console.log({message: logline});
        this.refs.console.return();
    }

    render() {
        return(
            <div className="content-wrapper">
                <section className="content-header">
                    <h1>
                        Server Console
                        <small>Send commands and messages to the Factorio server</small>
                    </h1>
                    <ol className="breadcrumb">
                        <li><IndexLink to="/"><i className="fa fa-dashboard"></i>Server Control</IndexLink></li>
                        <li className="active">Here</li>
                    </ol>
                </section>

                <section className="content">

                    <Console ref="console"
                        autofocus={true}
                        handler={this.handleCommand}
                    />

                </section>
            </div>
        );
    }
}

export default ConsoleContent;
