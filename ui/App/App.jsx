import React from 'react';
import Header from './components/Header.jsx';
import Sidebar from './components/Sidebar.jsx';
import Footer from './components/Footer.jsx';
import HiddenSidebar from './components/HiddenSidebar.jsx';


class App extends React.Component {
    constructor(props) {
        super(props);
        this.facServStatus = this.facServStatus.bind(this);
        this.getSaves = this.getSaves.bind(this);
        this.getStatus = this.getStatus.bind(this);
        this.state = {
            serverRunning: "stopped",
            serverStatus: {},
            saves: [],
        }
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
        return(
            <div className="wrapper" style={{height: "100%"}}>

                <Header />

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
                     getSaves: this.getSaves}
                )}

                <Footer />

                <HiddenSidebar 
                    serverStatus={this.state.serverStatus}
                />

            </div>
        )
    }
}

export default App
