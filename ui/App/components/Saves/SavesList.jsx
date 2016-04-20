import React from 'react';
import Save from './Save.jsx';

class SavesList extends React.Component {
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
                                </tr>
                            </thead>
                            <tbody>
                            {this.props.saves.map ( (save, i) => {
                                return(
                                    <Save
                                        key={i}
                                        save={save}
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
    dlSave: React.PropTypes.func.isRequired
}

export default SavesList
