import React from 'react';
import PropTypes from 'prop-types';
import {ReactSwalNormal} from 'Utilities/customSwal';
import FontAwesomeIcon from "../FontAwesomeIcon";

class ServerCtl extends React.Component {
    constructor(props) {
        super(props);
        this.startServer = this.startServer.bind(this);
        this.stopServer = this.stopServer.bind(this);
        this.killServer = this.killServer.bind(this);

        this.incrementPort = this.incrementPort.bind(this);
        this.decrementPort = this.decrementPort.bind(this);

        this.state = {
            gameBindIP: "0.0.0.0",
            savefile: "",
            port: 34197,
        }
    }

    startServer(e) {
        e.preventDefault();
        let serverSettings = {
            bindip: this.refs.gameBindIP.value,
            savefile: this.refs.savefile.value,
            port: Number(this.refs.port.value),
        }
        $.ajax({
            type: "POST",
            url: "/api/server/start",
            dataType: "json",
            data: JSON.stringify(serverSettings),
            success: (resp) => {
                this.props.facServStatus();
                this.props.getStatus();
                if (resp.success === true) {
                    ReactSwalNormal.fire({
                        title: "Factorio server started",
                        text: resp.data,
                        type: "success"
                    });
                } else {
                    ReactSwalNormal.fire({
                        title: "Error starting Factorio server",
                        text: resp.data,
                        type: "error"
                    });
                }
            }
        });

        this.setState({
            savefile: this.refs.savefile.value,
            port: Number(this.refs.port.value),
        });
    }

    stopServer(e) {
        $.ajax({
            type: "GET",
            url: "/api/server/stop",
            dataType: "json",
            success: (resp) => {
                this.props.facServStatus();
                this.props.getStatus();

                ReactSwalNormal.fire({
                    title: resp.data
                });
            }
        });
        e.preventDefault();
    }

    killServer(e) {
        $.ajax({
            type: "GET",
            url: "/api/server/kill",
            dataType: "json",
            success: (resp) => {
                this.props.facServStatus();
                this.props.getStatus();

                ReactSwalNormal.fire({
                    title: resp.data
                });
            }
        });
        e.preventDefault();
    }

    incrementPort() {
        let port = this.state.port + 1;
        this.setState({port: port})
    }

    decrementPort() {
        let port = this.state.port - 1;
        this.setState({port: port})
    }

    render() {
        return (
            <div id="serverCtl" className="box">
                <div className="box-header">
                    <h3 className="box-title">Server Control</h3>
                </div>

                <div className="box-body">
                    <form action="" onSubmit={this.startServer}>
                        <div className="form-group">
                            <div className="row">
                                <div className="col-md-4">
                                    <button className="btn btn-block btn-success" type="submit">
                                        <FontAwesomeIcon icon="play" className="fa-fw"/>Start Factorio Server
                                    </button>
                                </div>

                                <div className="col-md-4">
                                    <button className="btn btn-block btn-warning" type="button" onClick={this.stopServer}>
                                        <FontAwesomeIcon icon="stop" className="fa-fw"/>Stop &amp; Save Factorio Server
                                    </button>
                                </div>

                                <div className="col-md-4">
                                    <button className="btn btn-block btn-danger" type="button" onClick={this.killServer}>
                                        <FontAwesomeIcon icon="close" className="fa-fw"/>Stop Factorio Server without Saving
                                    </button>
                                </div>
                            </div>

                            <hr/>
                            <label>Select Save File</label>
                            <select ref="savefile" className="form-control">
                                {this.props.saves.map((save, i) => {
                                    return (
                                        <option key={save.name} value={save.name}>{save.name}</option>
                                    )

                                })}
                            </select>
                        </div>

                        <div className="box box-success advanced">
                            <button type="button"
                                    className="btn btn-box-tool"
                                    data-toggle="collapse"
                                    data-target="#serverCtlAdvanced"
                                    aria-expanded="false"
                                    aria-controls="serverCtlAdvanced"
                            >
                                <div className="box-header with-border">
                                    <FontAwesomeIcon icon="plus" className="fa-fw"/>
                                    <div className="box-title">Advanced</div>
                                </div>
                            </button>
                            <div id="serverCtlAdvanced" className="box-body collapse">
                                <label htmlFor="port">Factorio Server IP</label>
                                <div id="port" className="input-group">
                                    <input ref="gameBindIP"
                                           name="gameBindIP"
                                           id="gameBindIP"
                                           type="text"
                                           className="form-control"
                                           defaultValue={this.state.gameBindIP}
                                           placeholder={this.state.gameBindIP}/>
                                </div>
                                <label htmlFor="port">Factorio Server Port</label>
                                <div id="port" className="input-group">
                                    <input ref="port"
                                           name="port"
                                           id="port"
                                           type="text"
                                           className="form-control"
                                           defaultValue={this.state.port}
                                           placeholder={this.state.port}
                                    />
                                    <div className="input-group-btn">
                                        <button type="button" className="btn btn-primary" onClick={this.incrementPort}>
                                            <FontAwesomeIcon icon="arrow-up"/>
                                        </button>
                                        <button type="button" className="btn btn-primary" onClick={this.decrementPort}>
                                            <FontAwesomeIcon icon="arrow-down"/>
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>

        )
    }
}

ServerCtl.propTypes = {
    facServStatus: PropTypes.func.isRequired,
    getStatus: PropTypes.func.isRequired,
}

export default ServerCtl

