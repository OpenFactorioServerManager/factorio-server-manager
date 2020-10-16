import React from 'react';
import PropTypes from 'prop-types';
import Save from './Save.jsx';
import {ReactSwalDanger, ReactSwalNormal} from 'Utilities/customSwal';

class SavesList extends React.Component {
   constructor(props) {
        super(props);
        this.updateSavesList = this.updateSavesList.bind(this);
        this.removeSave = this.removeSave.bind(this);
   }

    updateSavesList () {
        this.props.getSaves();
    }

    removeSave(saveName, e) {
        let self = this;
        ReactSwalDanger.fire({
            title: "Are you sure?",
            html: <p>Save: {saveName} will be deleted</p>,
            icon: "question",
            showCancelButton: true,
            confirmButtonText: "Yes, delete it!",
            showLoaderOnConfirm: true,
            preConfirm: () => {
                return new Promise((resolve, reject) => {
                    $.ajax({
                        url: "/api/saves/rm/" + saveName,
                        dataType: "json",
                        success: (resp) => {
                            if (resp.success === true) {
                                resolve(resp.data);
                            } else {
                                reject("Unknown occurred!");
                            }
                        },
                        error: () => {
                            reject("Unknown occurred!");
                        }
                    });
                });
            },
        }).then((result) => {
            if(result.value) {
                ReactSwalNormal.fire({
                    title: "Deleted!",
                    text: result.value,
                    icon: "success"
                });
            }
            self.updateSavesList();
        }).catch((result) => {
            ReactSwalNormal.fire({
                title: "An error occurred!",
                text: result,
                icon: "error"
            });
        });
    }

    render() {
        let savesList;
        if (this.props.saves.length === 0) {
            savesList = <tr></tr>
        } else {
            savesList = this.props.saves.map ( (save, i) => {
                return(
                    <Save
                        key={i}
                        saves={this.props.saves}
                        index={i}
                        save={save}
                        removeSave={this.removeSave}
                    />
                )
            });
            
        }

        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Save Files</h3>
                </div>
                
                <div className="box-body">

                    <div className="table-responsive">
                        <table className="table table-striped">
                            <thead>
                                <tr>
                                    <th>Filname</th>
                                    <th>Last Modified Time</th>
                                    <th>Filesize</th>
                                    <th>Download</th>
                                </tr>
                            </thead>
                            <tbody>
                                {savesList}
                            </tbody>
                        </table>        
                    </div>
                </div>
            </div>
        )
    }
}

SavesList.propTypes = {
    saves: PropTypes.array.isRequired,
    dlSave: PropTypes.func.isRequired,
    getSaves: PropTypes.func.isRequired
};

export default SavesList
