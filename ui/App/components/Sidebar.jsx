import React from 'react';
import {Link, NavLink} from 'react-router-dom';
import PropTypes from 'prop-types';
import FontAwesomeIcon from './FontAwesomeIcon.jsx';

class Sidebar extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        if(this.props.serverRunning === "running") {
            var serverStatus = <span className="status text-success"><FontAwesomeIcon icon="circle"/></span>
        } else {
            var serverStatus = <span className="status text-danger"><FontAwesomeIcon icon="circle"/></span>
        }

        return (
            <aside className="main-sidebar sidebar-dark-primary elevation-4">
                <Link className="brand-link logo-switch" to="/">
                    <span className="logo-xl">
                        <img src="./images/factorio.jpg" className="img-circle" alt="User Image"/>&nbsp;
                        <div className="info">
                            <b>Factorio</b>SM&nbsp;
                            {serverStatus}
                        </div>
                    </span>
                    <span className="logo-xs">FSM {serverStatus}</span>
                </Link>

                <div className="sidebar">
                    <nav className="mt-2">
                        <ul className="nav nav-pills nav-sidebar flex-column" data-widget="treeview" role="menu" data-accordion="false">
                            <li className="nav-header">MENU</li>
                            <li className="nav-item">
                                <NavLink exact to="/" activeClassName="active" className="nav-link">
                                    <FontAwesomeIcon icon="tachometer-alt" className="nav-icon"/><p>Server Control</p>
                                </NavLink>
                            </li>
                            <li className="nav-item">
                                <NavLink to="/mods" activeClassName="active" className="nav-link">
                                    <FontAwesomeIcon icon="pencil-alt" className="nav-icon"/><p>Mods</p>
                                </NavLink>
                            </li>
                            <li className="nav-item">
                                <NavLink to="/logs" activeClassName="active" className="nav-link">
                                    <FontAwesomeIcon icon="file-alt" className="nav-icon" prefix="far"/><p>Logs</p>
                                </NavLink>
                            </li>
                            <li className="nav-item">
                                <NavLink to="/saves" activeClassName="active" className="nav-link">
                                    <FontAwesomeIcon icon="save" className="nav-icon" prefix="far"/><p>Saves</p>
                                </NavLink>
                            </li>
                            <li className="nav-item">
                                <NavLink to="/config" activeClassName="active" className="nav-link">
                                    <FontAwesomeIcon icon="cogs" className="nav-icon"/><p>Game Configuration</p>
                                </NavLink>
                            </li>
                            <li className="nav-item">
                                <NavLink to="/settings" activeClassName="active" className="nav-link">
                                    <FontAwesomeIcon icon="cog" className="nav-icon"/><p>Settings</p>
                                </NavLink>
                            </li>
                            <li className="nav-item">
                                <NavLink to="/console" activeClassName="active" className="nav-link">
                                    <FontAwesomeIcon icon="terminal" className="nav-icon"/><p>Console</p>
                                </NavLink>
                            </li>
                        </ul>
                    </nav>
                </div>
            </aside>
        )
    }
}

Sidebar.propTypes = {
    serverStatus: PropTypes.func.isRequired,
    serverRunning: PropTypes.string.isRequired,
}

export default Sidebar
