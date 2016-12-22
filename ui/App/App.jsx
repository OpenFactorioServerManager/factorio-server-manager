import React from 'react';
import {browserHistory} from 'react-router';
import Header from './components/Header.jsx';
import Sidebar from './components/Sidebar.jsx';
import Footer from './components/Footer.jsx';
import HiddenSidebar from './components/HiddenSidebar.jsx';
import Socket from '../socket.js';

class App extends React.Component {
    constructor(props) {
        super(props);
        this.checkLogin = this.checkLogin.bind(this);
        this.flashMessage = this.flashMessage.bind(this);
        this.facServStatus = this.facServStatus.bind(this);
        this.getSaves = this.getSaves.bind(this);
        this.getStatus = this.getStatus.bind(this);
        this.connectWebSocket = this.connectWebSocket.bind(this);
        this.state = {
            serverRunning: "stopped",
            serverStatus: {},
            saves: [],
            loggedIn: false,
            username: "",
            messages: [],
            showMessage: false,
        }
    }

    componentDidMount() {
        this.checkLogin();
        // Wait 1 second before redirecting to login page
        setTimeout(() => {
            if (!this.state.loggedIn) {
                browserHistory.push("/login");
            }
        }, 1000);
        this.connectWebSocket();
    }
    
    connectWebSocket() {
        var ws_scheme = window.location.protocol == "https:" ? "wss" : "ws";
        let ws = new WebSocket(ws_scheme + "://" + window.location.host + "/ws");
        let socket = this.socket = new Socket(ws);
    }

    flashMessage(message) {
        var m = this.state.messages;
        m.push(message);
        this.setState({messages: m, showMessage: true});
    }

    checkLogin() {
        $.ajax({
            url: "/api/user/status",
            dataType: "json",
            success: (data) => {
                if (data.success === true) {
                    this.setState({loggedIn: true,
                        username: data.data.Username})
                }
            }
        })
    }

    facServStatus() {
        $.ajax({
            url: "/api/server/status",
            dataType: "json",
            success: (data) => {
                this.setState({serverRunning: data.data.status})
            }
        })
    }

    getSaves() {
        $.ajax({
            url: "/api/saves/list",
            dataType: "json",
            success: (data) => {
                if (data.success === true) {
                    this.setState({saves: data.data})
                } else {
                    this.setState({saves: []})
                }
            },
            error: (xhr, status, err) => {
                console.log('api/saves/list', status, err.toString());
            }
        })

        if (!this.state.saves) {
            this.setState({saves:[]})
        }
    }

    getStatus() {
        $.ajax({
            url: "/api/server/status",
            dataType: "json",
            success: (data) => {
                this.setState({serverStatus: data.data})
            },
            error: (xhr, status, err) => {
                console.log('api/server/status', status, err.toString());
            }
        })
    }

    render() {
        // render main application, 
        // if logged in show application
        // if not logged in show Not logged in message
        var resp;
        if (this.state.loggedIn) {
            var resp = 
                <div>
                    <Header 
                        username={this.state.username}
                        loggedIn={this.state.loggedIn}
                        messages={this.state.messages}
                    />

                    <Sidebar 
                        serverStatus={this.facServStatus}
                        serverRunning={this.state.serverRunning}
                    />
                    
                    // Render react-router components and pass in props
                    {React.cloneElement(
                        this.props.children,
                        {message: "",
                        messages: this.state.messages,
                        flashMessage: this.flashMessage,
                        facServStatus: this.facServStatus,
                        serverStatus: this.state.serverStatus,
                        getStatus: this.getStatus,
                        saves: this.state.saves,
                        getSaves: this.getSaves,
                        username: this.state.username,
                        socket: this.socket}
                    )}

                    <Footer />

                    <HiddenSidebar 
                        serverStatus={this.state.serverStatus}
                        username={this.state.username}
                        loggedIn={this.state.loggedIn}
                        checkLogin={this.checkLogin}
                    />
                </div>
        } else {
            var resp = <div><p>Not Logged in</p></div>;
        }

        return(
            <div className="wrapper">
            {resp}
            </div>
        )
    }
}

export default App
