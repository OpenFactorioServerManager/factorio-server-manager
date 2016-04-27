import React from 'react';

class HiddenSidebar extends React.Component {
    constructor(props) {
        super(props);
    }

    capitalizeFirstLetter(string) {
        return string.charAt(0).toUpperCase() + string.slice(1);
    }

    render() {
        return(
            <aside className="control-sidebar control-sidebar-dark">
                <div className="table-responsive">
                <table className="table table-border">
                    <thead>
                        <tr>
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
            <div className="control-sidebar-bg"></div>
            </aside>
        )
    }
}

HiddenSidebar.propTypes = {
    serverStatus: React.PropTypes.object.isRequired,
}


export default HiddenSidebar
