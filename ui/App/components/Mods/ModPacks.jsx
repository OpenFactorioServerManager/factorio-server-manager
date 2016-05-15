import React from 'react';

class ModPacks extends React.Component {
    constructor(props) {
        super(props)
        this.createModPack = this.createModPack.bind(this);
        this.removeModPack = this.removeModPack.bind(this);
    }

    createModPack(e, modpack) {
        e.preventDefault();
        self = this;

        let modpacks = {};
        modpacks["mods"] = [];
        modpacks["title"] = this.refs.modpackName.value;

        for (var m in this.refs) {
            if (this.refs[m].checked) {
                modpacks["mods"].push(this.refs[m].id)
            }
        }
        $.ajax({
            url: "/api/mods/packs/add",
            dataType: "json",
            type: "POST",
            data: JSON.stringify(modpacks),
            success: (resp) => {
                if (resp.success === true) {
                    swal("Added modpack", "Modpack: " + modpacks["title"] + " added successfully", "success")
                    self.props.loadModPackList();
                } else {
                    swal("Error", "Could not create modpack " + modpacks["title"], "error")
                }
            }
        })
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
                    <div className="box box-success collapsed-box">
                        <button type="button" className="btn btn-box-tool" data-widget="collapse">
                            <div className="box-header with-border">
                            <i className="fa fa-plus fa-fw"></i><h4 className="box-title">Create Mod Pack</h4>
                            </div>
                        </button>
                        <div className="box-body" style={{display: "none"}}>
                            <form onSubmit={this.createModPack}>
                                <label for="modPack">Mod Pack Name</label>
                                <div id="modpack" className="input-group">
                                    <input ref="modpackName" name="modpack" id="modpack" type="text" className="form-control" placeholder="Enter Mod Pack Name..." />
                                </div>
                                {this.props.installedMods.map( (mod, i) => {
                                    return(
                                    <div className="checkbox" key={i}>
                                        <label for={mod}>
                                            <input id={mod} ref={"mod-"+mod} type="checkbox" />
                                            {mod}
                                        </label>
                                    </div>
                                    )
                                } )}
                                <div className="col-md-4">
                                    <button className="btn btn-block btn-success" type="submit"><i className="fa fa-save fa-fw"></i>Create Mod Pack</button>
                                </div>
                            </form>
                        </div>
                    </div>
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
    installedMods: React.PropTypes.array.isRequired,
    modPacks: React.PropTypes.array.isRequired, 
    loadModPackList: React.PropTypes.func.isRequired,
}

export default ModPacks
