import React from 'react';
import {Link, NavLink} from 'react-router-dom';
import PropTypes from 'prop-types';

class Sidebar extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        if (this.props.serverRunning === "running") {
            var serverStatus = 
                <Link to="/"><i className="fa fa-circle text-success"></i>Server Online</Link>
        } else {
            var serverStatus = 
                <Link to="/"><i className="fa fa-circle text-danger"></i>Server Offline</Link>
        }

        return(
            <aside className="main-sidebar">
                <section className="sidebar" style={{height: "100%"}}>

                <div className="user-panel">
                    <div className="pull-left image">
                        <img src="./dist/dist/img/factorio.jpg" className="img-circle" alt="User Image" />
                    </div>
                    <div className="pull-left info">
                        <p>Factorio Server Manager</p>
                        {serverStatus}
                    </div>
                </div>

                <form action="#" method="get" className="sidebar-form">
                    <div className="input-group">
                    <input type="text" name="q" className="form-control" placeholder="Search..." />
                        <span className="input-group-btn">
                            <button type="submit" name="search" id="search-btn" className="btn btn-flat"><i className="fa fa-search"></i>
                            </button>
                        </span>
                    </div>
                </form>

                <ul className="sidebar-menu">
                    <li className="header">MENU</li>
                    <li>
                        <NavLink exact to="/" activeClassName="active">
                            <i className="fa fa-tachometer"></i><span>Server Control</span>
                        </NavLink>
                    </li>
                    <li>
                        <NavLink to="/mods" activeClassName="active">
                            <i className="fa fa-pencil"></i><span>Mods</span>
                        </NavLink>
                    </li>
                    <li>
                        <NavLink to="/logs" activeClassName="active">
                            <i className="fa fa-file-text-o"></i><span>Logs</span>
                        </NavLink>
                    </li>
                    <li>
                        <NavLink to="/saves" activeClassName="active">
                            <i className="fa fa-floppy-o"></i><span>Saves</span>
                        </NavLink>
                    </li>
                    <li>
                        <NavLink to="/config" activeClassName="active">
                            <i className="fa fa-cogs"></i><span>Game Configuration</span>
                        </NavLink>
                    </li>
                    <li>
                        <NavLink to="/settings" activeClassName="active">
                            <i className="fa fa-cog"></i><span>Settings</span>
                        </NavLink>
                    </li>
                    <li>
                        <NavLink to="/console" activeClassName="active">
                            <i className="fa fa-terminal"></i><span>Console</span>
                        </NavLink>
                    </li>
                </ul>
                </section>
                <div style={{height: "100%"}}></div>
            </aside>
        )
    }
}

Sidebar.propTypes = {
    serverStatus: PropTypes.func.isRequired,
    serverRunning: PropTypes.string.isRequired,
}

export default Sidebar
