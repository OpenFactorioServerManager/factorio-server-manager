import React from 'react';
import {IndexLink} from 'react-router';

class LoginContent extends React.Component {
    render() {
        return(
            <div className="content-wrapper">
                <section className="content-header">
                <h1>
                    Saves
                    <small>Factorio Save Files</small>
                </h1>
                <ol className="breadcrumb">
                    <li><IndexLink to="/"><i className="fa fa-dashboard fa-fw"></i>Server Control</IndexLink></li>
                    <li className="active">Here</li>
                </ol>
                </section>

                <section className="content">
                <div className="row">
                    <div className="login-box-body">
                        <p className="login-box-msg">Sign in to start your session</p>

                        <form action="">
                            <div className="form-group has-feedback">
                                <input type="email" className="form-control" placeholder="Email" />
                                <span className="fa fa-envelope form-control-feedback"></span>
                            </div>
                            <div className="form-group has-feedback">
                                <input type="password" className="form-control" placeholder="Password" />
                                <span className="fa fa-lock form-control-feedback"></span>
                            </div>
                            <div className="row">
                                <div className="col-xs-8">
                                    <div className="checkbox">
                                        <label className="">
                                            <div className="" aria-checked="false" aria-disabled="false" style="position: relative;">
                                            <input type="checkbox" style="position: absolute; top: -20%; left: -20%; display: block; width: 140%; height: 140%; margin: 0px; padding: 0px; border: 0px; opacity: 0; background: rgb(255, 255, 255);" />
                                            </div> Remember Me
                                        </label>
                                    </div>
                                </div>
                                <div className="col-xs-4">
                                    <button type="submit" className="btn btn-primary btn-block btn-flat">Sign In</button>
                                </div>
                            </div>
                        </form>

                        <a href="#">I forgot my password</a><br />
                    </div>       
                </div>
                </section>
            </div>
        )
    }
}

export default LoginContent
