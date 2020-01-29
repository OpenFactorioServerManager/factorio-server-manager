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

        this.gameBindIPRef = React.createRef();
        this.saveFileRef = React.createRef();
        this.portRef = React.createRef();
    }

    startServer(e) {
        e.preventDefault();
        let serverSettings = {
            bindip: this.gameBindIPRef.current.value,
            savefile: this.saveFileRef.current.value,
            port: Number(this.portRef.current.value),
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
                        icon: "success"
                    });
                } else {
                    ReactSwalNormal.fire({
                        title: "Error starting Factorio server",
                        text: resp.data,
                        icon: "error"
                    });
                }
            }
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
        this.portRef.current.value = Number(this.portRef.current.value) + 1;
    }

    decrementPort() {
        this.portRef.current.value = Number(this.portRef.current.value - 1);
    }

    render() {
        return (
            <div id="serverCtl" className="card">
                <div className="card-header">
                    <h3 className="card-title">Server Control</h3>
                </div>

                <div className="card-body">
                    <form action="" onSubmit={this.startServer}>
                        <div className="row">
                            <div className="col-md-4">
                                <button className="btn btn-block btn-success" type="submit">
                                    <FontAwesomeIcon icon="play" className="fa-fw"/>Start Factorio Server
                                </button>
                            </div>

                            <div className="col-md-4">
                                <button className="btn btn-block btn-warning" type="button"
                                        onClick={this.stopServer}>
                                    <FontAwesomeIcon icon="stop" className="fa-fw"/>Stop &amp; Save Factorio Server
                                </button>
                            </div>

                            <div className="col-md-4">
                                <button className="btn btn-block btn-danger" type="button"
                                        onClick={this.killServer}>
                                    <FontAwesomeIcon icon="close" className="fa-fw"/>Stop Factorio Server without
                                    Saving
                                </button>
                            </div>
                        </div>

                        <hr/>
                        <div className="form-group">
                            <label>Select Save File</label>
                            <select ref={this.saveFileRef} className="form-control">
                                {this.props.saves.map((save, i) => {
                                    return (
                                        <option key={save.name} value={save.name}>{save.name}</option>
                                    )

                                })}
                            </select>
                        </div>

                        <div className="form-group">
                            <label htmlFor="gameBindIP">Factorio Server IP</label>
                            <div className="input-group">
                                <input ref={this.gameBindIPRef}
                                       name="gameBindIP"
                                       id="gameBindIP"
                                       type="text"
                                       className="form-control"
                                       defaultValue="0.0.0.0"
                                       placeholder="0.0.0.0"/>
                            </div>
                        </div>

                        <div className="form-group">
                            <label htmlFor="port">Factorio Server Port</label>
                            <div className="input-group">
                                <input ref={this.portRef}
                                       name="port"
                                       id="port"
                                       type="text"
                                       className="form-control"
                                       defaultValue="34197"
                                       placeholder="34197"
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

