import React from 'react';
import {Link} from 'react-router';

class Sidebar extends React.Component {
    render() {
        return(
            <aside className="main-sidebar">

                <section className="sidebar">

                <div className="user-panel">
                    <div className="pull-left image">
                    <img src="./dist/dist/img/factorio.jpg" className="img-circle" alt="User Image" />
                    </div>
                    <div className="pull-left info">
                    <p>Factorio Server Manager</p>
                    <a href="#"><i className="fa fa-circle text-success"></i>Server Online</a>
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
                    <li><Link to="/mods" activeClassName="active"><i className="fa fa-link"></i><span>Mods</span></Link></li>
                    <li><Link to="/logs" activeClassName="active"><i className="fa fa-link"></i> <span>Logs</span></Link></li>
                    <li><Link to="/saves" activeClassName="active"><i className="fa fa-link"></i> <span>Saves</span></Link></li>
                </ul>
                </section>
            </aside>
        )
    }
}

export default Sidebar
