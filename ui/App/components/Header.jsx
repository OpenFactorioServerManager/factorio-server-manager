import React from 'react';
import {Link, withRouter} from 'react-router-dom';
import PropTypes from 'prop-types';
import FontAwesomeIcon from "./FontAwesomeIcon";

class Header extends React.Component {
    constructor(props) {
        super(props);
        this.onLogout = this.onLogout.bind(this);
    }
    
    onLogout(e) {
        e.preventDefault();
        $.ajax({
            url: "/api/logout",
            dataType: "json",
            success: (resp) => {
                console.log(resp)
            }
        });

        // Wait for 1 second for logout callback to complete
        setTimeout(() => {
            this.props.history.push("/login")
        }, 1000);
    }

    render() {
        var loginMenu; 
        if (this.props.loggedIn) {
            loginMenu = 
                <ul className="navbar-nav ml-auto">
                    <li className="nav-item">
                        <Link className="nav-link" to="/settings">
                            <FontAwesomeIcon icon="cogs" className="fa-fw"/>Settings
                        </Link>
                    </li>
                    <li className="nav-item">
                        <a href="javascript:void(0)" onClick={this.onLogout} className="nav-link">
                            <FontAwesomeIcon icon="lock" className="fa-fw"/>Logout
                        </a>
                    </li>
                </ul>
        }
        return(
            <nav className="main-header navbar navbar-expand navbar-light border-bottom">
                <ul className="navbar-nav">
                    <li className="nav-item">
                        <a className="nav-link" data-widget="pushmenu" href="#">
                            <FontAwesomeIcon icon="bars"/>
                        </a>
                    </li>
                </ul>

                {loginMenu}
            </nav>
        )
    }
}

Header.propTypes = {
    username: PropTypes.string.isRequired,
    loggedIn: PropTypes.bool.isRequired,
}

export default withRouter(Header);
