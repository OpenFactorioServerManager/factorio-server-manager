import React from 'react';
import {IndexLink} from 'react-router';
import ModOverview from './Mods/overview';

class ModsContent extends React.Component {
    constructor(props) {
        super(props);

        this.componentDidMount = this.componentDidMount.bind(this);
        this.loadModList = this.loadModList.bind(this);

        this.state = {
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
                    />
                </section>
            </div>
        )
    }
}

export default ModsContent;