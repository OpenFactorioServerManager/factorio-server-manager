import React from 'react';

class ModUpload extends React.Component {
    componentDidMount() {
        $("#mod_upload_input").fileinput({
            uploadUrl: '/api/mods/upload',
            // hiddenThumbnailContent: true,
            // showPreview: false,
            showCancel: false,
            showUploadedThumbs: false,
            browseOnZoneClick: true,
            uploadAsync: false,
            allowedFileExtensions: ['zip'],
            browseLabel: "Select Mods...",
            browseIcon: '<i class="fa fa-upload text-muted" style="color: white;"></i>&nbsp;',
        }).on('filebatchuploadsuccess fileuploaded', this.props.uploadModSuccessHandler);
    }

    render() {
        return(
            <div className="box-body">
                <div className="alert alert-warning alert-dismissible" role="alert">
                    <button type="button" className="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    The mods you upload will override every mod, that is already uploaded!<br/>
                    Uploaded mods will be treated same, as if you install mods with the mod-portal-api.<br/><br/>
                    Only zip-files are allowed to upload !!
                    {/*TODO make better english*/}
                </div>

                <label className="control-label">Select File</label>
                <input id="mod_upload_input" name="mod_file" multiple type="file" />
            </div>
        );
    }
}

ModUpload.PropTypes = {
    uploadModSuccessHandler: React.PropTypes.func.isRequired
};

export default ModUpload;