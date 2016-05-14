import React from 'react';

class ModPacks extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Mod Packs</h3>
                </div>
                     
                <div className="box-body">
                    <table className="table table-striped">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Download</th>
                                <th>Delete</th>
                            </tr>
                        </thead>
                        <tbody>
                        {this.props.modPacks.map( (mod, i) => {
                            return(
                                <tr key={i}>
                                    <td>{mod}</td>    
                                </tr>
                            )
                        })}
                        </tbody>
                    </table>
                </div>
            </div>
        )
    }
}

ModPacks.propTypes = {
    modPacks: React.PropTypes.array.isRequired, 
}

export default ModPacks
