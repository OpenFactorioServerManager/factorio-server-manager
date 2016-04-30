import React from 'react';
import {Link, browserHistory} from 'react-router';

class HiddenSidebar extends React.Component {
    constructor(props) {
        super(props);
        this.onLogout = this.onLogout.bind(this);
    }

    capitalizeFirstLetter(string) {
        return string.charAt(0).toUpperCase() + string.slice(1);
    }

    onLogout(e) {
        e.preventDefault();
        $.ajax({
            url: "/api/logout",
            dataType: "json",
            success: (resp) => {
                alert(resp)
            }
        });
        // Wait for 1 second for logout callback to complete
        setTimeout(() => {
            browserHistory.push("/login");
        }, 1000);
    }

    render() {
        var username;
        if (this.props.loggedIn) {
            username = <p>{this.props.username}</p>
        }

        return(
            <aside className="control-sidebar control-sidebar-dark">
                <ul className="control-sidebar-menu">
                    <li>
                        <Link to="/login" activeClassName="active">
                        <i className="menu-icon fa fa-lock bg-green"></i>
                        Login
                        </Link>
                    </li>
                    <li>
                        <a onClick={this.onLogout}>
                        <i className="menu-icon fa fa-lock bg-red"></i>
                        Login
                        </a>
                    </li>
                </ul>
                Current user: {username}
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
                <div className="control-sidebar-bg" style={{position: "fixed", height: "auto"}}></div>
            </aside>
        )
    }
}

HiddenSidebar.propTypes = {
    serverStatus: React.PropTypes.object.isRequired,
    username: React.PropTypes.string.isRequired,
    loggedIn: React.PropTypes.bool.isRequired,
    checkLogin: React.PropTypes.func.isRequired,
}


export default HiddenSidebar
