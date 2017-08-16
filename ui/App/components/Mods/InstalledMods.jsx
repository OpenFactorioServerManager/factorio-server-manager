import React from 'react';
import swal from 'sweetalert';

class InstalledMods extends React.Component {
    componentDidMount() {
        this.uploadFile = this.uploadFile.bind(this);
        this.removeMod = this.removeMod.bind(this);
    }

    updateInstalledMods() {
        this.props.loadInstalledModList();
    }

    uploadFile(e) {
        e.preventDefault();
        var fd = new FormData();
        fd.append('modfile', this.refs.file.files[0]);

        $.ajax({
            url: "/api/mods/upload",
            type: "POST",
            data: fd,
            processData: false,
            contentType: false,
            success: (data) => {
                var response = JSON.parse(data)
                if (response.success === true) {
                    this.updateInstalledMods();
                }
            }
        });
    }

    removeMod(i) {
        var self = this;
        swal({   
            title: "Are you sure?",  
            text: "Save: " + self.props.installedMods[i] + " will be deleted",   
            type: "warning",   
            showCancelButton: true,   
            confirmButtonColor: "#DD6B55",   
            confirmButtonText: "Yes, delete it!",   
            closeOnConfirm: false 
        }, 
        () => {
            $.ajax({
                url: "/api/mods/rm/" + self.props.installedMods[i],
                success: (resp) => {
                    if (resp.success === true) {
                        swal("Deleted!", resp.data, "success"); 
                        self.updateInstalledMods();
                    }
                }
            })
        });
    }

    render() {
        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Installed Mods</h3>
                </div>
                     
                <div className="box-body">
                    <h4>Upload Mod</h4>
                    <form ref="uploadForm" className="form-inline" encType='multipart/form-data'>
                        <div className="form-group">
                            <label for="modfile">Upload Mod File...</label>
                            <input className="form-control btn btn-default" ref="file" type="file" name="modfile" id="modfile" />
                        </div>
                        <div className="form-group">
                            <input className="form-control btn btn-default" type="button" ref="button" value="Upload" onClick={this.uploadFile} />
                        </div>
                    </form>
                    
                    <div className="table-responsive">
                    <table className="table table-striped">
                        <thead>
                            <tr>
                                <th>Mod Name</th>
                                <th>Download</th>
                                <th>Delete</th>
                            </tr>
                        </thead>
                        <tbody>
                        {this.props.installedMods.map ( (mod, i) => {
                            var saveLocation = "/api/mods/dl/" + mod;
                            return(
                                <tr key={i}>
                                    <td>
                                        {mod}
                                    </td>
                                    <td>
                                        <a className="btn btn-default" href={saveLocation}>Download</a>
                                    </td>
                                    <td>
                                        <button
                                            className="btn btn-danger btn-small" 
                                            ref="modInput"
                                            type="button" 
                                            onClick={this.removeMod.bind(this, i)}>
                                        <i className="fa fa-trash"></i>
                                        &nbsp;
                                        Delete
                                        </button>
                                    </td>
                                </tr>
                            )                                            
                        })}
                        </tbody>
                    </table>
                    </div>
                </div>
            </div>
        )
    }
}

InstalledMods.propTypes = {
    installedMods: React.PropTypes.array.isRequired,
    loadInstalledModList: React.PropTypes.func.isRequired
}

export default InstalledMods
