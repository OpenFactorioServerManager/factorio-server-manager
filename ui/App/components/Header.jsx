import React from 'react';
import {IndexLink} from 'react-router';

class Header extends React.Component {
    render() {
        return(
            <header className="main-header">
                
                <IndexLink className="logo" to="/"><span className="logo-lg"><b>Factorio</b>SM</span></IndexLink>
                
                <nav className="navbar navbar-static-top" role="navigation">
                <a href="#" className="sidebar-toggle" data-toggle="offcanvas" role="button">
                <span className="sr-only">Toggle navigation</span>
                </a>
                <div className="navbar-custom-menu">
                    <ul className="nav navbar-nav">

                    <li className="dropdown notifications-menu">
                        <a href="#" className="dropdown-toggle" data-toggle="dropdown">
                        <i className="fa fa-bell-o"></i>
                        <span className="label label-warning">10</span>
                        </a>
                        <ul className="dropdown-menu">
                        <li className="header">You have 10 notifications</li>
                        <li>
                            <ul className="menu">
                                <a href="#">
                                <i className="fa fa-users text-aqua"></i> 5 new members joined today
                                </a>
                            </ul>
                        </li>
                        <li className="footer"><a href="#">View all</a></li>
                        </ul>
                    </li>
                    <li>
                        <a href="#" data-toggle="control-sidebar"><i className="fa fa-gears"></i></a>
                    </li>
                    </ul>
                </div>
                </nav>
            </header>
        )
    }
}

export default Header
