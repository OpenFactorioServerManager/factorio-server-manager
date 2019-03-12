import React from "react";
import {Link} from "react-router-dom";
import FontAwesomeIcon from "./FontAwesomeIcon";
import {ReactSwalDanger} from "../../js/customSwal";

export default class ModConfigContent extends React.Component {
    constructor(props) {
        super(props);

        this.updateConfig = this.updateConfig.bind(this);

        this.state = {
            config: null
        };
    }

    componentDidMount() {
        this.loadConfig();
    }

    loadConfig() {
        $.ajax({
            url: "/api/mods/settings",
            type: "GET",
            dataType: "json",
            success: (data) => {
                if(data.success) {
                    this.setState({
                        config: data.data
                    });
                } else {
                    ReactSwalDanger.fire({
                        title: "Loading mod-settings failed",
                        text: data.data,
                        type: "error"
                    });
                }
            },
            error: (xhr, status, err) => {
                console.log('api/server/status', status, err.toString());
                let json_data = JSON.parse(jqXHR.responseJSON.data);

                ReactSwalNormal.fire({
                    title: json_data.detail,
                    type: "error"
                });
            }
        });
    }

    updateConfig(e) {
        let config = this.state.config;
        let configElem = $('.mod-config');

        $.each(configElem, (_, v) => {
            let configName = v.id;
            let configValue = v.value;
            let category = v.getAttribute("data-cat");
            let type = v.getAttribute("data-type");

            switch (type) {
                case "boolean":
                    configValue = configValue == "true" ? true : false;
                    break;
                case "number":
                    configValue = parseFloat(configValue);
                    break;
            }

            if (!config[category]) {
                config[category] = [];
            }
            if (!config[category][configName]) {
                config[category][configName] = {};
            }
            config[category][configName].value = configValue;
        });

        let configString = JSON.stringify(config);
        $.ajax({
            url: "/api/mods/settings",
            type: "POST",
            data: {
                data: configString
            },
            dataType: "json",
            success: (data) => {
                if(data.success) {
                    this.setState({
                        config: config
                    });
                } else {
                    ReactSwalDanger.fire({
                        title: "Loading mod-settings failed",
                        text: data.data,
                        type: "error"
                    });
                }
            },
            error: (xhr, status, err) => {
                console.log('api/server/status', status, err.toString());
                let json_data = JSON.parse(jqXHR.responseJSON.data);

                ReactSwalNormal.fire({
                    title: json_data.detail,
                    type: "error"
                });
            }
        });
    }

    render() {
        let categories = [];
        if(this.state.config) {
            for(let confCat in this.state.config) {
                let singles = [];

                for(let confSingle in this.state.config[confCat]) {
                    let singleType = typeof this.state.config[confCat][confSingle].value;
                    let input;
                    switch (singleType) {
                        case "boolean":
                            input = <select id={confSingle}
                                            className="form-control mod-config"
                                            data-cat={confCat}
                                            data-type={singleType}
                                            defaultValue={this.state.config[confCat][confSingle].value}
                            >
                                <option value="true">True</option>
                                <option value="false">False</option>
                            </select>
                            break;
                        case "number":
                            input = <input id={confSingle}
                                           className="form-control mod-config"
                                           type="number"
                                           defaultValue={this.state.config[confCat][confSingle].value}
                                           data-cat={confCat}
                                           data-type={singleType}
                            />
                            break;
                        case "string":
                            input = <input id={confSingle}
                                           className="form-control mod-config"
                                           type="text"
                                           defaultValue={this.state.config[confCat][confSingle].value}
                                           data-cat={confCat}
                                           data-type={singleType}
                            />
                            break;
                        default:
                            input = <input id={confSingle}
                                           className="form-control mod-config"
                                           type="text"
                                           defaultValue={ this.state.config[confCat][confSingle].value }
                                           data-cat={confCat}
                                           data-type="default"
                            />
                            break;
                    }

                    singles.push(<div key={confSingle} className="form-group">
                        <label htmlFor="tests" className="control-label col-md-3">{confSingle}</label>
                        <div className="col-md-6">
                            {input}
                        </div>
                    </div>);
                }

                categories.push(<div className="box" key={confCat}>
                    <div className="box-header">
                        <h3 className="box-title">{confCat}</h3>
                    </div>

                    <div className="box-body">
                        {singles}
                    </div>
                </div>);
            }
        }

        return (
            <div id="mod-config" className="content-wrapper">
                <section className="content-header">
                    <h1>
                        Mod-Configuration

                        <small className="float-sm-right">
                            <ol className="breadcrumb">
                                <li className="breadcrumb-item">
                                    <Link to="/"><FontAwesomeIcon icon="tachometer-alt" className="fa-fw"/>Server Control</Link>
                                </li>
                                <li className="breadcrumb-item active">
                                    <FontAwesomeIcon icon="cogs"/>Mod-Configuration
                                </li>
                            </ol>
                        </small>
                    </h1>
                </section>

                <section className="content">
                    <div className="alert alert-warning alert-dismissible" role="alert">
                        Mod settings only can have specific values (like enums). If invalid values are set, they will get reset,
                        when the factorio-server will load those. Currently it's not implemented, to read the possible values!
                    </div>

                    {categories}

                    <input type="submit" className="btn btn-success" value="Update Settings" onClick={this.updateConfig}/>
                </section>
            </div>
        );
    }
}
