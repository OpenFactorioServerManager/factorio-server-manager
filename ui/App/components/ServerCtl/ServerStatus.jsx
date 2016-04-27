import React from 'react';

class ServerStatus extends React.Component {
    constructor(props) {
        super(props);
    }

    capitalizeFirstLetter(string) {
        return string.charAt(0).toUpperCase() + string.slice(1);
    }

    render() {
        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Server Status</h3>
                </div>
                
                <div className="box-body">
                    <div className="table-responsive">
                    <table className="table table-striped">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Status</th>
                            </tr>
                        </thead>
                        <tbody>
                            {Object.keys(this.props.serverStatus).map(function(key) {
                                return(
                                    <tr key={key}>
                                        <td>{this.capitalizeFirstLetter(key)}</td>
                                        <td>{this.props.serverStatus[key]}</td>
                                    </tr>
                                )                                                  
                            }, this)}        
                        </tbody>
                    </table>
                    </div>
                </div>
            </div>
        )
    }
}

ServerStatus.propTypes = {
    serverStatus: React.PropTypes.object.isRequired,
}


export default ServerStatus
