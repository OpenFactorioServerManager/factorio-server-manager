import React from 'react';
import {Link, browserHistory} from 'react-router';

class HiddenSidebar extends React.Component {
    constructor(props) {
        super(props);
    }

    capitalizeFirstLetter(string) {
        return string.charAt(0).toUpperCase() + string.slice(1);
    }

    render() {
        var username;
        if (this.props.loggedIn) {
            username = <p>{this.props.username}</p>
        }

        return(
        <aside id="control-sidebar" className="control-sidebar control-sidebar-dark">
            <ul className="nav nav-tabs nav-justified control-sidebar-tabs">
            <li className="active"><a href="control-sidebar-theme-demo-options-tab" data-toggle="tab"><i className="fa fa-wrench"></i></a></li>
            <li><a href="#control-sidebar-home-tab" data-toggle="tab"><i className="fa fa-home"></i></a></li>
            <li><a href="#control-sidebar-settings-tab" data-toggle="tab"><i className="fa fa-gears"></i></a></li>
            </ul>
            <div className="tab-content">
                <div className="tab-pane" id="control-sidebar-home-tab">
                    <h3 className="control-sidebar-heading">Recent Activity</h3>
                    <ul className="control-sidebar-menu">
                        <li>
                            <Link to="/login" activeClassName="active">
                            <i className="menu-icon fa fa-birthday-cake bg-red"></i>

                            <div className="menu-info">
                                <i classNameName="menu-icon fa fa-lock bg-green"></i>
                                <h4 className="control-sidebar-subheading">Login</h4>
                            </div>
                            </Link>
                        </li>
                        <li>
                            <a href="javascript:void(0)">
                            <i className="menu-icon fa fa-lock bg-yellow"></i>

                            <div className="menu-info">
                                
                                <h4 className="control-sidebar-subheading">Logout</h4>
                            </div>
                            </a>
                        </li>
                    </ul>

                    <h3 className="control-sidebar-heading">Tasks Progress</h3>
                        <div classNameName="table-responsive">
                        <table classNameName="table table-border">
                            <thead>
                                <tr>
                                </tr>
                            </thead>
                            <tbody>
                                {Object.keys(this.props.serverStatus).map(function(key) {
                                    return(
                                        <tr key={key}>
                                            <td>{this.capitalizeFirstLetter(key)}</td>
                                            <td>{this.props.serverStatus[key]}</td>
                                        </tr>
                                    )                                                  
                                }, this)}        
                            </tbody>
                        </table>
                        </div>
                    <ul className="control-sidebar-menu">
                    <li>
                        <a href="javascript:void(0)">
                        <h4 className="control-sidebar-subheading">
                            Custom Template Design
                            <span className="label label-danger pull-right">70%</span>
                        </h4>

                        <div className="progress progress-xxs">
                            <div className="progress-bar progress-bar-danger"></div>
                        </div>
                        </a>
                    </li>
                    <li>
                        <a href="javascript:void(0)">
                        <h4 className="control-sidebar-subheading">
                            Update Resume
                            <span className="label label-success pull-right">95%</span>
                        </h4>

                        <div className="progress progress-xxs">
                            <div className="progress-bar progress-bar-success"></div>
                        </div>
                        </a>
                    </li>
                    <li>
                        <a href="javascript:void(0)">
                        <h4 className="control-sidebar-subheading">
                            Laravel Integration
                            <span className="label label-warning pull-right">50%</span>
                        </h4>

                        <div className="progress progress-xxs">
                            <div className="progress-bar progress-bar-warning"></div>
                        </div>
                        </a>
                    </li>
                    <li>
                        <a href="javascript:void(0)">
                        <h4 className="control-sidebar-subheading">
                            Back End Framework
                            <span className="label label-primary pull-right">68%</span>
                        </h4>

                        <div className="progress progress-xxs">
                            <div className="progress-bar progress-bar-primary"></div>
                        </div>
                        </a>
                    </li>
                    </ul>
                </div>

                <div className="tab-pane" id="control-sidebar-settings-tab">
                    <form method="post">
                    <h3 className="control-sidebar-heading">General Settings</h3>

                    <div className="form-group">
                        <label className="control-sidebar-subheading">
                        Report panel usage
                        <input type="checkbox" className="pull-right" checked="" />
                        </label>

                        <p>
                        Some information about this general settings option
                        </p>
                    </div>

                    <div className="form-group">
                        <label className="control-sidebar-subheading">
                        Allow mail redirect
                        <input type="checkbox" className="pull-right" checked="" />
                        </label>

                        <p>
                        Other sets of options are available
                        </p>
                    </div>

                    <div className="form-group">
                        <label className="control-sidebar-subheading">
                        Expose author name in posts
                        <input type="checkbox" className="pull-right" checked="" />
                        </label>

                        <p>
                        Allow the user to show his name in blog posts
                        </p>
                    </div>

                    <h3 className="control-sidebar-heading">Chat Settings</h3>

                    <div className="form-group">
                        <label className="control-sidebar-subheading">
                        Show me as online
                        <input type="checkbox" className="pull-right" checked="" />
                        </label>
                    </div>

                    <div className="form-group">
                        <label className="control-sidebar-subheading">
                        Turn off notifications
                        <input type="checkbox" className="pull-right" />
                        </label>
                    </div>

                    <div className="form-group">
                        <label className="control-sidebar-subheading">
                        Delete chat history
                        <a href="javascript:void(0)" className="text-red pull-right"><i className="fa fa-trash-o"></i></a>
                        </label>
                    </div>
                    </form>
                </div>
            </div>
        </aside>
        )
    }
}

HiddenSidebar.propTypes = {
    serverStatus: React.PropTypes.object.isRequired,
    username: React.PropTypes.string.isRequired,
    loggedIn: React.PropTypes.bool.isRequired,
    checkLogin: React.PropTypes.func.isRequired,
}


export default HiddenSidebar
