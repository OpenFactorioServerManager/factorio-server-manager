import React from 'react';
import ServerCtl from './ServerCtl/ServerCtl.jsx';
import ServerStatus from './ServerCtl/ServerStatus.jsx';
import FontAwesomeIcon from "./FontAwesomeIcon";

class Index extends React.Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {
        this.props.facServStatus();
        this.props.getSaves();
        this.props.getStatus();
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

                        <small className="float-sm-right">
                            <ol className="breadcrumb">
                                <li className="breadcrumb-item active">
                                    <FontAwesomeIcon icon="tachometer-alt"/>Server Control
                                </li>
                            </ol>
                        </small>
                    </h1>
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
