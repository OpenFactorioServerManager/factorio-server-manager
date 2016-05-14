import React from 'react';

class ModPacks extends React.Component {
    constructor(props) {
        super(props)
    }

    removeModPack(modpack, e) {
        var self = this;
        swal({   
            title: "Are you sure?",  
            text: "Modpack: " + modpack + " will be deleted",   
            type: "warning",   
            showCancelButton: true,   
            confirmButtonColor: "#DD6B55",   
            confirmButtonText: "Yes, delete it!",   
            closeOnConfirm: false 
        }, 
        () => {
            $.ajax({
                url: "/api/mods/packs/rm/" + modpack,
                dataType: "json",
                success: (resp) => {
                    if (resp.success === true) {
                        swal("Deleted!", resp.data, "success"); 
                        self.props.loadModPackList();
                    }
                }
            })
        });

    }

    render() {
        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Mod Packs</h3>
                </div>
                     
                <div className="box-body">
                    <table className="table table-striped">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Download</th>
                                <th>Delete</th>
                            </tr>
                        </thead>
                        <tbody>
                        {this.props.modPacks.map( (mod, i) => {
                            console.log(mod)
                            let dlURL = "/api/mods/packs/dl/" + mod
                            return(
                                <tr key={i}>
                                    <td>{mod}</td>    
                                    <td><a className="btn btn-default" href={dlURL}></a></td>
                                    <td>
                                        <button
                                            className="btn btn-danger"
                                            ref="modpack"
                                            type="button"
                                            onClick={this.removeModPack.bind(this, mod)}>
                                        </button>
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

ModPacks.propTypes = {
    modPacks: React.PropTypes.array.isRequired, 
    loadModPackList: React.PropTypes.func.isRequired,
}

export default ModPacks
