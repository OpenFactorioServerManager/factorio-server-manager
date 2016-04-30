import React from 'react';
import {browserHistory} from 'react-router';
import Header from './components/Header.jsx';
import Sidebar from './components/Sidebar.jsx';
import Footer from './components/Footer.jsx';
import HiddenSidebar from './components/HiddenSidebar.jsx';


class App extends React.Component {
    constructor(props) {
        super(props);
        this.checkLogin = this.checkLogin.bind(this);
        this.facServStatus = this.facServStatus.bind(this);
        this.getSaves = this.getSaves.bind(this);
        this.getStatus = this.getStatus.bind(this);
        this.state = {
            serverRunning: "stopped",
            serverStatus: {},
            saves: [],
            loggedIn: false,
            username: "",
        }
    }

    componentDidMount() {
        this.checkLogin();
        setTimeout(() => {
            if (!this.state.loggedIn) {
                browserHistory.push("/login");
            }
        }, 1000);
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
                this.setState({saves: data.data})
            },
            error: (xhr, status, err) => {
                console.log('api/mods/list', status, err.toString());
            }
        })
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
        // render main application, if not logged in show Not logged in message
        // if logged in show application
        var resp;
        if (this.state.loggedIn) {
            var resp = 
                <div>
                    <Header 
                        username={this.state.username}
                        loggedIn={this.state.loggedIn}
                    />

                    <Sidebar 
                        serverStatus={this.facServStatus}
                        serverRunning={this.state.serverRunning}
                    />
                    
                    {React.cloneElement(
                        this.props.children,
                        {message: "",
                        facServStatus: this.facServStatus,
                        serverStatus: this.state.serverStatus,
                        getStatus: this.getStatus,
                        saves: this.state.saves,
                        getSaves: this.getSaves,
                        username: this.state.username}
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
