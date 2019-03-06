import React from 'react';
import PropTypes from 'prop-types';
import {ReactSwalDanger, ReactSwalNormal} from 'Utilities/customSwal';

class UserTable extends React.Component {
    constructor(props) {
        super(props);
        this.removeUser = this.removeUser.bind(this);
    }

    removeUser(user) {
        ReactSwalDanger.fire({
            title: "Are you sure?",
            html: <p>User {user} will be deleted!</p>,
            type: "question",
            showCancelButton: true,
            confirmButtonText: "Yes, delete it!",
            showLoaderOnConfirm: true,
            preConfirm: () => {
                return new Promise((resolve, reject) => {
                    user = {username: user};
                    $.ajax({
                        type: "POST",
                        url: "/api/user/remove",
                        dataType: "json",
                        data: JSON.stringify(user),
                        success: (resp) => {
                            if (resp.success === true) {
                                resolve(resp.data);
                            } else {
                                reject("Unknown error");
                            }
                        },
                        error: () => {
                            reject("Unknown error");
                        }
                    });
                });
            }
        }).then((result) => {
            if (result.value) {
                ReactSwalNormal.fire({
                    title: "Deleted!",
                    text: result.value,
                    type: "success"
                });
                this.props.listUsers();
            }
        }).catch((result) => {
            ReactSwalNormal.fire({
                title: "An error occurred!",
                text: result,
                type: "error"
            });
        });

    }

    render() {
        return(
                <div className="box">
                    <div className="box-header">
                        <h3 className="box-title">Users</h3>
                    </div>
                    
                    <div className="box-body">
                        <table className="table table-striped">
                            <thead>
                                <tr>
                                    <th>User</th>
                                    <th>Role</th>
                                    <th>Email</th>
                                    <th>Delete</th>
                                </tr>
                            </thead>
                            <tbody>
                                {this.props.users.map( (user, i) => {
                                    return(
                                        <tr key={user.username}>
                                            <td>{user.username}</td>
                                            <td>{user.role}</td>
                                            <td>{user.email}</td>
                                            <td>
                                                 <button className="btn btn-danger" onClick={()=>{this.removeUser(user.username)}}>Delete</button>
                                            </td>
                                        </tr>
                                    )                                    
                                })}
                            </tbody>
                        </table>
                    </div>

                </div>
        )
    }
}

UserTable.propTypes = {
    users: PropTypes.array.isRequired,
    listUsers: PropTypes.func.isRequired,
}

export default UserTable
