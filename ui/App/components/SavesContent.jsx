import React from 'react';
import SavesList from './Saves/SavesList.jsx';

class SavesContent extends React.Component {
    constructor(props) {
        super(props);
        this.dlSave = this.dlSave.bind(this);
        this.getSaves = this.getSaves.bind(this);
        this.state = {
            saves: []
        }
    }

    componentDidMount() {
        this.getSaves();
    }

    getSaves() {
        $.ajax({
            url: "/api/saves/list",
            dataType: "json",
            success: (data) => {
                this.setState({saves: data.data})
            },
            error: (xhr, status, err) => {
                console.log('api/mods/list', status, err.toString());
            }
        })
    }

    dlSave(saveName) {
        $.ajax({
            url: "/api/saves/dl/" + saveName,
            dataType: "json",
            success: (data) => {
                console.log("Downloading save: " + saveName)
            },
            error: (xhr, status, err) => {
                console.log('api/mods/list', status, err.toString());
            }
        })
    }

    render() {
        return(
            <div className="content-wrapper">
                <section className="content-header">
                <h1>
                    Saves
                    <small>Factorio Save Files</small>
                </h1>
                <ol className="breadcrumb">
                    <li><a href="#"><i className="fa fa-dashboard"></i> Level</a></li>
                    <li className="active">Here</li>
                </ol>
                </section>

                <section className="content">

                <SavesList 
                    {...this.state}
                    dlSave={this.dlSave}
                    getSaves={this.getSaves}
                />


                </section>
            </div>
        )
    }
}

export default SavesContent
