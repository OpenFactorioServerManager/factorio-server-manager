import React from 'react';
import ReactDOMServer from 'react-dom/server';
import {instanceOfModsContent} from "./ModsPropTypes";

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

                    let modVersion = mod.version.major + "." + mod.version.minor + "." + mod.version.build;
                    let singleCheckbox = <tr key={mod.name}>
                        <td>
                            {mod.name}
                            <input type="hidden" name="mod_name" value={mod.name}/>
                        </td>
                        <td>
                            {modVersion}
                            <input type="hidden" name="mod_version" value={modVersion}/>
                        </td>
                    </tr>

                    checkboxes.push(singleCheckbox);
                });

                let table = <div>
                    All Mods will be installed
                    <div style={{display: "flex", width: "100%", justifyContent: "center"}}>
                        <form id="swalForm">
                            <table>
                                <thead>
                                    <tr>
                                        <th>
                                            ModName
                                        </th>
                                        <th>
                                            ModVersion
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

                swal({
                    title: "Mods to install",
                    text: ReactDOMServer.renderToStaticMarkup(table),
                    html: true,
                    type: 'info',
                    showCancelButton: true,
                    closeOnConfirm: false,
                    confirmButtonText: "Download Mods!",
                    cancelButtonText: "Cancel",
                    showLoaderOnConfirm: true
                }, this.loadModsSwalHandler);
            },
            error: (jqXHR) => {
                swal({
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
                swal({
                    title: "All Mods installed successfully!",
                    type: "success"
                });

                this.props.modContentClass.setState({
                    installedMods: data.data.mods
                });
            },
            error: (jqXHR) => {
                let json_data = JSON.parse(jqXHR.responseJSON.data);

                swal({
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

        return (
            <div className="box-body">
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
}

export default ModLoadSave;
