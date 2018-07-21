import React from 'react';
import {Link} from 'react-router-dom';
import ServerCtl from './ServerCtl/ServerCtl.jsx';
import ServerStatus from './ServerCtl/ServerStatus.jsx';

class Index extends React.Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {
        this.props.facServStatus();
        this.props.getSaves();
        this.props.getStatus();
        console.log(this.props.serverStatus);
    }

    componentWillUnmount() {
        this.props.facServStatus();
    }

    render() {
        return(
            <div className="content-wrapper" style={{height: "100%"}}>
                <section className="content-header" style={{height: "100%"}}>
                    <h1>
                        Factorio Server
                        <small>Control your Factorio server</small>
                    </h1>
                    <ol className="breadcrumb">
                        <li><Link to="/"><i className="fa fa-dashboard"></i>Server Control</Link></li>
                    </ol>
                </section>

                <section className="content">
                    <ServerStatus
                        serverStatus={this.props.serverStatus}
                        facServStatus={this.props.facServStatus}
                        getStatus={this.props.getStatus}
                    />

                    <ServerCtl
                        getStatus={this.props.getStatus}
                        saves={this.props.saves}
                        getSaves={this.props.getSaves}
                        serverStatus={this.props.serverStatus}
                        facServStatus={this.props.facServStatus}
                    />
                </section>
            </div>
        )
    }
}

export default Index
