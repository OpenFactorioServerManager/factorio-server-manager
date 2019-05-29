import React from 'react';
import PropTypes from 'prop-types';

class UploadSave extends React.Component {
    constructor(props) {
        super(props)
        this.updateSavesList = this.updateSavesList.bind(this);
    }

    componentDidMount() {
        $("#savefile").fileinput({
            showPreview: false,
            uploadUrl: '/api/saves/upload',
            showCancel: false,
            allowedFileExtensions: ['zip'],
            theme: "fas",
            removeClass: "btn btn-default",
            uploadClass: "btn btn-default",
        }).on('filebatchuploadsuccess fileuploaded', this.updateSavesList);
    }

    updateSavesList() {
        $('#savefile').fileinput('reset');
        this.props.getSaves();
    }

    render() {
        return(
            <div className="box">
                <div className="box-header">
                    <h4 className="box-title">Upload Save File</h4>
                </div>
                <div className="box-body">
                    <input id="savefile" name="savefile" type="file" className="file"/>
                </div>
            </div>
        )
    }
}

UploadSave.propTypes = {
    getSaves: PropTypes.func.isRequired,
}

export default UploadSave
