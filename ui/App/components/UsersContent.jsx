import React from 'react';
import {IndexLink} from 'react-router';

class UsersContent extends React.Component {
    constructor(props) {
        super(props);
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
        let user = {
            username: this.refs.username,
            role: this.refs.role,
            password: this.refs.password,
            email: this.refs.email,
        }
        $.ajax({
            type: "POST",
            url: "/api/user/add",
            dataType: "json",
            data: JSON.stringify(user),
            success: (resp) => {
                if (resp.success === true) {
                    alert("User: ", user.username, " added successfully.")
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
                                <h3>{user.Username}</h3>        
                                <h3>{user.Role}</h3>        
                                <h3>{user.Email}</h3>        
                                </div>
                            )                                    
                        })}

                        </div>
                    </div>
                </section>
            </div>
        )
    }
}

export default UsersContent;
