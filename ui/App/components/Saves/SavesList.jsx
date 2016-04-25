import React from 'react';
import Save from './Save.jsx';

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
        $.ajax({
            url: "/api/saves/rm/" + saveName,
            success: (data) => {
                alert(data)
            }
        })
        this.updateSavesList();
    }

    render() {
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
                            {this.props.saves.map ( (save, i) => {
                                return(
                                    <Save
                                        key={i}
                                        saves={this.props.saves}
                                        index={i}
                                        save={save}
                                        removeSave={this.removeSave}
                                    />
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

SavesList.propTypes = {
    saves: React.PropTypes.array.isRequired,
    dlSave: React.PropTypes.func.isRequired,
    getSaves: React.PropTypes.func.isRequired
}

export default SavesList
