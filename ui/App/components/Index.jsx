import React from 'react';
import ServerCtl from './ServerCtl/ServerCtl.jsx';

class Index extends React.Component {
    constructor(props) {
        super(props);

    }

    componentDidMount() {
        this.props.facServerStatus();
        this.props.getSaves();
    }

    componentWillUnmount() {
        this.props.facServerStatus();
    }

    render() {
        return(
            <div className="content-wrapper" style={{height: "100%"}}>
                <section className="content-header">
                <h1>
                    Index
                    <small>Optional description</small>
                </h1>
                <ol className="breadcrumb">
                    <li><a href="#"><i className="fa fa-dashboard"></i> Level</a></li>
                    <li className="active">Here</li>
                </ol>
                </section>

                <section className="content">

                <ServerCtl 
                    saves={this.props.saves}
                    getSaves={this.props.getSaves}
                />


                </section>
            </div>
        )
    }
}

export default Index
