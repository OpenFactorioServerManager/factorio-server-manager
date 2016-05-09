import React from 'react';
import {IndexLink} from 'react-router';
import UserTable from './Users/UserTable.jsx';
import AddUser from './Users/AddUser.jsx';

class UsersContent extends React.Component {
    constructor(props) {
        super(props);
        this.listUsers = this.listUsers.bind(this);
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
                    this.setState({users: resp.data})
                } else {
                    console.log("error listing users.")
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
                    <UserTable
                        users={this.state.users}
                        listUsers={this.listUsers}
                    /> 
                    <AddUser 
                        listUsers={this.listUsers}
                    />

                </section>
            </div>
        )
    }
}

export default UsersContent;
