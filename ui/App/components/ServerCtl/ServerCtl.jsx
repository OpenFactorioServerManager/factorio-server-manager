import React from 'react';

class ServerCtl extends React.Component {
    constructor(props) {
        super(props);
        this.startServer = this.startServer.bind(this);
        this.incrementAutosave = this.incrementAutosave.bind(this);
        this.decrementAutosave = this.decrementAutosave.bind(this);

        this.incrementAutosaveSlots = this.incrementAutosaveSlots.bind(this);
        this.decrementAutosaveSlots = this.decrementAutosaveSlots.bind(this);

        this.incrementPort = this.incrementPort.bind(this);
        this.decrementPort = this.decrementPort.bind(this);
        
        this.incrementLatency = this.incrementLatency.bind(this);
        this.decrementLatency = this.decrementLatency.bind(this);

        this.toggleAllowCmd = this.toggleAllowCmd.bind(this);
        this.toggleP2P = this.toggleP2P.bind(this);
        this.toggleAutoPause = this.toggleAutoPause.bind(this);

        this.state = {
            latency: 100,
            autosaveInterval: 5,
            autosaveSlots: 10,
            port: 34197,
            disallowCmd: false,
            peer2peer: false,
            autoPause: false,
        }
    }

    startServer(e) {
        let serverSettings = {
            latency: Number(this.refs.latency.value), 
            autosave_interval: Number(this.refs.autosaveInterval.value),
            autosave_slots: Number(this.refs.autosaveSlots.value),
            port: Number(this.refs.port.value),
            disallow_cmd: this.refs.allowCmd.checked,
            peer2peer: this.refs.p2p.checked,
            auto_pause: this.refs.autoPause.checked,
        }
        console.log(serverSettings);
        $.ajax({
            type: "POST",
            url: "/api/server/start",
            dataType: "json",
            data: JSON.stringify(serverSettings),
            success: (resp) => {
                alert(resp)
            }
        })
        e.preventDefault();
    }

    incrementAutosave() {
        let saveInterval = this.state.autosaveInterval + 1;
        this.setState({autosaveInterval: saveInterval})
    }

    decrementAutosave() {
        let saveInterval = this.state.autosaveInterval - 1;
        this.setState({autosaveInterval: saveInterval})
    }

    incrementAutosaveSlots() {
        let saveSlots = this.state.autosaveSlots + 1;
        this.setState({autosaveSlots: saveSlots})
    }

    decrementAutosaveSlots() {
        let saveSlots = this.state.autosaveSlots - 1;
        this.setState({autosaveSlots: saveSlots})
    }

    incrementPort() {
        let port = this.state.port + 1;
        this.setState({port: port})
    }

    decrementPort() {
        let port = this.state.port - 1;
        this.setState({port: port})
    }
    
    incrementLatency() {
        let latency = this.state.latency + 1;
        this.setState({latency: latency})
    }

    decrementLatency() {
        let latency= this.state.latency- 1;
        this.setState({latency: latency})
    }

    toggleAllowCmd() {
        let cmd = !this.state.disallowCmd
        this.setState({disallowCmd: cmd})
    }

    toggleP2P() {
        let p2p = !this.state.peer2peer;
        this.setState({peer2peer: p2p})
    }

    toggleAutoPause() {
        let pause = !this.state.autoPause;
        this.setState({autoPause: pause})
    }

    render() {
        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Server Control</h3>
                </div>
                
                <div className="box-body">

                <form action="" onSubmit={this.startServer}>
                    <label for="latency">Server latency setting (ms)</label>
                    <div id="latency" className="input-group">
                        <input ref="latency" name="latency" id="latency" type="text" className="form-control" onchange={this.state.latency} value={this.state.latency} placeholder={this.state.latency} />
                        <div className="input-group-btn">
                        <button type="button" className="btn btn-primary" onClick={this.incrementLatency}><i className="fa fa-arrow-up"></i></button>
                        <button type="button" className="btn btn-primary" onClick={this.decrementLatency}><i className="fa fa-arrow-down"></i></button>
                        </div>
                    </div>
                    <label for="autosaveInterval">Autosave Interval (mins)</label>
                    <div id="autosaveInterval" className="input-group">
                        <input ref="autosaveInterval" name="autosaveInterval" id="autosaveInterval" type="text" className="form-control" onchange={this.state.autosaveInterval} value={this.state.autosaveInterval} placeholder={this.state.autosaveInterval} />
                        <div className="input-group-btn">
                        <button type="button" className="btn btn-primary" onClick={this.incrementAutosave}><i className="fa fa-arrow-up"></i></button>
                        <button type="button" className="btn btn-primary" onClick={this.decrementAutosave}><i className="fa fa-arrow-down"></i></button>
                        </div>
                    </div>
                    <label for="autosaveSlots">Autosave Slots</label>
                    <div id="autosaveSlots" className="input-group">
                        <input ref="autosaveSlots" name="autosaveSlots" id="autosaveSlots" type="text" className="form-control" onChange={this.state.autosaveSlots} value={this.state.autosaveSlots} placeholder={this.state.autosaveSlots} />
                        <div className="input-group-btn">
                        <button type="button" className="btn btn-primary" onClick={this.incrementAutosaveSlots}><i className="fa fa-arrow-up"></i></button>
                        <button type="button" className="btn btn-primary" onClick={this.decrementAutosaveSlots}><i className="fa fa-arrow-down"></i></button>
                        </div>
                    </div>
                    <label for="port">Factorio Server Port</label>
                    <div id="port" className="input-group">
                        <input ref="port" name="port" id="port" type="text" className="form-control" onChange={this.state.port} value={this.state.port} placeholder={this.state.port} />
                        <div className="input-group-btn">
                        <button type="button" className="btn btn-primary" onClick={this.incrementPort}><i className="fa fa-arrow-up"></i></button>
                        <button type="button" className="btn btn-primary" onClick={this.decrementPort}><i className="fa fa-arrow-down"></i></button>
                        </div>
                    </div>
                    <div class="form-group">
                        <div class="checkbox">
                            <input ref="autoPause" type="checkbox" onClick={this.toggleAutoPause} />
                            <label>
                            Auto Pause when no players connected
                            </label>
                        </div>

                        <div class="checkbox">
                            <input ref="p2p" type="checkbox" onClick={this.toggleP2P} />
                            <label>
                            Peer to peer connection method
                            </label>
                        </div>

                        <div class="checkbox">
                            <input ref="allowCmd" type="checkbox" onClick={this.toggleAllowCmd} />
                            <label>
                            Allow commands on the server
                            </label>
                        </div>
                    </div>

                    <div class="form-group">
                    <select ref="savefile" class="form-control">
                        {this.props.saves.map( (save, i) => {
                            return(
                                <option key={save.name} value={save.name}>{save.name}</option>
                            )                                
                            
                        })}
                    </select>
                    <label>Select Save File</label>
                    </div>

                    <button className="btn btn-block btn-success" type="submit">Start Factorio Server</button>
                </form>

                </div>
            </div>

        )
    }
}

export default ServerCtl
