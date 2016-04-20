import React from 'react';
import Mod from './Mod.jsx'

class ModList extends React.Component {
    componentDidMount() {
        console.log(this.props.listMods);
    }

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
                        {this.props.listMods.map ( (mod, i) => {
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
        )
    }
}

ModList.propTypes = {
    listMods: React.PropTypes.array.isRequired,
    toggleMod: React.PropTypes.func.isRequired
}

export default ModList
