import React from 'react';
import {Switch, Route, withRouter} from 'react-router-dom';
import Header from './components/Header.jsx';
import Sidebar from './components/Sidebar.jsx';
import Footer from './components/Footer.jsx';
import Socket from '../socket.js';
import Index from "./components/Index";
import UsersContent from "./components/UsersContent";
import ModsContent from "./components/ModsContent";
import LogsContent from "./components/LogsContent";
import SavesContent from "./components/SavesContent";
import ConfigContent from "./components/ConfigContent";
import ConsoleContent from "./components/ConsoleContent";

class App extends React.Component {
    constructor(props) {
        super(props);
        this.checkLogin = this.checkLogin.bind(this);
        this.flashMessage = this.flashMessage.bind(this);
        this.facServStatus = this.facServStatus.bind(this);
        this.getSaves = this.getSaves.bind(this);
        this.getStatus = this.getStatus.bind(this);
        this.connectWebSocket = this.connectWebSocket.bind(this);
        this.getFactorioVersion = this.getFactorioVersion.bind(this);

        this.state = {
            serverRunning: "stopped",
            serverStatus: {},
            factorioVersion: "",
            saves: [],
            loggedIn: false,
            username: "",
            messages: [],
            showMessage: false,
        }
    }

    componentDidMount() {
        this.checkLogin();
    }

    connectWebSocket() {
        let ws_scheme = window.location.protocol == "https:" ? "wss" : "ws";
        let ws = new WebSocket(ws_scheme + "://" + window.location.host + "/ws");
        this.socket = new Socket(ws);
    }

    flashMessage(message) {
        var m = this.state.messages;
        m.push(message);
        this.setState({messages: m, showMessage: true});
    }

    checkLogin() {
        $.ajax({
            url: "/api/user/status",
            type: "GET",
            dataType: "json",
            success: (data) => {
                if (data.success === true) {
                    this.setState({
                        loggedIn: true,
                        username: data.data.Username
                    });

                    this.connectWebSocket();
                    this.getFactorioVersion(); //Init serverStatus, so i know, which factorio-version is installed
                } else {
                    this.props.history.push("/login");
                }
            },
            error: () => {
                this.props.history.push("/login");
            }
        })
    }

    facServStatus() {
        $.ajax({
            url: "/api/server/status",
            dataType: "json",
            success: (data) => {
                this.setState({
                    serverRunning: data.data.status
                })
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
        });

        if (!this.state.saves) {
            this.setState({saves:[]});
        }
    }

    getStatus() {
        $.ajax({
            url: "/api/server/status",
            dataType: "json",
            success: (data) => {
                this.setState({
                    serverStatus: data.data
                })
            },
            error: (xhr, status, err) => {
                console.log('api/server/status', status, err.toString());
            }
        })
    }

    getFactorioVersion() {
        $.ajax({
            url: "/api/server/facVersion",
            // dataType: "json",
            success: (data) => {
                console.log(data);
                this.setState({
                    factorioVersion: data.data.base_mod_version
                });
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
        let appProps = {
            message: "",
            messages: this.state.messages,
            flashMessage: this.flashMessage,
            facServStatus: this.facServStatus,
            serverStatus: this.state.serverStatus,
            factorioVersion: this.state.factorioVersion,
            getStatus: this.getStatus,
            saves: this.state.saves,
            getSaves: this.getSaves,
            username: this.state.username,
            socket: this.socket
        };

        let resp;
        if (this.state.loggedIn) {
            resp =
                <div className="wrapper">
                    <Header
                        username={this.state.username}
                        loggedIn={this.state.loggedIn}
                        messages={this.state.messages}
                    />

                    <Sidebar
                        serverStatus={this.facServStatus}
                        serverRunning={this.state.serverRunning}
                    />

                    {/*Render react-router components and pass in props*/}
                    <Switch>
                        <Route path="/server" render={(props) => {return <Index {...props} {...appProps}/>}}/>
                        <Route path="/settings" render={(props) => {return <UsersContent {...props} {...appProps}/>}}/>
                        <Route path="/mods" render={(props) => {return <ModsContent {...props} {...appProps}/>}}/>
                        <Route path="/logs" render={(props) => {return <LogsContent {...props} {...appProps}/>}}/>
                        <Route path="/saves" render={(props) => {return <SavesContent {...props} {...appProps}/>}}/>
                        <Route path="/config" render={(props) => {return <ConfigContent {...props} {...appProps}/>}}/>
                        <Route path="/console" render={(props) => {return <ConsoleContent {...props} {...appProps}/>}}/>
                        <Route exact path="/" render={(props) => {return <Index {...props} {...appProps} />}}/>
                    </Switch>

                    <Footer />
                </div>
        } else {
            resp = <div><p>Not Logged in</p></div>;
        }

        return resp;
    }
}

export default withRouter(App);
