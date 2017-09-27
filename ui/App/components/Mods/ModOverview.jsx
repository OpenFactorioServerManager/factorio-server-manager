import React from 'react';
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

    render() {
        return(
            <div>
                <div className="box collapsed-box" id="add-mod-box">
                    <div className="box-header" data-widget="collapse" style={{cursor: "pointer"}}>
                        <i className="fa fa-plus"></i>
                        <h3 className="box-title">Add Mod</h3>
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

                <div className="box">
                    <div className="box-header" data-widget="collapse" style={{cursor: "pointer"}}>
                        <i className="fa fa-minus"></i>
                        <h3 className="box-title">Manage Mods</h3>
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
    updateMod: React.PropTypes.func.isRequired,
    uploadModSuccessHandler: React.PropTypes.func.isRequired,

    modContentClass: instanceOfModsContent.isRequired,
};

export default ModOverview;