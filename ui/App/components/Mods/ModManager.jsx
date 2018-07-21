import React from "react";
import PropTypes from 'prop-types';
import Mod from "./Mod.jsx";

class ModManager extends React.Component {
    constructor(props) {
        super(props);
    }

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
                            <th>Factorio Version</th>
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
    installedMods: PropTypes.array,
    toggleMod: PropTypes.func.isRequired,
    deleteMod: PropTypes.func.isRequired,
    updateMod: PropTypes.func.isRequired,
    updateCountAdd: PropTypes.func,
}

export default ModManager;