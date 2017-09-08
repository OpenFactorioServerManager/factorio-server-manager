import React from 'react';
import ModManager from "../ModManager.jsx";
import NativeListener from 'react-native-listener';
import {instanceOfModsContent} from "../ModsPropTypes.js";

class ModPackOverview extends React.Component {
    constructor(props) {
        super(props);

        this.createModPack = this.createModPack.bind(this);
        this.deleteModPack = this.deleteModPack.bind(this);
        this.loadModPack = this.loadModPack.bind(this);

        this.state = {
            listPacks: []
        }
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
                console.log(data);
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
                    this_class.setState({
                        listPacks: data.data.mod_packs
                    });

                    swal({
                        title: "modpack created successfully",
                        type: "success"
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
                    this_class.setState({
                        listPacks: data.data.mod_packs
                    });

                    swal({
                        title: "Modpack deleted successfully",
                        type: "success"
                    });
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
                        installedMods: data.data
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
                                    <div className="box collapsed-box">
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
                                                deleteMod={this.test} //TODO
                                                toggleMod={this.test}
                                                updateMod={this.test}
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
