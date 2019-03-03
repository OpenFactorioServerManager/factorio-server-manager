import React from 'react';
import {withRouter} from 'react-router-dom';
import FontAwesomeIcon from "./FontAwesomeIcon";

class LoginContent extends React.Component {
    constructor(props) {
        super(props);
        this.loginUser = this.loginUser.bind(this);
    }

    componentDidMount() {}

    loginUser(e) {
        e.preventDefault();
        let user = {
            username: this.refs.username.value,
            password: this.refs.password.value,
        };

        $.ajax({
            type: "POST",
            url: "/api/login",
            dataType: "json",
            data: JSON.stringify(user),
            success: (resp) => {
                console.log(resp);
                this.props.history.push("/");
            }
        });
    }

    render() {
        return(
            <div className="container" id="login">
                <div className="d-flex justify-content-center h-100">
                    <div className="card">
                        <div className="card-header">
                            <h1>
                                <img src="./images/factorio.jpg" className="img-circle" alt="User Image"/>
                                Factorio Server Manager
                            </h1>
                        </div>

                        <div className="car-body">
                            <form onSubmit={this.loginUser}>
                                <label className="input-group form-group">
                                    <div className="input-group-prepend">
                                        <span className="input-group-text">
                                            <FontAwesomeIcon icon="user"/>
                                        </span>
                                    </div>
                                    <input className="form-control" type="text" ref="username" placeholder="Username"/>
                                </label>

                                <label className="input-group form-group">
                                    <div className="input-group-prepend">
                                        <span className="input-group-text">
                                            <FontAwesomeIcon icon="lock"/>
                                        </span>
                                    </div>
                                    <input className="form-control" type="password" ref="password" placeholder="Password"/>
                                </label>

                                <label className="remember-me">
                                    <input type="checkbox"/>&nbsp;
                                    Remember me
                                </label>

                                <input type="submit" value="Sign In" className="btn btn-primary btn-block btn-flat"/>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}

export default withRouter(LoginContent);
