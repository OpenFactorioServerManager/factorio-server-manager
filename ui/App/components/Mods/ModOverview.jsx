import React from 'react';
import Mod from './Mod.jsx';
import ModSearch from './ModSearch.jsx';

class ModOverview extends React.Component {
    constructor(props) {
        super(props);

        this.handlerFactorioLogin = this.handlerFactorioLogin.bind(this);

        this.state = {
            username: "",
            userKey: "",
            shownModList: []
        }
    }

    handlerSearchMod(e) {
        console.log($(e.target).find("input").val());
        e.preventDefault();

        $.ajax({
            url: "/api/mods/search",
            method: "GET",
            data: $(e.target).serialize(),
            dataType: "JSON",
            success: (data) => {
                console.log(data);
            },
            error: (jqXHR) => {
                console.log(jqXHR.statusText);
            }
        })
    }

    handlerFactorioLogin(e) {
        e.preventDefault();

        let $form = $(e.target);
        let username = $form.find('input[name=username]').val();

        $.ajax({
            url: "/api/mods/factorio/login",
            method: "POST",
            data: $form.serialize(),
            dataType: "JSON",
            success: (data) => {
                swal({
                    title: "Logged in Successfully",
                    type: "success"
                });

                this.setState({
                    "username": username,
                    "userKey": (JSON.parse(data.data))[0]
                });
            },
            error: (jqXHR) => {
                let json_data = JSON.parse(jqXHR.responseJSON.data);

                swal({
                    title: json_data.message,
                    type: "error"
                });
            }
        });
    }

    render() {
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

                        <div className="box-body">
                            <ModSearch
                                {...this.state}
                                submitSearchMod={this.handlerSearchMod}
                                submitFactorioLogin={this.handlerFactorioLogin}
                            />
                        </div>
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
    installedMods: React.PropTypes.array.isRequired
};

export default ModOverview;