import React from 'react';
import NativeListener from 'react-native-listener';
import ModSearch from './search/ModSearch.jsx';
import ModUpload from "./ModUpload.jsx";
import ModManager from "./ModManager.jsx";
import ModPacks from "./packs/ModPackOverview.jsx";
import {instanceOfModsContent} from "./ModsPropTypes.js";

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

                swal({
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
                    <div className="box-header" data-widget="collapse" style={{cursor: "pointer"}}>
                        <i className="fa fa-plus"></i>
                        <h3 className="box-title">Add Mod</h3>
                        {this.props.logged_in ?
                            <div className="box-tools pull-right">
                                <NativeListener onClick={this.props.factorioLogoutHandler}>
                                    <button className="btn btn-box-tool btn-danger" style={{color: "#fff"}}>Logout
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

                <div className="box collapsed-box">
                    <div className="box-header" data-widget="collapse" style={{cursor: "pointer"}}>
                        <i className="fa fa-plus"></i>
                        <h3 className="box-title">Upload Mod</h3>
                    </div>

                    <ModUpload
                        {...this.props}
                    />
                </div>

                <div className="box" id="manage-mods">
                    <div className="box-header" data-widget="collapse" style={{cursor: "pointer"}}>
                        <i className="fa fa-minus"></i>
                        <h3 className="box-title">Manage Mods</h3>
                        <div className="box-tools pull-right">
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
                                this.props.updates_available > 0 ?
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
                    />
                </div>

                <div className="box collapsed-box">
                    <div className="box-header" data-widget="collapse" style={{cursor: "pointer"}}>
                        <i className="fa fa-plus"></i>
                        <h3 className="box-title">Manage Modpacks</h3>
                    </div>

                    <ModPacks
                        {...this.props}
                    />
                </div>
            </div>
        );
    }
}

ModOverview.propTypes = {
    installedMods: React.PropTypes.array,
    submitFactorioLogin: React.PropTypes.func.isRequired,
    toggleMod: React.PropTypes.func.isRequired,
    deleteMod: React.PropTypes.func.isRequired,
    deleteAll: React.PropTypes.func.isRequired,
    updateMod: React.PropTypes.func.isRequired,
    uploadModSuccessHandler: React.PropTypes.func.isRequired,
    logged_in: React.PropTypes.bool.isRequired,
    factorioLogoutHandler: React.PropTypes.func.isRequired,
    updates_available: React.PropTypes.number.isRequired,
    updateAllMods: React.PropTypes.func.isRequired,
    updateCountAdd: React.PropTypes.func.isRequired,

    modContentClass: instanceOfModsContent.isRequired,
};

export default ModOverview;