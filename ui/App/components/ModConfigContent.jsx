import React from "react";
import {Link} from "react-router-dom";
import FontAwesomeIcon from "./FontAwesomeIcon";

export default class ModConfigContent extends React.Component {
    componentDidMount() {

    }

    render() {
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
                    <div className="box">
                        <div className="box-header">
                            <h3 className="box-title">Configure Mods</h3>
                        </div>

                        <div className="box-body">
                            <div key="test" className="form-group">
                                <label htmlFor="tests" className="control-label col-md-3">Test</label>
                                <div className="col-md-6">
                                    <input id="tests" className="form-control" type="text"/>
                                    <p className="help-block">Help text for Tests</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </section>
            </div>
        );
    }
}
