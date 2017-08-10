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
        this.swalSubmitHandler = this.swalSubmitHandler.bind(this);

        this.state = {
            username: "",
            userKey: "",
            installedMods: []
        }
    }

    componentDidMount() {
        this.loadModList();
        //TODO get base stuff
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

    swalSubmitHandler() {
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
        // swal({
        //     title: "Select the version, that will be installed",
        //     showLoaderOnConfirm: true
        // });
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
                }, this.swalSubmitHandler);
            },
            error: (xhr, status, err) => {
                console.log('api/mods/details', status, err.toString());
                $loader.remove();
            }
        })
        // swal({
        //     title: "Ajax request example",
        //     text: "Submit to run ajax request",
        //     type: "info",
        //     showCancelButton: true,
        //     closeOnConfirm: false,
        //     showLoaderOnConfirm: true,
        // },
        // function(){
        //     setTimeout(function(){
        //         swal("Ajax request finished!");
        //     }, 2000);
        // });
    };

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
                    />
                </section>
            </div>
        )
    }
}

export default ModsContent;