import React from 'react';

class Header extends React.Component {
    render() {
        return(
            <header className="main-header">

                <a href="/" className="logo">
                <span className="logo-mini"><b>F</b>SM</span>
                <span className="logo-lg"><b>Factorio</b>SM</span>
                </a>

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
                    <li className="dropdown user user-menu">
                        <a href="#" className="dropdown-toggle" data-toggle="dropdown">
                        <img src="./dist/dist/img/user2-160x160.jpg" className="user-image" alt="User Image" />
                        <span className="hidden-xs">Alexander Pierce</span>
                        </a>
                        <ul className="dropdown-menu">
                        <li className="user-header">
                            <img src="./dist/dist/img/user2-160x160.jpg" className="img-circle" alt="User Image" />

                            <p>
                            Alexander Pierce - Web Developer
                            <small>Member since Nov. 2012</small>
                            </p>
                        </li>
                        <li className="user-body">
                            <div className="row">
                            <div className="col-xs-4 text-center">
                                <a href="#">Followers</a>
                            </div>
                            <div className="col-xs-4 text-center">
                                <a href="#">Sales</a>
                            </div>
                            <div className="col-xs-4 text-center">
                                <a href="#">Friends</a>
                            </div>
                            </div>
                        </li>
                        <li className="user-footer">
                            <div className="pull-left">
                            <a href="#" className="btn btn-default btn-flat">Profile</a>
                            </div>
                            <div className="pull-right">
                            <a href="#" className="btn btn-default btn-flat">Sign out</a>
                            </div>
                        </li>
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
