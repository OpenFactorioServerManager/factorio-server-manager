import React from 'react';
import {IndexLink} from 'react-router';
import Settings from './Config/Settings.jsx';
import ServerSettings from './Config/ServerSettings.jsx';

class ConfigContent extends React.Component {
    constructor(props) {
        super(props);
        this.getConfig = this.getConfig.bind(this);
        this.getServerSettings = this.getServerSettings.bind(this);
        this.state = {
            config: {},
            serverSettings: {}
        }
    }

    componentDidMount() {
        this.getConfig();
        this.getServerSettings();
        console.log(this.state.serverSettings);
    }
    
    getConfig() {
        $.ajax({
            url: "/api/config",
            dataType: "json",
            success: (resp) => {
                if (resp.success === true) {
                    this.setState({config: resp.data})
                }
            },
            error: (xhr, status, err) => {
                console.log('/api/config/get', status, err.toString());
            }
        });
    }

    getServerSettings() {
        $.ajax({
            url: "/api/settings",
            dataType: "json",
            success: (resp) => {
                if (resp.success === true) {
                    this.setState({serverSettings: resp.data})
                }
                console.log(this.state.serverSettings);
            },
            error: (xhr, status, err) => {
                console.log('/api/settings/get', status, err.toString());
            }
        });
    }

    render() {
        return(
            <div className="content-wrapper">
                <section className="content-header">
                <h1>
                    Config 
                    <small>Manage game configuration</small>
                </h1>
                <ol className="breadcrumb">
                    <li><IndexLink to="/"><i className="fa fa-dashboard"></i>Server Control</IndexLink></li>
                    <li className="active">Here</li>
                </ol>
                </section>
                
                <section className="content">
                    <div className="box">
                        <div className="box-header">
                            <h3 className="box-title">Server Settings</h3>
                        </div>

                        <div className="box-body">
                        <div className="row">
                            <div className="col-md-6">
                                    <div className="server-settings-section">
                                        <div className="table-responsive">
                                            <table>
                                                    <thead>
                                                        <tr>
                                                            <th>Setting name</th>
                                                            <th>Setting value</th>
                                                        </tr>
                                                    </thead>
                                                    <tbody>
                                                    {Object.keys(this.state.serverSettings).map(function(key) {
                                                        var setting = this.state.serverSettings[key]
                                                        console.log(setting)
                                                        return(
                                                            <ServerSettings 
                                                                name={key}
                                                                setting={setting}
                                                            />
                                                        )
                                                    }, this)}
                                                    </tbody>
                                            </table>
                                        </div>
                                    </div>
                            </div>
                        </div>
                        </div>
                    </div>
                </section>

                <section className="content">
                    <div className="box">
                        <div className="box-header">
                            <h3 className="box-title">Game Configuration</h3>
                        </div>
                        
                        <div className="box-body">
                        <div className="row">
                            <div className="col-md-6">
                            {Object.keys(this.state.config).map(function(key) {
                                var conf = this.state.config[key]
                                return(
                                    <div className="settings-section" key={key}>
                                    <h3>{key}</h3>
                                        <div className="table-responsive">
                                        <table className="table table-striped">
                                            <thead>
                                                <tr>
                                                    <th>Setting name</th>
                                                    <th>Setting value</th>
                                                </tr>
                                            </thead>
                                                <Settings
                                                    section={key}
                                                    config={conf}    
                                                />
                                        </table>
                                        </div>
                                    </div>
                                )
                            }, this)}
                            </div>
                        </div>
                        </div>
                    </div>
                </section>


            </div>
        )
    }
}

export default ConfigContent
