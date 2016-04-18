import React from 'react';

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
                    <a href="#"><i className="fa fa-circle text-success"></i> Online</a>
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
                    <li className="active"><a href="#"><i className="fa fa-link"></i> <span>Saves</span></a></li>
                    <li><a href="#"><i className="fa fa-link"></i> <span>Logs</span></a></li>
                    <li className="treeview">
                    <a href="#"><i className="fa fa-link"></i> <span>Mods</span> <i className="fa fa-angle-left pull-right"></i></a>
                    <ul className="treeview-menu">
                        <li><a href="#">Server Mods</a></li>
                        <li><a href="#">Get Mods</a></li>
                    </ul>
                    </li>
                </ul>
                </section>
            </aside>
        )
    }
}

export default Sidebar
