import React from 'react';
import {Link} from 'react-router-dom';
import UserTable from './Users/UserTable.jsx';
import AddUser from './Users/AddUser.jsx';
import FontAwesomeIcon from "./FontAwesomeIcon";

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

                        <small className="float-sm-right">
                            <ol className="breadcrumb">
                                <li className="breadcrumb-item">
                                    <Link to="/"><FontAwesomeIcon icon="tachometer-alt"/>Server Control</Link>
                                </li>
                                <li className="breadcrumb-item active">
                                    <FontAwesomeIcon icon="cog"/>Settings
                                </li>
                            </ol>
                        </small>
                    </h1>

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
