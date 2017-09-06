import React from 'react';
import ModManager from "../ModManager.jsx";

class ModPackOverview extends React.Component {
    constructor(props) {
        super(props);

        this.createModPack = this.createModPack.bind(this);

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

                    let json_response = jqXHR.responseJSON
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
        e.preventDefault();
        e.nativeEvent.stopImmediatePropagation();
    }

    test() {
        console.log("test called");
    }

    render() {
        return(
            <div className="box-body">
                {
                    this.state.listPacks.map(
                        (modpack, index) => {
                            return(
                                <div className="box collapsed-box">
                                    <div className="box-header" data-widget="collapse" style={{cursor: "pointer"}}>
                                        <i className="fa fa-plus"></i>
                                        <h3 className="box-title">{modpack.name}</h3>
                                        <div className="box-tools pull-right">
                                            <button className="btn btn-box-tool btn-danger" style={{color: "#fff"}} onClick={this.deleteModPack}>Delete</button>
                                        </div>
                                    </div>
                                    <div className="box-body">
                                        <ModManager
                                            installedMods={modpack.mods.mods}
                                            deleteMod={this.test}
                                            toggleMod={this.test}
                                            updateMod={this.test}
                                        />
                                    </div>
                                </div>
                            )
                        }
                    )
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

};

export default ModPackOverview
