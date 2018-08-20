import React from 'react';
import PropTypes from 'prop-types';
import {Link} from 'react-router-dom';
import ModOverview from './Mods/ModOverview.jsx';
import locks from "locks";
import SemVer from 'semver';
import {ReactSwalNormal, ReactSwalDanger} from 'Utilities/customSwal';
import FontAwesomeIcon from "./FontAwesomeIcon";

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
        this.factorioLogoutHandler = this.factorioLogoutHandler.bind(this);
        this.deleteAllHandler = this.deleteAllHandler.bind(this);
        this.updateAllModsHandler = this.updateAllModsHandler.bind(this);
        this.updatesAvailable = this.updatesAvailable.bind(this);
        this.updateCountSubtract = this.updateCountSubtract.bind(this);
        this.updateCountAdd = this.updateCountAdd.bind(this);


        this.state = {
            loggedIn: false,
            installedMods: null,
            updatesAvailable: 0,
        };

        this.mutex = locks.createMutex();
    }

    componentDidMount() {
        this.loadModList();
        this.checkLoginState();
    }

    loadModList() {
        $.ajax({
            url: "/api/mods/list/installed",
            dataType: "json",
            success: (data) => {
                this.setState({installedMods: data.data})
            },
            error: (xhr, status, err) => {
                console.log('api/mods/list/installed', status, err.toString());
            }
        });
    }

    handlerFactorioLogin(e) {
        e.preventDefault();

        let $form = $(e.target);

        $.ajax({
            url: "/api/mods/factorio/login",
            method: "POST",
            data: $form.serialize(),
            dataType: "JSON",
            success: (data) => {
                ReactSwalNormal.fire({
                    title: "Logged in Successfully",
                    type: "success"
                });

                this.setState({
                    "loggedIn": data.data
                });
            },
            error: (jqXHR) => {
                ReactSwalNormal.fire({
                    title: jqXHR.responseJSON.data,
                    type: "error"
                });
            }
        });
    }

    checkLoginState() {
        $.ajax({
            url: "/api/mods/factorio/status",
            method: "POST",
            dataType: "json",
            success: (data) => {
                this.setState({
                    "loggedIn": data.data
                })
            },
            error: (jqXHR) => {
                let json_data = JSON.parse(jqXHR.responseJSON.data);

                console.log("error checking login status", json_data)
            }
        })
    }

    factorioLogoutHandler(e) {
        e.preventDefault();
        e.stopPropagation();

        $.ajax({
            url: "/api/mods/factorio/logout",
            method: "POST",
            dataType: "JSON",
            success: (data) => {
                this.setState({
                    loggedIn: data.data
                })
            },
            error: (jqXHR) => {
                ReactSwalNormal.fire({
                    title: "error logging out of factorio",
                    text: jqXHR.responseJSON.data,
                    type: "error"
                });
            }
        })
    }

    loadDownloadListSwalHandler() {
        let $checkedInput = $('input[name=version]:checked');
        let link = $checkedInput.data("link");
        let filename = $checkedInput.data("filename");
        let modName = $checkedInput.data("modid");

        $.ajax({
            method: "POST",
            url: "/api/mods/install",
            dataType: "JSON",
            data: {
                link: link,
                filename: filename,
                modName: modName
            },
            success: (data) => {
                this.setState({
                    installedMods: data.data.mods
                });

                ReactSwalNormal.fire({
                    title: "Mod installed",
                    type: "success"
                });
            },
            error: (jqXHR, status, err) => {
                ReactSwalNormal.fire({
                    title: "An error occurred",
                    text: err.toString(),
                    type: "error"
                });
            }
        });
    }

    //TODO remove modIdInput, when the factorio-mod-portal-api is fixed
    // all outcommented needs to be reimplemented, when it will work again
    loadDownloadList(e) {
        e.preventDefault();
        // let $button = $(e.target);
        // let $loader = $("<div class='loader'></div>");
        // $button.prepend($loader);
        let modId = $(e.target).find("input[name=modId]").val();
        // let modId = $button.data("modId");

        $.ajax({
            method: "POST",
            url: "/api/mods/details",
            data: {
                modId: modId
            },
            dataType: "json",
            success: (data) => {
                // $loader.remove();

                let correctData = JSON.parse(data.data);

                let checkboxes = [];
                let compatibleReleaseFound = false;

                correctData.releases.reverse();
                correctData.releases.forEach((release) => {
                    let incompatibleClass = "";
                    let isChecked = false;

                    if(!SemVer.satisfies(this.props.factorioVersion, release.info_json.factorio_version + ".x")) {
                        incompatibleClass = "incompatible";
                    } else if(compatibleReleaseFound == false) {
                        compatibleReleaseFound = true;
                        isChecked = true;
                    }

                    let date = new Date(release.released_at);

                    let singleBox = <tr className={incompatibleClass} key={release.version}>
                        <td>
                            <input type="radio"
                                   name="version"
                                   data-link={release.download_url}
                                   data-filename={release.file_name}
                                   data-modid={modId}
                                   defaultChecked={isChecked}
                            />
                        </td>
                        <td>
                            {release.version}
                        </td>
                        <td>
                            {release.info_json.factorio_version}
                        </td>
                        <td>
                            {date.toLocaleDateString()}
                        </td>
                        {/*<td>*/}
                            {/*{release.downloads_count}*/}
                        {/*</td>*/}
                    </tr>;

                    checkboxes.push(singleBox);
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
                            {/*<th>*/}
                                {/*Downloads*/}
                            {/*</th>*/}
                        </tr>
                    </thead>

                    <tbody>
                        {checkboxes}
                    </tbody>
                </table>;

                ReactSwalDanger.fire({
                    title: "Choose version",
                    html: table,
                    type: "question",
                    showCancelButton: true,
                    confirmButtonText: "Download it!",
                    showLoaderOnConfirm: true,
                    preConfirm: this.loadDownloadListSwalHandler
                });
            },
            error: (xhr, status, err) => {
                console.log('api/mods/details', status, err.toString());
                // $loader.remove();
            }
        })
    }

    toggleModHandler(e, updatesInProgress) {
        e.preventDefault();

        if(updatesInProgress) {
            ReactSwalNormal.fire({
                title: "Toggle mods failed",
                text: "Can't toggle the mod, when an update is still in progress",
                type: "error"
            });
            return false;
        }

        let $button = $(e.target);
        let $row = $button.parents("tr");
        let modName = $row.data("mod-name");

        $.ajax({
            url: "/api/mods/toggle",
            method: "POST",
            data: {
                modName: modName
            },
            dataType: "JSON",
            success: (data) => {
                if(data.success) {
                    this.mutex.lock(() => {
                        let installedMods = this.state.installedMods;

                        $.each(installedMods, (k, v) => {
                            if(v.name == modName) {
                                installedMods[k].enabled = data.data;
                            }
                        });

                        this.setState({
                            installedMods: installedMods
                        });

                        this.mutex.unlock();
                    });
                }
            },
            error: (jqXHR, status, err) => {
                console.log('api/mods/toggle', status, err.toString());
                ReactSwalNormal.fire({
                    title: "Toggle Mod went wrong",
                    text: err.toString(),
                    type: "error"
                });
            }
        });
    }

    deleteModHandler(e, updatesInProgress) {
        e.preventDefault();

        if(updatesInProgress) {
            ReactSwalNormal.fire({
                title: "Delete failed!",
                text: "Can't delete the mod, when an update is still in progress",
                type: "error"
            });
            return false;
        }

        let $button = $(e.target);
        let $row = $button.parents("tr");
        let modName = $row.data("mod-name");

        ReactSwalDanger({
            title: "Delete Mod?",
            text: "This will delete the mod and can break the save file",
            type: "question",
            showCancelButton: true,
            confirmButtonText: "Delete it!",
            showLoaderOnConfirm: true,
            preConfirm: () => {
                $.ajax({
                    url: "/api/mods/delete",
                    method: "POST",
                    data: {
                        modName: modName
                    },
                    dataType: "JSON",
                    success: (data) => {
                        if (data.success) {
                            this.mutex.lock(() => {
                                ReactSwalNormal.fire({
                                    title: <p>Delete of mod {modName} successful</p>,
                                    type: "success"
                                });
                                let installedMods = this.state.installedMods;

                                installedMods.forEach((v, k) => {
                                    if (v.name == modName) {
                                        delete installedMods[k];
                                    }
                                });

                                this.setState({
                                    installedMods: installedMods
                                });

                                this.mutex.unlock();
                            });
                        }
                    },
                    error: (jqXHR, status, err) => {
                        console.log('api/mods/delete', status, err.toString());
                        ReactSwalNormal.fire({
                            title: "Delete Mod went wrong",
                            text: err.toString(),
                            type: "error"
                        });
                    }
                });
            }
        });
    }

    deleteAllHandler(e) {
        e.preventDefault();
        e.stopPropagation();

        ReactSwalDanger.fire({
            title: "Delete Mod?",
            html: <p>This will delete ALL mods and can't be redone!<br/> Are you sure?</p>,
            type: "question",
            showCancelButton: true,
            confirmButtonText: "Yes, Delete ALL!",
            showLoaderOnConfirm: true,
            preConfirm: () => {
                $.ajax({
                    url: "/api/mods/delete/all",
                    method: "POST",
                    dataType: "JSON",
                    success: (data) => {
                        ReactSwalNormal.fire({
                            title: "ALL mods deleted successful",
                            type: "success"
                        })
                        this.setState({
                            installedMods: data.data
                        });
                    },
                    error: (jqXHR, status, err) => {
                        console.log('api/mods/delete/all', status, err.toString());
                        ReactSwalNormal.fire({
                            title: "Delete all mods went wrong",
                            html: <p>{err.toString()}<br/>{jqXHR.responseJSON.data}</p>,
                            type: "error",
                        });
                    }
                })
            }
        });
    }

    updateModHandler(e, toggleUpdateStatus, removeVersionAvailableStatus) {
        e.preventDefault();

        if(!this.state.loggedIn) {
            ReactSwalNormal.fire({
                title: "Update failed",
                text: "Please login into Factorio to update mod",
                type: "error",
            });

            let $addModBox = $('#add-mod-box');
            if($addModBox.hasClass("collapsed-box")) {
                $addModBox.find(".box-header").click();
            }
        } else {
            let $button = $(e.currentTarget);
            let downloadUrl = $button.data("downloadUrl");
            let filename = $button.data("fileName");
            let modName = $button.parents("tr").data("modName");

            //make button spinning
            toggleUpdateStatus();

            $.ajax({
                url: "/api/mods/update",
                method: "POST",
                data: {
                    downloadUrl: downloadUrl,
                    filename: filename,
                    modName: modName,
                },
                success: (data) => {
                    toggleUpdateStatus();
                    removeVersionAvailableStatus();

                    this.updateCountSubtract();

                    if(data.success) {
                        this.mutex.lock(() => {
                            ReactSwalNormal.fire({
                                title: <p>Update of mod {modName} successful</p>,
                                type: "success"
                            })
                            let installedMods = this.state.installedMods;

                            installedMods.forEach((v, k) => {
                                if(v.name == modName) {
                                    installedMods[k] = data.data;
                                }
                            });

                            this.setState({
                                installedMods: installedMods
                            });

                            this.mutex.unlock();
                        });
                    }
                },
                error: (jqXHR, status, err) => {
                    console.log('api/mods/delete', status, err.toString());
                    toggleUpdateStatus();

                    ReactSwalNormal.fire({
                        title: "Update Mod went wrong",
                        text: err.toString(),
                        type: "error"
                    });
                }
            });
        }
    }

    updatesAvailable() {
        this.setState({
            updatesAvailable: true
        });
    }

    updateAllModsHandler(e) {
        e.preventDefault();
        e.stopPropagation();

        let updateButtons = $('#manage-mods').find(".update-button");
        $.each(updateButtons, (k, v) => {
            v.click();
        });
    }

    updateCountSubtract() {
        this.setState({
            updatesAvailable: this.state.updatesAvailable > 0 ? this.state.updatesAvailable - 1 : 0
        });
    }

    updateCountAdd() {
        this.setState({
            updatesAvailable: this.state.updatesAvailable + 1
        });
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

                        <small className="float-sm-right">
                            <ol className="breadcrumb">
                                <li className="breadcrumb-item">
                                    <Link to="/"><FontAwesomeIcon icon="tachometer-alt"/>Server Control</Link>
                                </li>
                                <li className="breadcrumb-item active">
                                    <FontAwesomeIcon icon="pencil-alt"/>Mods
                                </li>
                            </ol>
                        </small>
                    </h1>
                </section>

                <section className="content">
                    <ModOverview
                        {...this.state}
                        {...this.props}
                        loadDownloadList={this.loadDownloadList}
                        submitFactorioLogin={this.handlerFactorioLogin}
                        toggleMod={this.toggleModHandler}
                        deleteMod={this.deleteModHandler}
                        deleteAll={this.deleteAllHandler}
                        updateMod={this.updateModHandler}
                        updateCountAdd={this.updateCountAdd}
                        uploadModSuccessHandler={this.uploadModSuccessHandler}
                        updateAllMods={this.updateAllModsHandler}
                        modContentClass={this}
                        factorioLogoutHandler={this.factorioLogoutHandler}
                    />
                </section>
            </div>
        )
    }
}

ModsContent.propTypes = {
    factorioVersion: PropTypes.string,
};

export default ModsContent;