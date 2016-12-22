import React from 'react';
import swal from 'sweetalert';

class UserTable extends React.Component {
    constructor(props) {
        super(props);
        this.removeUser = this.removeUser.bind(this);
    }

    removeUser(user) {
        swal({   
            title: "Are you sure?",  
            text: "User: " + user + " will be deleted",   
            type: "warning",   
            showCancelButton: true,   
            confirmButtonColor: "#DD6B55",   
            confirmButtonText: "Yes, delete it!",   
            closeOnConfirm: false 
        }, 
        () => {
            user = {username: user}
            $.ajax({
                type: "POST",
                url: "/api/user/remove",
                dataType: "json",
                data: JSON.stringify(user),
                success: (resp) => {
                    if (resp.success === true) {
                        swal("Deleted!", resp.data, "success"); 
                    }
                }
            })
        });

        this.props.listUsers();

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
    users: React.PropTypes.array.isRequired,
    listUsers: React.PropTypes.func.isRequired,
}

export default UserTable
