import React from 'react';
import {IndexLink} from 'react-router';
import ModList from './Mods/ListMods.jsx';
import InstalledMods from './Mods/InstalledMods.jsx';

class ModsContent extends React.Component {
    constructor(props) {
        super(props);
        this.componentDidMount = this.componentDidMount.bind(this);
        this.toggleMod = this.toggleMod.bind(this);
        this.loadInstalledModList = this.loadInstalledModList.bind(this);
        this.state = {
            installedMods: [],
            listMods: []
        };
    }

    componentDidMount() {
        this.loadModList();
        this.loadInstalledModList();
    }

    loadModList() {
        $.ajax({
            url: "/api/mods/list",
            dataType: "json",
            success: (data) => {
                this.setState({listMods: data.data.mods})
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
                    <li><IndexLink to="/"><i className="fa fa-dashboard"></i> Level</IndexLink></li>
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

                </section>
            </div>
        )
    }
}

export default ModsContent
