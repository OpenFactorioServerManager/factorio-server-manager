import React from 'react';
import {IndexLink} from 'react-router';

class UsersContent extends React.Component {
    constructor(props) {
        super(props);
        this.listUsers = this.listUsers.bind(this);
        this.createUser = this.createUser.bind(this);
        this.state = {
            users: [],
        }
    }

    componentDidMount() {
        this.listUsers();
    }

    listUsers() {
        $.ajax({
            type: "GET",
            url: "/api/user/list",
            dataType: "json",
            success: (resp) => {
                if (resp.success === true) {
                    console.log("Listing users: ", resp.data)
                    this.setState({users: resp.data})
                    console.log(this.state)
                } else {
                    console.log("error listing users.")
                }
            }
        })
    }

    createUser(e) {
        e.preventDefault();
        console.log(this.refs);
        let user = {
            username: this.refs.username.value,
            // Add handler for listing roles
            role: "admin",
            password: this.refs.password.value,
            email: this.refs.email.value,
        }
        if (user.password !== this.refs.passwordConfirm.value) {
            console.log("passwords do not match")
            return
        }

        $.ajax({
            type: "POST",
            url: "/api/user/add",
            dataType: "json",
            data: JSON.stringify(user),
            success: (resp) => {
                if (resp.success === true) {
                    alert("User: " + user.username + " added successfully.");
                    this.listUsers();
                } else {
                    alert("Error deleting user: ", resp.data)
                }
            }
        })
    }

    render() {
        return(
            <div className="content-wrapper">
                <section className="content-header">
                <h1>
                    Settings
                    <small>Manage Factorio Server Manager settings</small>
                </h1>
                <ol className="breadcrumb">
                    <li><IndexLink to="/"><i className="fa fa-dashboard fa-fw"></i>Server Control</IndexLink></li>
                    <li className="active">Here</li>
                </ol>
                </section>

                <section className="content">
                    <div className="box">
                        <div className="box-header">
                            <h3 className="box-title">Users</h3>
                        </div>
                        
                        <div className="box-body">
                        
                        {this.state.users.map( (user, i) => {
                            return(
                                <div>
                                <h4>{user.Username}</h4>        


                                </div>
                            )                                    
                        })}

                        <h4>Add user</h4>

                        <form action="" onSubmit={this.createUser}>
                            <div className="form-group">
                                <label for="username">Username</label>
                                <input ref="username" type="text" className="form-control" id="username" placeholder="Enter username" />
                            </div> 
                            <div className="form-group">
                                <label for="password">Password</label>
                                <input ref="password" type="text" className="form-control" id="password" placeholder="Enter password" />
                            </div> 
                            <div className="form-group">
                                <label for="password">Password confirmation</label>
                                <input ref="passwordConfirm" type="text" className="form-control" id="password" placeholder="Enter password again" />
                            </div> 
                            <div className="form-group">
                                <label for="email">Email address</label>
                                <input ref="email" type="text" className="form-control" id="email" placeholder="Enter email again" />
                            </div> 
                        
                            <button className="btn btn-block btn-success" type="submit"><i className="fa fa-plus fa-fw"></i>Add User</button>
                        </form>
                        </div>
                    </div>
                </section>
            </div>
        )
    }
}

export default UsersContent;
