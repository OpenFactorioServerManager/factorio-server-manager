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
            var serverStatus = <Link to="/" className="d-block text-success"><FontAwesomeIcon icon="circle"/>Server Online</Link>
        } else {
            var serverStatus = <Link to="/" className="d-block text-danger"><FontAwesomeIcon icon="circle"/>Server Offline</Link>
        }

        return (
            <aside className="main-sidebar sidebar-dark-primary elevation-4">
                <Link className="brand-link" to="/">
                    <span className="logo-lg"><b>Factorio</b>SM</span>
                </Link>

                <div className="sidebar">
                    <div className="user-panel">
                        <div className="image">
                            <img src="./images/factorio.jpg" className="img-circle" alt="User Image"/>
                        </div>
                        <div className="info">
                            <div className="text-white">Factorio Server Manager</div>
                            {serverStatus}
                        </div>
                    </div>

                    {/*<form action="#" method="get" className="sidebar-form">
                        <div className="input-group">
                            <input type="text" name="q" className="form-control" placeholder="Search..."/>
                            <span className="input-group-btn">
                            <button type="submit" name="search" id="search-btn" className="btn btn-flat"><i
                                className="fa fa-search"></i>
                            </button>
                        </span>
                        </div>
                    </form>*/}

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
