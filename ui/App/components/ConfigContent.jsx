import React from 'react';
import Settings from './Config/Settings.jsx';

class ConfigContent extends React.Component {
    constructor(props) {
        super(props);
        this.getConfig = this.getConfig.bind(this);
        this.state = {
            config: {}
        }
    }

    componentDidMount() {
        this.getConfig();
    }
    
    getConfig() {
        $.ajax({
            url: "/api/config",
            dataType: "json",
            success: (data) => {
                this.setState({config: data.data})
                console.log(this.state.config)
            },
            error: (xhr, status, err) => {
                console.log('/api/config/get', status, err.toString());
            }
        });
    }
    render() {
        return(
            <div className="content-wrapper">
                <section className="content-header">
                <h1>
                    Config 
                    <small>Manage server configuration</small>
                </h1>
                <ol className="breadcrumb">
                    <li><a href="#"><i className="fa fa-dashboard"></i> Level</a></li>
                    <li className="active">Here</li>
                </ol>
                </section>

                <section className="content">
                    <div className="box">
                        <div className="box-header">
                            <h3 className="box-title">Manage Server Configuration</h3>
                        </div>
                        
                        <div className="box-body">
    
                        {Object.keys(this.state.config).map(function(key) {
                            var conf = this.state.config[key]
                            return(
                                <div className="settings-section" key={key}>
                                <h3>{key}</h3>
                                    <div className="table-responsive">
                                    <table className="table table-striped">
                                        <thead>
                                            <tr>
                                                <th>Setting name</th>
                                                <th>Setting value</th>
                                            </tr>
                                        </thead>
                                            <Settings
                                                section={key}
                                                config={conf}    
                                            />
                                    </table>
                                    </div>
                                </div>
                            )
                        }, this)}
                        </div>
                    </div>
                </section>
            </div>
        )
    }
}

export default ConfigContent
