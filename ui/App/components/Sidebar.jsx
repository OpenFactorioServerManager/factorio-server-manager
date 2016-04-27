import React from 'react';
import {Link, IndexLink} from 'react-router';

class Sidebar extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        if (this.props.serverRunning === "running") {
            var serverStatus = 
                <IndexLink to="/"><i className="fa fa-circle text-success"></i>Server Online</IndexLink>
        } else {
            var serverStatus = 
                <IndexLink to="/"><i className="fa fa-circle text-danger"></i>Server Offline</IndexLink>
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
                    <li><IndexLink to="/" activeClassName="active"><i className="fa fa-link"></i><span>Server Control</span></IndexLink></li>
                    <li><Link to="/mods" activeClassName="active"><i className="fa fa-link"></i><span>Mods</span></Link></li>
                    <li><Link to="/logs" activeClassName="active"><i className="fa fa-link"></i> <span>Logs</span></Link></li>
                    <li><Link to="/saves" activeClassName="active"><i className="fa fa-link"></i> <span>Saves</span></Link></li>
                    <li><Link to="/config" activeClassName="active"><i className="fa fa-link"></i> <span>Configuration</span></Link></li>
                </ul>
                </section>
                <div style={{height: "100%"}}></div>
            </aside>
        )
    }
}

Sidebar.propTypes = {
    serverStatus: React.PropTypes.func.isRequired,
    serverRunning: React.PropTypes.string.isRequired,
}

export default Sidebar
