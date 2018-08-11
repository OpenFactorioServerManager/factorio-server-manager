import React from 'react';
import {Link} from 'react-router-dom';
import LogLines from './Logs/LogLines.jsx';
import FontAwesomeIcon from "./FontAwesomeIcon";

class LogsContent extends React.Component {
    constructor(props) {
        super(props);
        this.componentDidMount = this.componentDidMount.bind(this);
        this.getLastLog = this.getLastLog.bind(this);
        this.state = {
            log: []
        }
    }

    componentDidMount() {
        this.getLastLog();
    }

    getLastLog() {
        $.ajax({
            url: "/api/log/tail",
            dataType: "json",
            success: (data) => {
                this.setState({log: data.data})
            },
            error: (xhr, status, err) => {
                console.log('api/mods/list', status, err.toString());
            }
        })
    }

    render() {
        return (
            <div className="content-wrapper">
                <section className="content-header">
                    <h1>
                        Logs
                        <small>Analyze Factorio Logs</small>

                        <small className="float-sm-right">
                            <ol className="breadcrumb">
                                <li className="breadcrumb-item">
                                    <Link to="/"><FontAwesomeIcon icon="tachometer-alt"/>Server Control</Link>
                                </li>
                                <li className="breadcrumb-item active">
                                    <FontAwesomeIcon icon="file-alt" prefix="far"/>Logs
                                </li>
                            </ol>
                        </small>
                    </h1>
                </section>

                <section className="content">
                    <LogLines
                        getLastLog={this.getLastLog}
                        {...this.state}
                    />

                </section>
            </div>
        )
    }
}

export default LogsContent
