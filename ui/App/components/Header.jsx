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
