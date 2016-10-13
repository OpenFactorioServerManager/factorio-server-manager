import React from 'react';
import {IndexLink} from 'react-router';
import Settings from './Config/Settings.jsx';

class ConfigContent extends React.Component {
    constructor(props) {
        super(props);
        this.getConfig = this.getConfig.bind(this);
        this.getServerSettings = this.getServerSettings.bind(this);
        this.updateServerSettings = this.updateServerSettings.bind(this);
        this.handleServerSettingsChange = this.handleServerSettingsChange.bind(this);
        this.state = {
            config: {},
            serverSettings: {}
        }
    }

    componentDidMount() {
        this.getConfig();
        this.getServerSettings();
    }

    handleServerSettingsChange(name, e) {
        var change = this.state.serverSettings;
        change[name] = e.target.value;
        this.setState({serverSettings: change});
        console.log(this.state);
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
            },
            error: (xhr, status, err) => {
                console.log('/api/settings/get', status, err.toString());
            }
        });
    }
    
    updateServerSettings(e) {
        e.preventDefault();
        var serverSettingsJSON = JSON.stringify(this.state.serverSettings)
        $.ajax({
            url: "/api/settings/update",
            datatype: "json",
            type: "POST",
            data: serverSettingsJSON,
            success: (data) => {
                console.log(data);
                if (data.success === true) {
                    console.log("settings updated") 
                }
            }
        })
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
                            <div className="col-md-10">
                                <div className="server-settings-section">
                                    <div className="table-responsive">
                                        <form ref="settingsForm" className="form-horizontal" onSubmit={this.updateServerSettings}>
                                            {Object.keys(this.state.serverSettings).map(function(key) {
                                                var setting = this.state.serverSettings[key]
                                                var setting_key = key.replace(/_/g, " ")
                                                return(
                                                <div className="form-group">
                                                    <label for={key} className="control-label col-md-3">{setting_key}</label>
                                                    <div className="col-md-6">
                                                        <input 
                                                            ref={key} 
                                                            id={key} 
                                                            className="form-control" 
                                                            defaultValue={setting} 
                                                            type="text" 
                                                            onChange={this.handleServerSettingsChange.bind(this, key)}
                                                        />
                                                    </div>
                                                </div>
                                                )
                                            }, this)}
                                            <div className="col-xs-6">
                                                <div className="form-group">
                                                    <input className="form-control btn btn-success" type="submit" ref="button" value="Update Settings" />
                                                </div>
                                            </div>
                                        </form>
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
                            <div className="col-md-10">
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
