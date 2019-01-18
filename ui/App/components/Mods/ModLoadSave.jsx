import React from 'react';
import {instanceOfModsContent} from "./ModsPropTypes";
import PropTypes from "prop-types";
import {ReactSwalNormal} from 'Utilities/customSwal';

class ModLoadSave extends React.Component {
    constructor(props) {
        super(props);

        this.loadMods = this.loadMods.bind(this);
        this.loadModsSwalHandler = this.loadModsSwalHandler.bind(this);
    }

    componentDidMount() {
        //Load Saves
        this.props.getSaves();
    }

    loadMods(e) {
        e.preventDefault();

        $.ajax({
            url: "/api/mods/save/load",
            method: "POST",
            data: $(e.target).serialize(),
            dataType: "JSON",
            success: (data) => {
                let checkboxes = [];

                data.data.mods.forEach((mod) => {
                    if(mod.name == "base") return;

                    let singleCheckbox = <tr key={mod.name}>
                        <td>
                            {mod.name}
                            <input type="hidden" name="mod_name" value={mod.name}/>
                        </td>
                        <td>
                            {mod.version}
                            <input type="hidden" name="mod_version" value={mod.version}/>
                        </td>
                    </tr>

                    checkboxes.push(singleCheckbox);
                });

                if(checkboxes.length == 0) {
                    ReactSwalNormal.fire({
                        title: "No mods in this save!",
                        type: "error"
                    });
                    return;
                }

                let table = <div>
                    All Mods will be installed
                    <div style={{display: "flex", width: "100%", justifyContent: "center"}}>
                        <form id="swalForm">
                            <table>
                                <thead>
                                    <tr>
                                        <th>
                                            Name
                                        </th>
                                        <th>
                                            Version
                                        </th>
                                    </tr>
                                </thead>
                                <tbody>
                                        {checkboxes}
                                </tbody>
                            </table>
                        </form>
                    </div>
                </div>

                ReactSwalNormal.fire({
                    title: "Mods to install",
                    html: table,
                    type: 'question',
                    showCancelButton: true,
                    confirmButtonText: "Download Mods!",
                    showLoaderOnConfirm: true,
                    preConfirm: this.loadModsSwalHandler
                });
            },
            error: (jqXHR) => {
                ReactSwalNormal.fire({
                    title: jqXHR.responseJSON.data,
                    html: true,
                    type: "error",
                });
            }
        });
    }

    loadModsSwalHandler() {
        $.ajax({
            url: "/api/mods/install/multiple",
            method: "POST",
            dataType: "JSON",
            data: $("#swalForm").serialize(),
            success: (data) => {
                ReactSwalNormal.fire({
                    title: "All Mods installed successfully!",
                    type: "success"
                });

                this.props.modContentClass.setState({
                    installedMods: data.data.mods
                });
            },
            error: (jqXHR) => {
                let json_data = JSON.parse(jqXHR.responseJSON.data);

                ReactSwalNormal.fire({
                    title: json_data.detail,
                    type: "error",
                });
            }
        })
    }

    render() {
        let saves = [];
        this.props.saves.forEach((value, index) => {
            if(index != this.props.saves.length - 1) {
                saves.push(
                    <option key={index} value={value.name}>
                        {value.name}
                    </option>
                )
            }
        });

        let classes = "box-body" + " " + this.props.className;
        let ids = this.props.id;

        return (
            <div id={ids} className={classes}>
                <form action="" onSubmit={this.loadMods}>
                    <div className="input-group">
                        <select className="custom-select form-control" name="saveFile">
                            {saves}
                        </select>
                        <div className="input-group-append">
                            <button className="btn btn-outline-secondary" type="submit">Load Mods</button>
                        </div>
                    </div>
                </form>
            </div>
        )
    }
}

ModLoadSave.propTypes = {
    modContentClass: instanceOfModsContent.isRequired,
    className: PropTypes.string,
    id: PropTypes.string
}

export default ModLoadSave;
