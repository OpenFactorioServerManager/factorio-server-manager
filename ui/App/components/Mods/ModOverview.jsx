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
                    type: "error"
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
                <div className="box collapsed-box" id="add-mod-box">
                    <div className="box-header"
                         data-toggle="collapse"
                         data-target="#modSearch"
                         aria-expanded="false"
                         aria-controls="modSearch"
                         role="button"
                         style={{cursor: "pointer"}}
                    >
                        <FontAwesomeIcon icon="plus"/>
                        <h3 className="box-title">Add Mod</h3>
                        {this.props.loggedIn ?
                            <div className="box-tools pull-right">
                                <NativeListener onClick={this.props.factorioLogoutHandler}>
                                    <button className="btn btn-box-tool btn-danger" style={{color: "#fff"}}>
                                        Logout
                                    </button>
                                </NativeListener>
                            </div>
                        : null}
                    </div>

                    <ModSearch
                        {...this.state}
                        {...this.props}
                        className="collapse"
                        id="modSearch"
                        submitSearchMod={this.handlerSearchMod}
                        submitFactorioLogin={this.props.submitFactorioLogin}
                    />
                </div>

                <div className="box collapsed-box">
                    <div className="box-header"
                         data-toggle="collapse"
                         data-target="#modUpload"
                         aria-expanded="false"
                         aria-controls="modUpload"
                         role="button"
                         style={{cursor: "pointer"}}
                    >
                        <FontAwesomeIcon icon="plus"/>
                        <h3 className="box-title">Upload Mod</h3>
                    </div>

                    <ModUpload
                        {...this.props}
                        className="collapse"
                        id="modUpload"
                    />
                </div>

                <div className="box collapsed-box">
                    <div className="box-header"
                         data-toggle="collapse"
                         data-target="#modLoadSave"
                         aria-expanded="false"
                         aria-controls="modLoadSave"
                         role="button"
                         style={{cursor: "pointer"}}
                    >
                        <FontAwesomeIcon icon="plus"/>
                        <h3 className="box-title">Load Mods From Save</h3>
                    </div>

                    <ModLoadSave
                        {...this.props}
                        className="collapse"
                        id="modLoadSave"
                    />
               </div>

                <div className="box" id="manage-mods">
                    <div className="box-header"
                         data-toggle="collapse"
                         data-target="#modManager"
                         aria-expanded="true"
                         aria-controls="modManager"
                         role="button"
                         style={{cursor: "pointer"}}
                    >
                        <FontAwesomeIcon icon="minus"/>
                        <h3 className="box-title">Manage Mods</h3>
                        <div className="box-tools float-sm-right">
                            {
                                this.props.installedMods != null ?
                                    <NativeListener onClick={this.downloadAllHandler}>
                                        <a className="btn btn-box-tool btn-default" style={{marginRight: 20}} href={"/api/mods/download"} download>
                                            Download all Mods
                                        </a>
                                    </NativeListener>
                                    : null
                            }
                            {
                                this.props.updatesAvailable > 0 ?
                                    <NativeListener onClick={this.props.updateAllMods}>
                                        <button className="btn btn-box-tool btn-default" style={{marginRight: 20}}>
                                            Update all Mods
                                        </button>
                                    </NativeListener>
                                    : null
                            }
                            {
                                this.props.installedMods != null ?
                                    <NativeListener onClick={this.props.deleteAll}>
                                        <button className="btn btn-box-tool btn-danger" style={{color: "#fff"}}>
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
                        className="show"
                    />
                </div>

                <div className="box collapsed-box">
                    <div className="box-header"
                         data-toggle="collapse"
                         data-target="#modPacks"
                         aria-expanded="false"
                         aria-controls="modPacks"
                         role="button"
                         style={{cursor: "pointer"}}
                    >
                        <FontAwesomeIcon icon="plus"/>
                        <h3 className="box-title">Manage Modpacks</h3>
                    </div>

                    <ModPacks
                        {...this.props}
                        className="collapse"
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