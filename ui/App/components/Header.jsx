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
                <ul className="navbar-nav ml-auto">
                    <li className="nav-item">
                        <Link className="nav-link" to="/settings">
                            <i className="fa fa-gears fa-fw"></i>Settings
                        </Link>
                    </li>
                    <li className="nav-item">
                        <a href="javascript:void(0)" onClick={this.onLogout} className="nav-link">
                            <i className="fa fa-lock fa-fw"></i>Logout
                        </a>
                    </li>
                </ul>
        }
        return(
            <nav className="main-header navbar navbar-expand bg-white navbar-light border-bottom">
                <ul className="navbar-nav">
                    <li className="nav-item">
                        <a className="nav-link" data-widget="pushmenu" href="#">
                            <i className="fa fa-bars"></i>
                        </a>
                    </li>
                </ul>

                {loginMenu}
            </nav>

            // <header className="main-header">
            //     <Link className="logo" to="/"><span className="logo-lg"><b>Factorio</b>SM</span></Link>
            //
            //     <nav className="navbar navbar-static-top" role="navigation">
            //         <a href="#" className="sidebar-toggle" data-toggle="offcanvas" role="button">
            //             <span className="sr-only">Toggle navigation</span>
            //         </a>
            //         <div className="navbar-custom-menu">
            //             {loginMenu}
            //         </div>
            //     </nav>
            // </header>
        )
    }
}

Header.propTypes = {
    username: PropTypes.string.isRequired,
    loggedIn: PropTypes.bool.isRequired,
}

export default withRouter(Header);
