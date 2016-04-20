import React from 'react';

class LogLines extends React.Component {
    updateLog() {
        this.props.getLastLog();
    }

    render() {
        this.props.log.reverse();
        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Factorio Log</h3>
                </div>
                <div className="box-body">
                <input className="btn btn-default" type='button' onClick={this.updateLog.bind(this)} value="Refresh" />
                <h5>Latest log line at the top</h5>
                <samp>
                    {this.props.log.map ( (line, i) => {
                        return(
                            <p key={i}>{line}</p>
                        )                                        
                    })}
                </samp>
                </div>
            </div>
        )
    }
}

LogLines.propTypes = {
    log: React.PropTypes.array.isRequired,
    getLastLog: React.PropTypes.func.isRequired
}

export default LogLines
