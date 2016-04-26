import React from 'react';
import SavesList from './Saves/SavesList.jsx';
import CreateSave from './Saves/CreateSave.jsx';
import UploadSave from './Saves/UploadSave.jsx';

class SavesContent extends React.Component {
    constructor(props) {
        super(props);
        this.dlSave = this.dlSave.bind(this);
    }

    componentDidMount() {
        this.props.getSaves();
    }


    dlSave(saveName) {
        $.ajax({
            url: "/api/saves/dl/" + saveName,
            dataType: "json",
            success: (data) => {
                console.log("Downloading save: " + saveName)
            },
            error: (xhr, status, err) => {
                console.log('api/mods/list', status, err.toString());
            }
        })
    }

    render() {
        return(
            <div className="content-wrapper">
                <section className="content-header">
                <h1>
                    Saves
                    <small>Factorio Save Files</small>
                </h1>
                <ol className="breadcrumb">
                    <li><a href="#"><i className="fa fa-dashboard"></i> Level</a></li>
                    <li className="active">Here</li>
                </ol>
                </section>

                <section className="content">
                <div className="row">
                    <div className="col-md-6">
                        <CreateSave 
                            getSaves={this.props.getSaves}
                        />
                    </div>
                    <div className="col-md-6">
                        <UploadSave 
                            getSaves={this.props.getSaves}
                        />
                    </div>
                </div>

                <SavesList 
                    {...this.state}
                    saves={this.props.saves}
                    dlSave={this.dlSave}
                    getSaves={this.props.getSaves}
                />


                </section>
            </div>
        )
    }
}

export default SavesContent
