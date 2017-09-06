import React from "react";
import Mod from "./Mod.jsx";

class ModManager extends React.Component {
    render() {
        return (
            <div className="box-body">
                <div className="table-responsive">
                    <table className="table table-striped">
                        <thead>
                        <tr>
                            <th>Name</th>
                            <th>Status</th>
                            <th>Version</th>
                            <th>Toggle/Remove</th>
                        </tr>
                        </thead>
                        <tbody>
                        {(this.props.installedMods != null) ?
                            this.props.installedMods.map((mod, i) => {
                                if(mod.name !== "base")
                                    return(
                                        <Mod
                                            key={mod.name}
                                            mod={mod}
                                            {...this.props}
                                        />
                                    )
                            }):null}
                        </tbody>
                    </table>
                </div>
            </div>
        )
    }
}

ModManager.propTypes = {
    installedMods: React.PropTypes.array,
    toggleMod: React.PropTypes.func.isRequired,
    deleteMod: React.PropTypes.func.isRequired,
    updateMod: React.PropTypes.func.isRequired,
}

export default ModManager;