import React from 'react';
import {Link, withRouter} from 'react-router-dom';
import PropTypes from 'prop-types';

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
                <ul className="nav navbar-nav">
                    <li>
                        <Link to="/settings"><i className="fa fa-gears fa-fw"></i>Settings</Link>
                    </li>
                    <li>
                        <a href="javascript:void(0)" onClick={this.onLogout}><i className="fa fa-lock fa-fw"></i>Logout</a>
                    </li>
                </ul>
        }
        return(
            <header className="main-header">
                
                <Link className="logo" to="/"><span className="logo-lg"><b>Factorio</b>SM</span></Link>
                
                <nav className="navbar navbar-static-top" role="navigation">
                <a href="#" className="sidebar-toggle" data-toggle="offcanvas" role="button">
                <span className="sr-only">Toggle navigation</span>
                </a>
                <div className="navbar-custom-menu">
                    {loginMenu}
                </div>
                </nav>
            </header>
        )
    }
}

Header.propTypes = {
    username: PropTypes.string.isRequired,
    loggedIn: PropTypes.bool.isRequired,
}

export default withRouter(Header);
