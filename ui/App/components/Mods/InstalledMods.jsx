import React from 'react';

class InstalledMods extends React.Component {
    render() {
        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Installed Mods</h3>
                </div>
                     
                <div className="box-body">
                {this.props.installedMods.map ( (mod, i) => {
                    return(
                        <p>{mod}</p>
                    )                                            
                })}
                </div>
            </div>
        )
    }
}

InstalledMods.propTypes = {
    installedMods: React.PropTypes.array.isRequired
}

export default InstalledMods
