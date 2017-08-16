import React from 'react';
import Mod from './Mod.jsx';
import ModSearch from './search/ModSearch.jsx';

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
        console.log(this.props.installedMods)
        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Manage Mods</h3>
                </div>

                <div className="box-body">
                    <div className="box box-success collapsed-box">
                        <button className="btn btn-box-tool" type="button" data-widget="collapse">
                            <div className="box-header with-border">
                                <i className="fa fa-plus"></i>
                                <h4 className="box-title">Add Mods</h4>
                            </div>
                        </button>

                        <ModSearch
                            {...this.state}
                            {...this.props}
                            submitSearchMod={this.handlerSearchMod}
                            submitFactorioLogin={this.props.submitFactorioLogin}
                        />

                    </div>

                    <div className="table-responsive">
                        <table className="table table-striped">
                            <thead>
                            <tr>
                                <th>Name</th>
                                <th>Status</th>
                                <th>Toggle/Remove</th>
                            </tr>
                            </thead>
                            <tbody>
                            {this.props.installedMods.map ( (mod, i) => {
                                if(mod.name !== "base")
                                    return(
                                        <Mod
                                            key={i}
                                            mod={mod}
                                            {...this.props}
                                        />
                                    )
                            })}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        );
    }
}

ModOverview.propTypes = {
    installedMods: React.PropTypes.array.isRequired,
    submitFactorioLogin: React.PropTypes.func.isRequired,
    toggleMod: React.PropTypes.func.isRequired
};

export default ModOverview;