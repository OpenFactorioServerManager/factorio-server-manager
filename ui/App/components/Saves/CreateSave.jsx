import React from 'react';

class CreateSave extends React.Component {
    constructor(props) {
        super(props);
        this.createSaveFile = this.createSaveFile.bind(this);
        this.updateSavesList = this.updateSavesList.bind(this)
        this.state = {
            loading: false,
        }
    }

    updateSavesList() {
        this.props.getSaves();
    }

    createSaveFile(e) {
        this.setState({loading: true});
        $.ajax({
            url: "/api/saves/create/" + this.refs.savename.value,
            dataType: "json",
            success: (data) => {
                console.log(data);
                if (data.success === true) {
                    alert(data.data)
                    this.updateSavesList();
                    this.setState({loading: false});
                } else {
                    alert(data.data)
                    document.getElementById('uploadsave').removeChild(loading);
                    this.setState({loading: false});
                }
            }
        })
    }

    render() {
        var loadingOverlay
        if (this.state.loading) {
            loadingOverlay = 
                <div className="overlay">
                    <i className="fa fa-refresh fa-spin"></i>
                </div>
        } else {
            loadingOverlay = ""
        }

        return(
            <div className="box" id="uploadsave">
                <div className="box-header">
                    <h4 className="box-title">Create Save File</h4>
                </div>

                <div className="box-body">
                    <form>
                        <div className="form-group">
                            <label for="savefile">Enter Savefile Name... </label>
                            <input className="form-control" ref="savename" type="text" name="savefile" id="savefilename" />
                        </div>
                        <div className="form-group">
                            <input className="form-control btn btn-default" type="button" ref="button" value="Create" onClick={this.createSaveFile} />
                        </div>
                    </form> 
                </div>
                {loadingOverlay}
            </div>
        )
    }
}

CreateSave.propTypes = {
    getSaves: React.PropTypes.func.isRequired,
}

export default CreateSave
