import React from 'react';

class ServerCtl extends React.Component {
    constructor(props) {
        super(props);
        this.startServer = this.startServer.bind(this);
        this.stopServer = this.stopServer.bind(this);
        this.killServer = this.killServer.bind(this);

        this.incrementPort = this.incrementPort.bind(this);
        this.decrementPort = this.decrementPort.bind(this);

        this.state = {
            savefile: "",
            port: 34197,
        }
    }

    startServer(e) {
        e.preventDefault();
        let serverSettings = {
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
                    swal("Factorio Server Started", resp.data)
                } else {
                    var err = "Error starting Factorio Server: " + resp.data
                    swal("Error", err, "error")
                }
            }
        })
        this.setState({
            savefile: this.refs.savefile.value,
            port: Number(this.refs.port.value),
        })
    }

    stopServer(e) {
        $.ajax({
            type: "GET",
            url: "/api/server/stop",
            dataType: "json",
            success: (resp) => {
                this.props.facServStatus();
                this.props.getStatus();
                console.log(resp)
                swal(resp.data)
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
                console.log(resp)
                swal(resp.data)
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
        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Server Control</h3>
                </div>
                
                <div className="box-body">

                    <form action="" onSubmit={this.startServer}>
                        <div className="form-group">
                            <div className="row">
                                <div className="col-md-4">
                                <button className="btn btn-block btn-success" type="submit"><i className="fa fa-play fa-fw"></i>Start Factorio Server</button>
                                </div>
                                
                                <div className="col-md-4">
                                <button className="btn btn-block btn-warning" type="button" onClick={this.stopServer}><i className="fa fa-stop fa-fw"></i>Stop &amp; Save Factorio Server</button>
                                </div>

                                <div className="col-md-4">
                                <button className="btn btn-block btn-danger" type="button" onClick={this.killServer}><i className="fa fa-close fa-fw"></i>Stop Factorio Server without Saving</button>
                                </div>
                            </div>

                            <hr />
                            <label>Select Save File</label>
                            <select ref="savefile" className="form-control">
                                {this.props.saves.map( (save, i) => {
                                    return(
                                        <option key={save.name} value={save.name}>{save.name}</option>
                                    )                                
                                    
                                })}
                            </select>
                        </div>

                        <div className="box box-success collapsed-box">
                                <button type="button" className="btn btn-box-tool" data-widget="collapse">
                                    <div className="box-header with-border">
                                    <i className="fa fa-plus fa-fw"></i><h4 className="box-title">Advanced</h4>
                                    </div>
                                </button>
                                <div className="box-body" style={{display: "none"}}>
                                <label htmlFor="port">Factorio Server Port</label>
                                <div id="port" className="input-group">
                                    <input ref="port" name="port" id="port" type="text" className="form-control" onChange={this.state.port} value={this.state.port} placeholder={this.state.port} />
                                    <div className="input-group-btn">
                                    <button type="button" className="btn btn-primary" onClick={this.incrementPort}><i className="fa fa-arrow-up"></i></button>
                                    <button type="button" className="btn btn-primary" onClick={this.decrementPort}><i className="fa fa-arrow-down"></i></button>
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
    facServStatus: React.PropTypes.func.isRequired,
    getStatus: React.PropTypes.func.isRequired,
}

export default ServerCtl

