import React from 'react';

class UploadSave extends React.Component {
    constructor(props) {
        super(props)
        this.updateSavesList = this.updateSavesList.bind(this);
        this.uploadFile = this.uploadFile.bind(this);
    }

    updateSavesList() {
        this.props.getSaves();
    }

    uploadFile(e) {
        var fd = new FormData();
        fd.append('savefile', this.refs.file.files[0]);

        $.ajax({
            url: "/api/saves/upload",
            type: "POST",
            data: fd,
            processData: false,
            contentType: false,
            success: (data) => {
                alert(data)
            }
        });
        e.preventDefault();
        this.updateSavesList();
    }

    render() {
        return(
            <div className="box">
                <div className="box-header">
                    <h4 className="box-title">Upload Save File</h4>
                </div>
                <div className="box-body">
                    <form ref="uploadForm" className="form-inline" encType='multipart/form-data'>
                        <div className="form-group">
                            <label for="savefile">Upload Save File...</label>
                            <input className="form-control btn btn-default" ref="file" type="file" name="savefile" id="savefile" />
                        </div>
                        <div className="form-group">
                            <input className="form-control btn btn-default" type="button" ref="button" value="Upload" onClick={this.uploadFile} />
                        </div>
                    </form>
                </div>
            </div>
        )
    }
}

UploadSave.propTypes = {
    getSaves: React.PropTypes.func.isRequired,
}

export default UploadSave
