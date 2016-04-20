import React from 'react';
import LogLines from './Logs/LogLines.jsx';

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
                this.setState({log: data})
            },
            error: (xhr, status, err) => {
                console.log('api/mods/list', status, err.toString());
            }
        })
    }

    render() {
        return(
            <div className="content-wrapper">
                <section className="content-header">
                <h1>
                    Logs
                    <small>Optional description</small>
                </h1>
                <ol className="breadcrumb">
                    <li><a href="#"><i className="fa fa-dashboard"></i> Level</a></li>
                    <li className="active">Here</li>
                </ol>
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
