import React from 'react';
import PropTypes from 'prop-types';
import NativeListener from 'react-native-listener';
import ModSearch from './search/ModSearch.jsx';
import ModUpload from "./ModUpload.jsx";
import ModManager from "./ModManager.jsx";
import ModPacks from "./packs/ModPackOverview.jsx";
import {instanceOfModsContent} from "./ModsPropTypes.js";
import ModLoadSave from "./ModLoadSave.jsx";
import {ReactSwalNormal} from 'Utilities/customSwal';
import FontAwesomeIcon from "../FontAwesomeIcon";

class ModOverview extends React.Component {
    constructor(props) {
        super(props);

        this.handlerSearchMod = this.handlerSearchMod.bind(this);

        this.state = {
            shownModList: []
        }
    }

    handlerSearchMod(e) {
        e.preventDefault();

        $.ajax({
            url: "/api/mods/search",
            method: "GET",
            data: $(e.target).serialize(),
            dataType: "JSON",
            success: (data) => {
                let parsed_data = JSON.parse(data.data);

                this.setState({
                    "shownModList": parsed_data.results
                });
            },
            error: (jqXHR) => {
                let json_data = JSON.parse(jqXHR.responseJSON.data);

                ReactSwalNormal.fire({
                    title: json_data.detail,
                    icon: "error"
                });
            }
        })
    }

    downloadAllHandler(e) {
        e.stopPropagation();
    }

    render() {
        return(
            <div>
                <div className="card collapsed-card" id="add-mod-box">
                    <div className="card-header">
                        <button type="button" className="btn btn-tool btn-collapse" data-card-widget="collapse">
                            <FontAwesomeIcon icon="plus"/>
                        </button>
                        <h3 className="card-title">Add Mod</h3>
                        {this.props.loggedIn ?
                            <div className="card-tools">
                                <NativeListener onClick={this.props.factorioLogoutHandler}>
                                    <button className="btn btn-tool btn-danger">
                                        Logout
                                    </button>
                                </NativeListener>
                            </div>
                        : null}
                    </div>

                    <ModSearch
                        {...this.state}
                        {...this.props}
                        submitSearchMod={this.handlerSearchMod}
                        submitFactorioLogin={this.props.submitFactorioLogin}
                    />
                </div>

                <div className="card collapsed-card">
                    <div className="card-header">
                        <button type="button" className="btn btn-tool btn-collapse" data-card-widget="collapse">
                            <FontAwesomeIcon icon="plus"/>
                        </button>
                        <h3 className="card-title">Upload Mod</h3>
                    </div>

                    <ModUpload
                        {...this.props}
                    />
                </div>

                <div className="card collapsed-card">
                    <div className="card-header">
                        <button type="button" className="btn btn-tool btn-collapse" data-card-widget="collapse">
                            <FontAwesomeIcon icon="plus"/>
                        </button>
                        <h3 className="card-title">Load Mods From Save</h3>
                    </div>

                    <ModLoadSave
                        {...this.props}
                    />
               </div>

                <div className="card" id="manage-mods">
                    <div className="card-header">
                        <button type="button" className="btn btn-tool btn-collapse" data-card-widget="collapse">
                            <FontAwesomeIcon icon="minus"/>
                        </button>
                        <h3 className="card-title">Manage Mods</h3>
                        <div className="card-tools">
                            {
                                this.props.installedMods != null ?
                                    <NativeListener onClick={this.downloadAllHandler}>
                                        <a className="btn btn-tool btn-default" href={"/api/mods/download"} download>
                                            Download all Mods
                                        </a>
                                    </NativeListener>
                                    : null
                            }
                            {
                                this.props.updatesAvailable > 0 ?
                                    <NativeListener onClick={this.props.updateAllMods}>
                                        <button className="btn btn-tool btn-default">
                                            Update all Mods
                                        </button>
                                    </NativeListener>
                                    : null
                            }
                            {
                                this.props.installedMods != null ?
                                    <NativeListener onClick={this.props.deleteAll}>
                                        <button className="btn btn-tool btn-danger">
                                            Delete ALL Mods
                                        </button>
                                    </NativeListener>
                                    : null
                            }
                        </div>
                    </div>

                    <ModManager
                        {...this.props}
                        id="modManager"
                    />
                </div>

                <div className="card collapsed-card">
                    <div className="card-header">
                        <button type="button" className="btn btn-tool btn-collapse" data-card-widget="collapse">
                            <FontAwesomeIcon icon="plus"/>
                        </button>
                        <h3 className="card-title">Manage Modpacks</h3>
                    </div>

                    <ModPacks
                        {...this.props}
                        id="modPacks"
                    />
                </div>
            </div>
        );
    }
}

ModOverview.propTypes = {
    installedMods: PropTypes.array,
    submitFactorioLogin: PropTypes.func.isRequired,
    toggleMod: PropTypes.func.isRequired,
    deleteMod: PropTypes.func.isRequired,
    deleteAll: PropTypes.func.isRequired,
    updateMod: PropTypes.func.isRequired,
    uploadModSuccessHandler: PropTypes.func.isRequired,
    loggedIn: PropTypes.bool.isRequired,
    factorioLogoutHandler: PropTypes.func.isRequired,
    updatesAvailable: PropTypes.number.isRequired,
    updateAllMods: PropTypes.func.isRequired,
    updateCountAdd: PropTypes.func.isRequired,

    modContentClass: instanceOfModsContent.isRequired,
};

export default ModOverview;