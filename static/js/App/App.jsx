import React from 'react';
import {RouteHandler} from 'react-router';
import Header from './components/Header.jsx';
import Sidebar from './components/Sidebar.jsx';
import ModsContent from './components/ModsContent.jsx';


class App extends React.Component {
    render() {
        return(
            <div className="wrapper">

            <Header />

            <Sidebar />

            <ModsContent />

            <footer className="main-footer">
                <div className="pull-right hidden-xs">
                Anything you want
                </div>
                <strong>Copyright &copy; 2015 <a href="#">Company</a>.</strong> All rights reserved.
            </footer>

            <aside className="control-sidebar control-sidebar-dark">
                <ul className="nav nav-tabs nav-justified control-sidebar-tabs">
                <li className="active"><a href="#control-sidebar-home-tab" data-toggle="tab"><i className="fa fa-home"></i></a></li>
                <li><a href="#control-sidebar-settings-tab" data-toggle="tab"><i className="fa fa-gears"></i></a></li>
                </ul>
                <div className="tab-content">
                <div className="tab-pane active" id="control-sidebar-home-tab">
                    <h3 className="control-sidebar-heading">Recent Activity</h3>
                    <ul className="control-sidebar-menu">
                    <li>
                        <a href="javascript::;">
                        <i className="menu-icon fa fa-birthday-cake bg-red"></i>

                        <div className="menu-info">
                            <h4 className="control-sidebar-subheading">Langdon's Birthday</h4>

                            <p>Will be 23 on April 24th</p>
                        </div>
                        </a>
                    </li>
                    </ul>

                    <h3 className="control-sidebar-heading">Tasks Progress</h3>
                    <ul className="control-sidebar-menu">
                    <li>
                        <a href="javascript::;">
                        <h4 className="control-sidebar-subheading">
                            Custom Template Design
                            <span className="label label-danger pull-right">70%</span>
                        </h4>

                        <div className="progress progress-xxs">
                            <div className="progress-bar progress-bar-danger" style={{width: "70%"}}></div>
                        </div>
                        </a>
                    </li>
                    </ul>

                </div>
                <div className="tab-pane" id="control-sidebar-stats-tab">Stats Tab Content</div>
                <div className="tab-pane" id="control-sidebar-settings-tab">
                    <form method="post">
                    <h3 className="control-sidebar-heading">General Settings</h3>

                    <div className="form-group">
                        <label className="control-sidebar-subheading">
                        Report panel usage
                        <input type="checkbox" className="pull-right" />
                        </label>

                        <p>
                        Some information about this general settings option
                        </p>
                    </div>
                    </form>
                </div>
                </div>
            </aside>
            <div className="control-sidebar-bg"></div>
        </div>
        )
    }
}

export default App
