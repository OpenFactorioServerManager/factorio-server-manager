import React from 'react';
import Mod from './Mod';

class ModOverview extends React.Component {
    render() {
        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Manage Mods</h3>
                </div>

                <div className="box-body">
                    <div className="table-responsive">
                        <table className="table table-striped">
                            <thead>
                            <tr>
                                <th>Name</th>
                                <th>Status</th>
                                <th>Toggle Status</th>
                            </tr>
                            </thead>
                            <tbody>
                            {this.props.installedMods.map ( (mod, i) => {
                                return(
                                    <Mod
                                        key={i}
                                        mod={mod}
                                        {...this.props}
                                    />
                                )
                            })}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        );
    }
}

ModOverview.propTypes = {
    installedMods: React.PropTypes.array.isRequired
}

export default ModOverview;