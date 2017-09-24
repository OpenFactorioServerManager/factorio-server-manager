import React from 'react';
import ReactDOMServer from 'react-dom/server';
import {IndexLink} from 'react-router';
import ModOverview from './Mods/ModOverview.jsx';

class ModsContent extends React.Component {
    constructor(props) {
        super(props);

        this.componentDidMount = this.componentDidMount.bind(this);
        this.loadModList = this.loadModList.bind(this);
        this.handlerFactorioLogin = this.handlerFactorioLogin.bind(this);
        this.loadDownloadList = this.loadDownloadList.bind(this);
        this.loadDownloadListSwalHandler = this.loadDownloadListSwalHandler.bind(this);
        this.toggleModHandler = this.toggleModHandler.bind(this);
        this.deleteModHandler = this.deleteModHandler.bind(this);
        this.updateModHandler = this.updateModHandler.bind(this);
        this.uploadModSuccessHandler = this.uploadModSuccessHandler.bind(this);

        this.state = {
            username: sessionStorage.getItem("username") || "",
            userKey: sessionStorage.getItem("userKey") || "",
            installedMods: [],
        }
    }

    componentDidMount() {
        this.loadModList();
    }

    loadModList() {
        $.ajax({
            url: "/api/mods/list/installed",
            dataType: "json",
            success: (data) => {
                this.setState({installedMods: data.data})
            },
            error: (xhr, status, err) => {
                console.log('api/mods/list', status, err.toString());
            }
        });
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

                let user_key = (JSON.parse(data.data))[0];

                sessionStorage.setItem("username", username);
                sessionStorage.setItem("userKey", user_key);

                this.setState({
                    "username": username,
                    "userKey": user_key
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

    loadDownloadListSwalHandler() {
        let $checked_input = $('input[name=version]:checked');
        let link = $checked_input.data("link");
        let filename = $checked_input.data("filename");
        let mod_name = $checked_input.data("modid");

        $.ajax({
            method: "POST",
            url: "/api/mods/install",
            dataType: "JSON",
            data: {
                username: this.state.username,
                userKey: this.state.userKey,
                link: link,
                filename: filename,
                modName: mod_name
            },
            success: (data) => {
                this.setState({
                    installedMods: data.data
                })

                swal({
                    type: "success",
                    title: "Mod installed"
                });
            },
            error: (jqXHR, status, err) => {
                swal({
                    type: "error",
                    title: "some error occured",
                    text: err.toString()
                });
            }
        });
    }

    loadDownloadList(e) {
        e.preventDefault();
        let $button = $(e.target);
        let $loader = $("<div class='loader'></div>");
        $button.prepend($loader);
        let mod_id = $button.data("modId");

        $.ajax({
            method: "POST",
            url: "/api/mods/details",
            data: {
                mod_id: mod_id
            },
            dataType: "json",
            success: (data) => {
                $loader.remove();

                let correct_data = JSON.parse(data.data);

                let checkboxes = []
                correct_data.releases.forEach((release, index) => {
                    let date = new Date(release.released_at);

                    let single_box = <tr>
                        <td>
                            <input type="radio"
                                   name="version"
                                   data-link={release.download_url}
                                   data-filename={release.file_name}
                                   data-modid={mod_id}
                                   checked={index == 0 ? true : false}
                            />
                        </td>
                        <td>
                            {release.info_json.version}
                        </td>
                        <td>
                            {release.info_json.factorio_version}
                        </td>
                        <td>
                            {date.toLocaleDateString()}
                        </td>
                        <td>
                            {release.downloads_count}
                        </td>
                    </tr>;

                    checkboxes.push(single_box);
                });

                let table = <table>
                    <thead>
                        <tr>
                            <th></th>
                            <th>
                                Version
                            </th>
                            <th>
                                Game Version
                            </th>
                            <th>
                                Release Date
                            </th>
                            <th>
                                Downloads
                            </th>
                        </tr>
                    </thead>

                    <tbody>
                        {checkboxes}
                    </tbody>
                </table>;

                swal({
                    title: "Choose version",
                    text: ReactDOMServer.renderToStaticMarkup(table),
                    html: true,
                    type: "info",
                    showCancelButton: true,
                    closeOnConfirm: false,
                    confirmButtonText: "Download it!",
                    cancelButtonText: "Close",
                    showLoaderOnConfirm: true,
                }, this.loadDownloadListSwalHandler);
            },
            error: (xhr, status, err) => {
                console.log('api/mods/details', status, err.toString());
                $loader.remove();
            }
        })
    }

    toggleModHandler(e) {
        e.preventDefault();
        let $button = $(e.target);
        let $row = $button.parents("tr");
        let mod_name = $row.data("mod-name");

        $.ajax({
            url: "/api/mods/toggle",
            method: "POST",
            data: {
                mod_name: mod_name
            },
            dataType: "JSON",
            success: (data) => {
                this.setState({
                    installedMods: data.data
                });
            },
            error: (jqXHR, status, err) => {
                console.log('api/mods/toggle', status, err.toString());
                swal({
                    title: "Toggle Mod went wrong",
                    text: err.toString(),
                    type: "error"
                });
            }
        });
    }

    deleteModHandler(e) {
        e.preventDefault();
        let $button = $(e.target);
        let $row = $button.parents("tr");
        let mod_name = $row.data("mod-name");
        let class_this = this;

        swal({
            title: "Delete Mod?",
            text: "This will delete the mod and can break the save file",
            type: "info",
            showCancelButton: true,
            closeOnConfirm: false,
            confirmButtonText: "Delete it!",
            cancelButtonText: "Close",
            showLoaderOnConfirm: true,
            confirmButtonColor: "#DD6B55",
        }, function () {
            $.ajax({
                url: "/api/mods/delete",
                method: "POST",
                data: {
                    mod_name: mod_name
                },
                dataType: "JSON",
                success: (data) => {
                    swal("Delete of mod " + mod_name + " successful", "", "success");
                    class_this.setState({
                        installedMods: data.data
                    });
                },
                error: (jqXHR, status, err) => {
                    console.log('api/mods/delete', status, err.toString());
                    swal({
                        title: "Delete Mod went wrong",
                        text: err.toString(),
                        type: "error"
                    });
                }
            });
        });
    }

    updateModHandler(e, toggleUpdateStatus, removeVersionAvailableStatus) {
        e.preventDefault();

        if(!this.state.userKey) {
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
            let modname = $button.parents("tr").data("modName");

            let this_class = this;

            //make button spinning
            toggleUpdateStatus();

            $.ajax({
                url: "/api/mods/update",
                method: "POST",
                data: {
                    username: this.state.username,
                    userKey: this.state.userKey,
                    downloadUrl: download_url,
                    filename: filename,
                    mod_name: modname,
                },
                success: (data) => {
                    toggleUpdateStatus();
                    removeVersionAvailableStatus();
                    this_class.setState({
                        installedMods: data.data.mods
                    });
                },
                error: (jqXHR, status, err) => {
                    console.log('api/mods/delete', status, err.toString());
                    toggleUpdateStatus();
                    swal({
                        title: "Update Mod went wrong",
                        text: err.toString(),
                        type: "error"
                    });
                }
            });
        }
    }

    uploadModSuccessHandler(event, data) {
        this.setState({
            installedMods: data.response.data.mods
        });
    }

    render() {
        return(
            <div className="content-wrapper">
                <section className="content-header">
                    <h1>
                        Mods
                        <small>Manage your mods</small>
                    </h1>
                    <ol className="breadcrumb">
                        <li><IndexLink to="/"><i className="fa fa-dashboard fa-fw"></i>Server Control</IndexLink></li>
                        <li className="active">Here</li>
                    </ol>
                </section>

                <section className="content">
                    <ModOverview
                        {...this.state}
                        loadDownloadList={this.loadDownloadList}
                        submitFactorioLogin={this.handlerFactorioLogin}
                        toggleMod={this.toggleModHandler}
                        deleteMod={this.deleteModHandler}
                        updateMod={this.updateModHandler}
                        uploadModSuccessHandler={this.uploadModSuccessHandler}
                        modContentClass={this}
                    />
                </section>
            </div>
        )
    }
}

export default ModsContent;