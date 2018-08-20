import React from 'react';
import ModManager from "../ModManager.jsx";
import NativeListener from 'react-native-listener';
import {instanceOfModsContent} from "../ModsPropTypes.js";
import locks from "locks";
import PropTypes from "prop-types";
import {ReactSwalNormal, ReactSwalDanger} from 'Utilities/customSwal';
import FontAwesomeIcon from "../../FontAwesomeIcon";

class ModPackOverview extends React.Component {
    constructor(props) {
        super(props);

        this.createModPack = this.createModPack.bind(this);
        this.deleteModPack = this.deleteModPack.bind(this);
        this.loadModPack = this.loadModPack.bind(this);
        this.modPackToggleModHandler = this.modPackToggleModHandler.bind(this);
        this.modPackDeleteModHandler = this.modPackDeleteModHandler.bind(this);
        this.modPackUpdateModHandler = this.modPackUpdateModHandler.bind(this);

        this.state = {
            listPacks: []
        }

        this.mutex = locks.createMutex();
    }

    componentDidMount() {
        this.getModPacks();
    }

    getModPacks() {
        //send ajax to get all modPacks and setState
        $.ajax({
            url: "/api/mods/packs/list",
            method: "GET",
            dataType: "JSON",
            success: (data) => {
                this.setState({
                    listPacks: data.data.mod_packs
                });
            },
            error: (jqXHR, status, err) => {
                console.log('api/mods/packs/list', status, err.toString());
            }
        })
    }

    createModPack() {
        ReactSwalNormal.fire({
            title: "Create modpack",
            html: "Please enter an unique modpack name:",
            input: "text",
            showCancelButton: true,
            inputPlaceholder: "Modpack name",
            inputAttributes: {
                required: "required"
            },
            inputValidator: (value) => {
                return new Promise(resolve => {
                    if(value) {
                        resolve();
                    } else {
                        resolve("You need to enter a name");
                    }
                });
            },
            showLoaderOnConfirm: true,
            preConfirm: (inputValue) => {
                // console.log(this);

                $.ajax({
                    url: "/api/mods/packs/create",
                    method: "POST",
                    data: {name: inputValue},
                    dataType: "JSON",
                    success: (data) => {
                        this.mutex.lock(() => {
                            let packList = this.state.listPacks;

                            data.data.mod_packs.forEach((v, k) => {
                                if(v.name == inputValue) {
                                    packList.push(data.data.mod_packs[k]);
                                    return false;
                                }
                            });

                            this.setState({
                                listPacks: packList
                            });

                            ReactSwalNormal.fire({
                                title: "modpack created successfully",
                                type: "success"
                            });

                            this.mutex.unlock();
                        });
                    },
                    error: (jqXHR, status, err) => {
                        console.log('api/mods/packs/create', status, err.toString());

                        let jsonResponse = jqXHR.responseJSON;

                        ReactSwalNormal.fire({
                            title: "Error on creating modpack",
                            text: jsonResponse.data,
                            type: "error"
                        });
                    }
                });
            }
        });
    }

    deleteModPack(e) {
        e.stopPropagation();

        let name = $(e.target).parent().prev().html();

        ReactSwalDanger.fire({
            title: "Are you sure?",
            html: <p>You really want to delete this modpack?<br/>There is no turning back, the modpack will be deleted forever (a very long time)!</p>,
            type: "question",
            showCancelButton: true,
            showLoaderOnConfirm: true,
            preConfirm: () => {
                $.ajax({
                    url: "/api/mods/packs/delete",
                    method: "POST",
                    data: {name: name},
                    dataType: "JSON",
                    success: (data) => {
                        if(data.success) {
                            this.mutex.lock(() => {
                                let modPacks = this.state.listPacks;

                                modPacks.forEach((v, k) => {
                                    if(v.name == name) {
                                        delete modPacks[k];
                                    }
                                });

                                this.setState({
                                    listPacks: modPacks
                                });

                                ReactSwalNormal.fire({
                                    title: "Modpack deleted successfully",
                                    type: "success"
                                });

                                this.mutex.unlock();
                            });
                        }
                    },
                    error: (jqXHR, status, err) => {
                        console.log('api/mods/packs/delete', status, err.toString());

                        let jsonResponse = jqXHR.responseJSON || err.toString();
                        jsonResponse = jsonResponse.data || err.toString();

                        ReactSwalNormal.fire({
                            title: "Error on creating modpack",
                            text: jsonResponse,
                            type: "error"
                        });
                    }
                })
            }
        });
    }

    loadModPack(e) {
        e.stopPropagation();

        let name = $(e.target).parent().prev().html();

        ReactSwalDanger.fire({
            title: "Are you sure?",
            text: "This operation will replace the current installed mods with the mods out of the selected ModPack!",
            type: "question",
            showCancelButton: true,
            showLoaderOnConfirm: true,
            preConfirm: () => {
                $.ajax({
                    url: "/api/mods/packs/load",
                    method: "POST",
                    data: {name: name},
                    dataType: "JSON",
                    success: (data) => {
                        ReactSwalNormal.fire({
                            title: "ModPack loaded!",
                            type: "success"
                        });

                        this.props.modContentClass.setState({
                            installedMods: data.data.mods
                        });
                    },
                    error: (jqXHR, status, err) => {
                        console.log('api/mods/packs/load', status, err.toString());

                        let jsonResponse = jqXHR.responseJSON || err.toString();
                        jsonResponse = jsonResponse.data || err.toString();

                        ReactSwalNormal.fire({
                            title: "Error on loading ModPack",
                            text: jsonResponse,
                            type: "error"
                        });
                    }
                })
            }
        });
    }

    downloadModPack(e) {
        e.stopPropagation();
    }

    modPackToggleModHandler(e, updatesInProgress) {
        e.preventDefault();


        if(updatesInProgress) {
            ReactSwalNormal.fir({
                title: "Toggle mod failed",
                text: "Can't toggle the mod, when an update is still in progress",
                type: "error"
            });
            return false;
        }

        let $button = $(e.target);
        let $row = $button.parents("tr");
        let modName = $row.data("mod-name");
        let modPackName = $row.parents(".single-modpack").find("h3").html();

        $.ajax({
            url: "/api/mods/packs/mod/toggle",
            method: "POST",
            data: {
                modName: modName,
                modPack: modPackName
            },
            dataType: "JSON",
            success: (data) => {
                if(data.success) {
                    this.mutex.lock(() => {
                        let packList = this.state.listPacks;

                        packList.forEach((modPack, modPackKey) => {
                            if(modPack.name == modPackName) {
                                packList[modPackKey].mods.mods.forEach((mod, modKey) => {
                                    if(mod.name == modName) {
                                        packList[modPackKey].mods.mods[modKey].enabled = data.data;
                                        return false;
                                    }
                                });
                            }
                        });

                        this.setState({
                            listPacks: packList
                        });

                        this.mutex.unlock();
                    });
                }
            },
            error: (jqXHR, status, err) => {
                console.log('api/mods/packs/mod/toggle', status, err.toString());
                ReactSwalNormal.fire({
                    title: "Toggle Mod went wrong",
                    text: err.toString(),
                    type: "error"
                });
            }
        });
    }

    modPackDeleteModHandler(e, updatesInProgress) {
        e.preventDefault();

        if(updatesInProgress) {
            ReactSwalNormal.fire({
                title: "Delete failed",
                text: "Can't delete the mod, when an update is still in progress",
                type: "error"
            });
            return false;
        }

        let $button = $(e.target);
        let $row = $button.parents("tr");
        let modName = $row.data("mod-name");
        let modPackName = $row.parents(".single-modpack").find("h3").html();

        ReactSwalDanger.fire({
            title: "Delete Mod?",
            text: "This will delete the mod forever",
            type: "question",
            showCancelButton: true,
            confirmButtonText: "Delete it!",
            cancelButtonText: "Close",
            showLoaderOnConfirm: true,
            preConfirm: () => {
                $.ajax({
                    url: "/api/mods/packs/mod/delete",
                    method: "POST",
                    data: {
                        modName: modName,
                        modPackName: modPackName
                    },
                    dataType: "JSON",
                    success: (data) => {
                        if(data.success) {
                            this.mutex.lock(() => {
                                ReactSwalNormal.fire({
                                    title: <p>Delete of mod {modName} inside modPack {modPackName} successful</p>,
                                    type: "success"
                                })

                                let packList = this.state.listPacks;

                                packList.forEach((modPack, modPackKey) => {
                                    if(modPack.name == modPackName) {
                                        packList[modPackKey].mods.mods.forEach((mod, modKey) => {
                                            if(mod.name == modName) {
                                                delete packList[modPackKey].mods.mods[modKey];
                                                return false;
                                            }
                                        });
                                    }
                                });

                                this.setState({
                                    listPacks: packList
                                });

                                this.mutex.unlock();
                            });
                        }
                    },
                    error: (jqXHR, status, err) => {
                        console.log('api/mods/packs/mod/delete', status, err.toString());
                        ReactSwalNormal.fire({
                            title: "Delete Mod went wrong",
                            text: jqXHR.responseJSON.data,
                            type: "error"
                        });
                    }
                });
            }
        });
    }

    modPackUpdateModHandler(e, toggleUpdateStatus, removeVersionAvailableStatus) {
        e.preventDefault();

        if(!this.props.modContentClass.state.loggedIn) {
            ReactSwalNormal.fire({
                title: "Update failed",
                text: "please login into Factorio to update mod",
                type: "error",
            });

            let $addModBox = $('#add-mod-box');
            if($addModBox.hasClass("collapsed-box")) {
                $addModBox.find(".box-header").click();
            }
        } else {
            let $button = $(e.currentTarget);
            let download_url = $button.data("downloadUrl");
            let filename = $button.data("fileName");
            let $row = $button.parents("tr");
            let modName = $row.data("modName");
            let modPackName = $row.parents(".single-modpack").find("h3").html();

            //make button spinning
            toggleUpdateStatus();

            $.ajax({
                url: "/api/mods/packs/mod/update",
                method: "POST",
                data: {
                    downloadUrl: download_url,
                    filename: filename,
                    modName: modName,
                    modPackName: modPackName
                },
                success: (data) => {
                    toggleUpdateStatus();
                    removeVersionAvailableStatus();

                    if(data.success) {
                        this.mutex.lock(() => {
                            let packList = this.state.listPacks;

                            packList.forEach((modPack, modPackKey) => {
                                if(modPack.name == modPackName) {
                                    packList[modPackKey].mods.mods.forEach((mod, modKey) => {
                                        if(mod.name == modName) {
                                            packList[modPackKey].mods.mods[modKey] = data.data;
                                            return false;
                                        }
                                    });
                                }
                            });

                            this.setState({
                                listPacks: packList
                            });

                            this.mutex.unlock();
                        });
                    }
                },
                error: (jqXHR, status, err) => {
                    console.log('api/mods/packs/mod/update', status, err.toString());
                    toggleUpdateStatus();
                    ReactSwalNormal.fire({
                        title: "Update Mod went wrong",
                        text: jqXHR.responseJSON.data,
                        type: "error"
                    });
                }
            });
        }
    }

    render() {
        let classes = "box-body" + " " + this.props.className;
        let ids = this.props.id;

        return(
            <div id={ids} className={classes}>
                {
                    this.state.listPacks != null ?
                        this.state.listPacks.map(
                            (modpack, index) => {
                                return(
                                    <div key={modpack.name} className="box single-modpack collapsed-box">
                                        <div className="box-header"
                                             data-toggle="collapse"
                                             data-target={"#" + modpack.name}
                                             aria-expanded="false"
                                             aria-controls={modpack.name}
                                             style={{cursor: "pointer"}}
                                        >
                                            <FontAwesomeIcon icon="plus"/>
                                            <h3 className="box-title">{modpack.name}</h3>
                                            <div className="box-tools pull-right">
                                                <NativeListener onClick={this.downloadModPack}>
                                                    <a className="btn btn-box-tool btn-default" style={{marginRight: 10}} href={"/api/mods/packs/download/" + modpack.name} download>Download</a>
                                                </NativeListener>

                                                <NativeListener onClick={this.loadModPack}>
                                                    <button className="btn btn-box-tool btn-default" style={{marginRight: 10}}>Load ModPack</button>
                                                </NativeListener>

                                                <NativeListener onClick={this.deleteModPack}>
                                                    <button className="btn btn-box-tool btn-danger" style={{color: "#fff"}}>Delete</button>
                                                </NativeListener>
                                            </div>
                                        </div>

                                        <ModManager
                                            {...this.props}
                                            className="collapse"
                                            id={modpack.name}
                                            installedMods={modpack.mods.mods}
                                            deleteMod={this.modPackDeleteModHandler}
                                            toggleMod={this.modPackToggleModHandler}
                                            updateMod={this.modPackUpdateModHandler}
                                        />
                                    </div>
                                )
                            }
                        )
                    : null
                }

                <div className="box">
                    <div className="box-header" style={{cursor: "pointer"}} onClick={this.createModPack}>
                        <FontAwesomeIcon icon="plus"/>
                        <h3 className="box-title">Add ModPack with current installed mods</h3>
                    </div>
                </div>
            </div>
        );
    }
}

ModPackOverview.propTypes = {
    modContentClass: instanceOfModsContent.isRequired,
    className: PropTypes.string,
    id: PropTypes.string
};

export default ModPackOverview
