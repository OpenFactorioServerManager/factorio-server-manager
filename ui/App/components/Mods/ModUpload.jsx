import React from 'react';
import PropTypes from 'prop-types';

class ModUpload extends React.Component {
    componentDidMount() {
        $("#mod_upload_input").fileinput({
            uploadUrl: '/api/mods/upload',
            showCancel: false,
            showUploadedThumbs: false,
            browseOnZoneClick: true,
            uploadAsync: false,
            allowedFileExtensions: ['zip'],
            browseLabel: "Select Mods...",
            browseIcon: '<i class="fa fa-upload text-muted" style="color: white;"></i>&nbsp;',
            theme: "fas",
            slugCallback: function(filename) {
                return filename;
            },
        }).on('filebatchuploadsuccess fileuploaded', this.props.uploadModSuccessHandler);
    }

    render() {
        let classes = "card-body" + " " + this.props.className;
        let ids = this.props.id;

        return(
            <div id={ids} className={classes}>
                <div className="alert alert-warning alert-dismissible" role="alert">
                    <button type="button" className="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    The mods you upload will override every mod, that is already uploaded!<br/>
                    Uploaded mods will be treated same, as if you install mods with the mod-portal-api.<br/><br/>
                    Only zip-files are allowed to upload !!
                </div>

                <label className="control-label">Select File</label>
                <input id="mod_upload_input" name="mod_file" multiple type="file" />
            </div>
        );
    }
}

ModUpload.propTypes = {
    uploadModSuccessHandler: PropTypes.func.isRequired,
    className: PropTypes.string,
    id: PropTypes.string
};

export default ModUpload;