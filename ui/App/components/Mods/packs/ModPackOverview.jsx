import React from 'react';
import ModManager from "../ModManager.jsx";
import NativeListener from 'react-native-listener';
import {instanceOfModsContent} from "../ModsPropTypes.js";
import locks from "locks";

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
        let this_class = this;
        swal({
            title: "Create modpack",
            text: "Please enter an unique modpack name:",
            type: "input",
            showCancelButton: true,
            closeOnConfirm: false,
            inputPlaceholder: "Modpack name",
            showLoaderOnConfirm: true
        },
        function(inputValue){
            if (inputValue === false) return false;

            if (inputValue === "") {
                swal.showInputError("A modpack needs a name!");
                return false
            }

            $.ajax({
                url: "/api/mods/packs/create",
                method: "POST",
                data: {name: inputValue},
                dataType: "JSON",
                success: (data) => {
                    this_class.mutex.lock(() => {
                        let packList = this_class.state.listPacks;

                        data.data.mod_packs.forEach((v, k) => {
                            if(v.name == inputValue) {
                                packList.push(data.data.mod_packs[k]);
                                return false;
                            }
                        });

                        this_class.setState({
                            listPacks: packList
                        });

                        swal({
                            title: "modpack created successfully",
                            type: "success"
                        });

                        this_class.mutex.unlock();
                    });
                },
                error: (jqXHR, status, err) => {
                    console.log('api/mods/packs/create', status, err.toString());

                    let json_response = jqXHR.responseJSON;
                    swal({
                        title: "Error on creating modpack",
                        text: json_response.data,
                        type: "error"
                    });
                }
            });
        });
    }

    deleteModPack(e) {
        e.stopPropagation();

        let this_class = this;
        let name = $(e.target).parent().prev().html();

        swal({
            title: "Are you sure?",
            text: "You really want to delete this modpack?\nThere is no turning back, the modpack will be deleted forever (a very long time)!",
            type: "info",
            showCancelButton: true,
            closeOnConfirm: false,
            showLoaderOnConfirm: true
        },
        function() {
            $.ajax({
                url: "/api/mods/packs/delete",
                method: "POST",
                data: {name: name},
                dataType: "JSON",
                success: (data) => {
                    if(data.success) {
                        this_class.mutex.lock(() => {
                            let mod_packs = this_class.state.listPacks;

                            mod_packs.forEach((v, k) => {
                                if(v.name == name) {
                                    delete mod_packs[k];
                                }
                            });

                            this_class.setState({
                                listPacks: mod_packs
                            });

                            swal({
                                title: "Modpack deleted successfully",
                                type: "success"
                            });

                            this_class.mutex.unlock();
                        });
                    }
                },
                error: (jqXHR, status, err) => {
                    console.log('api/mods/packs/delete', status, err.toString());

                    let json_response = jqXHR.responseJSON || err.toString();
                    json_response = json_response.data || err.toString();

                    swal({
                        title: "Error on creating modpack",
                        text: json_response,
                        type: "error"
                    });
                }
            })
        });
    }

    loadModPack(e) {
        e.stopPropagation();

        let this_class = this
        let name = $(e.target).parent().prev().html();

        swal({
            title: "Are you sure?",
            text: "This operation will replace the current installed mods with the mods out of the selected ModPack!",
            type: "info",
            showCancelButton: true,
            closeOnConfirm: false,
            showLoaderOnConfirm: true
        },
        function() {
            $.ajax({
                url: "/api/mods/packs/load",
                method: "POST",
                data: {name: name},
                dataType: "JSON",
                success: (data) => {
                    swal({
                        title: "ModPack loaded!",
                        type: "success"
                    });

                    this_class.props.modContentClass.setState({
                        installedMods: data.data.mods
                    });
                },
                error: (jqXHR, status, err) => {
                    console.log('api/mods/packs/load', status, err.toString());

                    let json_response = jqXHR.responseJSON || err.toString();
                    json_response = json_response.data || err.toString();

                    swal({
                        title: "Error on loading ModPack",
                        text: json_response,
                        type: "error"
                    });
                }
            })
        });
    }

    downloadModPack(e) {
        e.stopPropagation();
    }

    modPackToggleModHandler(e, updatesInProgress) {
        e.preventDefault();


        if(updatesInProgress) {
            swal("Toggle mod failed", "Can't toggle the mod, when an update is still in progress", "error");
            return false;
        }

        let $button = $(e.target);
        let $row = $button.parents("tr");
        let mod_name = $row.data("mod-name");
        let mod_pack = $row.parents(".single-modpack").find("h3").html();
        let this_class = this;

        $.ajax({
            url: "/api/mods/packs/mod/toggle",
            method: "POST",
            data: {
                mod_name: mod_name,
                mod_pack: mod_pack
            },
            dataType: "JSON",
            success: (data) => {
                if(data.success) {
                    this_class.mutex.lock(() => {
                        let packList = this_class.state.listPacks;

                        packList.forEach((modPack, modPackKey) => {
                            if(modPack.name == mod_pack) {
                                packList[modPackKey].mods.mods.forEach((mod, modKey) => {
                                    if(mod.name == mod_name) {
                                        packList[modPackKey].mods.mods[modKey].enabled = data.data;
                                        return false;
                                    }
                                });
                            }
                        });

                        this_class.setState({
                            listPacks: packList
                        });

                        this_class.mutex.unlock();
                    });
                }
            },
            error: (jqXHR, status, err) => {
                console.log('api/mods/packs/mod/toggle', status, err.toString());
                swal({
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
            swal("Delete failed", "Can't delete the mod, when an update is still in progress", "error");
            return false;
        }

        let $button = $(e.target);
        let $row = $button.parents("tr");
        let mod_name = $row.data("mod-name");
        let mod_pack = $row.parents(".single-modpack").find("h3").html();
        let class_this = this;

        swal({
            title: "Delete Mod?",
            text: "This will delete the mod forever",
            type: "info",
            showCancelButton: true,
            closeOnConfirm: false,
            confirmButtonText: "Delete it!",
            cancelButtonText: "Close",
            showLoaderOnConfirm: true,
            confirmButtonColor: "#DD6B55",
        }, function () {
            $.ajax({
                url: "/api/mods/packs/mod/delete",
                method: "POST",
                data: {
                    mod_name: mod_name,
                    mod_pack_name: mod_pack
                },
                dataType: "JSON",
                success: (data) => {
                    if(data.success) {
                        class_this.mutex.lock(() => {
                            swal("Delete of mod " + mod_name + " inside modPack " + mod_pack + " successful", "", "success");

                            let packList = class_this.state.listPacks;

                            packList.forEach((modPack, modPackKey) => {
                                if(modPack.name == mod_pack) {
                                    packList[modPackKey].mods.mods.forEach((mod, modKey) => {
                                        if(mod.name == mod_name) {
                                            delete packList[modPackKey].mods.mods[modKey];
                                            return false;
                                        }
                                    });
                                }
                            });

                            class_this.setState({
                                listPacks: packList
                            });

                            class_this.mutex.unlock();
                        });
                    }
                },
                error: (jqXHR, status, err) => {
                    console.log('api/mods/packs/mod/delete', status, err.toString());
                    swal({
                        title: "Delete Mod went wrong",
                        text: jqXHR.responseJSON.data,
                        type: "error"
                    });
                }
            });
        });
    }

    modPackUpdateModHandler(e, toggleUpdateStatus, removeVersionAvailableStatus) {
        e.preventDefault();

        if(!this.props.modContentClass.state.logged_in) {
            swal({
                type: "error",
                title: "Update failed",
                text: "please login into Factorio to update mod"
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
            let modname = $row.data("modName");
            let mod_pack = $row.parents(".single-modpack").find("h3").html();

            let this_class = this;

            //make button spinning
            toggleUpdateStatus();

            $.ajax({
                url: "/api/mods/packs/mod/update",
                method: "POST",
                data: {
                    downloadUrl: download_url,
                    filename: filename,
                    mod_name: modname,
                    mod_pack_name: mod_pack
                },
                success: (data) => {
                    toggleUpdateStatus();
                    removeVersionAvailableStatus();

                    if(data.success) {
                        this_class.mutex.lock(() => {
                            let packList = this_class.state.listPacks;

                            packList.forEach((modPack, modPackKey) => {
                                if(modPack.name == mod_pack) {
                                    packList[modPackKey].mods.mods.forEach((mod, modKey) => {
                                        if(mod.name == modname) {
                                            packList[modPackKey].mods.mods[modKey] = data.data;
                                            return false;
                                        }
                                    });
                                }
                            });

                            this_class.setState({
                                listPacks: packList
                            });

                            this_class.mutex.unlock();
                        });
                    }
                },
                error: (jqXHR, status, err) => {
                    console.log('api/mods/packs/mod/update', status, err.toString());
                    toggleUpdateStatus();
                    swal({
                        title: "Update Mod went wrong",
                        text: jqXHR.responseJSON.data,
                        type: "error"
                    });
                }
            });
        }
    }

    test() {
        console.log("test called");
    }

    render() {
        return(
            <div className="box-body">
                {
                    this.state.listPacks != null ?
                        this.state.listPacks.map(
                            (modpack, index) => {
                                return(
                                    <div key={modpack.name} className="box single-modpack collapsed-box">
                                        <div className="box-header" data-widget="collapse" style={{cursor: "pointer"}}>
                                            <i className="fa fa-plus"></i>
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
                                        <div className="box-body">
                                            <ModManager
                                                installedMods={modpack.mods.mods}
                                                deleteMod={this.modPackDeleteModHandler}
                                                toggleMod={this.modPackToggleModHandler}
                                                updateMod={this.modPackUpdateModHandler}
                                            />
                                        </div>
                                    </div>
                                )
                            }
                        )
                    : null
                }

                <div className="box">
                    <div className="box-header" style={{cursor: "pointer"}} onClick={this.createModPack}>
                        <i className="fa fa-plus"></i>
                        <h3 className="box-title">Add ModPack with current installed mods</h3>
                    </div>
                </div>
            </div>
        );
    }
}

ModPackOverview.propTypes = {
    modContentClass: instanceOfModsContent.isRequired,
};

export default ModPackOverview
