import React from 'react';
import {IndexLink} from 'react-router';
import ModList from './Mods/ListMods.jsx';
import InstalledMods from './Mods/InstalledMods.jsx';
import ModPacks from './Mods/ModPacks.jsx'

class ModsContent extends React.Component {
    constructor(props) {
        super(props);
        this.componentDidMount = this.componentDidMount.bind(this);
        this.toggleMod = this.toggleMod.bind(this);
        this.loadInstalledModList = this.loadInstalledModList.bind(this);
        this.loadModPackList = this.loadModPackList.bind(this);
        this.state = {
            installedMods: [],
            listMods: [],
            modPacks: [],
        };
    }

    componentDidMount() {
        this.loadModList();
        this.loadInstalledModList();
        this.loadModPackList();
    }

    loadModList() {
        $.ajax({
            url: "/api/mods/list",
            dataType: "json",
            success: (data) => {
                console.log(data)
                if (data.success == true) {
                    this.setState({listMods: data.data.mods})
                } else {
                    this.setState({listMods: []})
                }
            },
            error: (xhr, status, err) => {
                console.log('api/mods/list', status, err.toString());
            }
        });
    }

    loadInstalledModList() {
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

    loadModPackList() {
        $.ajax({
            url: "/api/mods/packs/list",
            dataType: "json",
            success: (resp) => {
                if (resp.success === true) {
                    this.setState({modPacks: resp.data})
                    console.log(this.state)
                } else {
                    this.setState({modPacks: []})
                }
            }
        })
    }

    toggleMod(modName) {
        $.ajax({
            url: "/api/mods/toggle/" + modName,
            dataType: "json",
            success: (data) => {
                this.setState({listMods: data.data.mods})
            },
            error: (xhr, status, err) => {
                console.log('api/mods/toggle', status, err.toString());
            }
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

                    <InstalledMods 
                        {...this.state}
                        loadInstalledModList={this.loadInstalledModList}
                    />
                    <ModList
                        {...this.state}
                        toggleMod={this.toggleMod}
                    />
                    <ModPacks 
                        {...this.state}
                        loadModPackList={this.loadModPackList}
                    />

                </section>
            </div>
        )
    }
}

export default ModsContent
