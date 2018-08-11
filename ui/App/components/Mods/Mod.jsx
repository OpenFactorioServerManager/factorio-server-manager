import React from 'react';
import PropTypes from 'prop-types';
import SemVer from 'semver';
import FontAwesomeIcon from "../FontAwesomeIcon";

class Mod extends React.Component {
    constructor(props) {
        super(props);

        this.toggleUpdateStatus = this.toggleUpdateStatus.bind(this);
        this.removeVersionAvailableStatus = this.removeVersionAvailableStatus.bind(this);

        this.state = {
            newVersionAvailable: false,
            updateInProgress: false
        }

        this.versionAjax = {};
    }

    componentDidMount() {
        this.checkForNewVersion();
    }

    componentDidUpdate(prevProps, prevState) {
        if(prevProps.mod.version != this.props.mod.version) {
            this.checkForNewVersion();
        }
    }

    componentWillUnmount() {
        this.versionAjax.abort();
    }

    checkForNewVersion() {
        //send AJAX that will check this
        this.versionAjax = $.ajax({
            url: "/api/mods/details",
            method: "POST",
            data: {
                modId: this.props.mod.name
            },
            dataType: "JSON",
            success: (data) => {
                let newData = JSON.parse(data.data);
                //get newest COMPATIBLE release
                let newestRelease;
                newData.releases.forEach((release) => {
                    //TODO change to info_json dependency (when mod-portal-api is working again)
                    if(SemVer.satisfies(this.props.factorioVersion, release.info_json.factorio_version + ".x")) {
                        if(!newestRelease) {
                            newestRelease = release;
                        } else if(SemVer.gt(release.version, newestRelease.version)) {
                            newestRelease = release;
                        }
                    }
                });

                if(newestRelease && newestRelease.version != this.props.mod.version) {
                    if(this.props.updateCountAdd)
                        this.props.updateCountAdd();

                    this.setState({
                        newVersionAvailable: true,
                        newVersion: {
                            downloadUrl: newestRelease.download_url,
                            file_name: newestRelease.file_name,
                            version: newestRelease.version
                        }
                    });
                } else {
                    this.setState({
                        newVersionAvailable: false,
                        newVersion: null
                    });
                }
            },
            error: (jqXHR, status, err) => {
                console.log('api/mods/details', status, err.toString());
            }
        });
    }

    toggleUpdateStatus() {
        console.log("update Status changed");
        this.setState({
            updateInProgress: !this.state.updateInProgress
        })
    }

    removeVersionAvailableStatus() {
        this.setState({
            newVersionAvailable: false,
            newVersion: null
        })
    }

    render() {
        let modStatus;
        if(this.props.mod.enabled === false) {
            modStatus = <span className="badge badge-danger">Disabled</span>
        } else {
            modStatus = <span className="badge badge-success">Enabled</span>
        }

        let version;
        if(this.state.newVersionAvailable) {
            let faArrow;
            if(SemVer.gt(this.state.newVersion.version, this.props.mod.version)) {
                faArrow = "arrow-circle-up";
            } else {
                faArrow = "arrow-circle-down";
            }

            version = <span>{this.props.mod.version}
                <a className="btn btn-xs btn-default update-button"
                   style={{
                       marginLeft: 10,
                       display: "inline-flex",
                       justifyContent: "center",
                       alignItems: "center",
                       width: 30,
                       height: 30
                   }}
                   href="#"
                   onClick={(event) => {
                       this.state.updateInProgress && this.props.updateMod != null ? null : this.props.updateMod(event, this.toggleUpdateStatus, this.removeVersionAvailableStatus);
                   }}
                   data-download-url={this.state.newVersion.downloadUrl}
                   data-file-name={this.state.newVersion.file_name}
                >
                    {
                        this.state.updateInProgress ?
                            <div className='loader' style={{width: 15, height: 15, marginRight: 0, borderWidth: 3,}}></div>
                            :
                            <FontAwesomeIcon icon={faArrow} title="Update Mod" style={{fontSize: "15pt"}}/>
                    }
                </a>
            </span>;
        } else {
            version = this.props.mod.version;
        }

        let factorioVersion;
        if(!this.props.mod.compatibility) {
            factorioVersion = <span style={{color: "red"}}>
                {this.props.mod.factorio_version}&nbsp;&nbsp;
                <sup>not compatible</sup>
            </span>
        } else {
            factorioVersion = this.props.mod.factorio_version;
        }

        return(
            <tr data-mod-name={this.props.mod.name}
                data-file-name={this.props.mod.file_name}
                data-mod-version={this.props.mod.version}
            >
                <td>{this.props.mod.title}</td>
                <td>{modStatus}</td>
                <td>{version}</td>
                <td>{factorioVersion}</td>
                <td>
                    <input className='btn btn-default btn-sm'
                        ref='modName'
                        type='submit'
                        value='Toggle'
                        onClick={(event) => this.props.toggleMod(event, this.state.updateInProgress)}
                        disabled={this.state.updateInProgress}
                    />

                    <input className="btn btn-danger btn-sm"
                        style={{marginLeft: 25}}
                        ref="modName"
                        type="submit"
                        value="Delete"
                        onClick={(event) => this.props.deleteMod(event, this.state.updateInProgress)}
                        disabled={this.state.updateInProgress}
                    />
                </td>
            </tr>
        )
    }
}

Mod.propTypes = {
    mod: PropTypes.object.isRequired,
    toggleMod: PropTypes.func.isRequired,
    deleteMod: PropTypes.func.isRequired,
    updateMod: PropTypes.func.isRequired,
    updateCountAdd: PropTypes.func,
    factorioVersion: PropTypes.string,
};

export default Mod
