import React from 'react';
import PropTypes from 'prop-types';

class ModSearch extends React.Component {
    render() {
        let classes = "card-body" + " " + this.props.className;
        let ids = this.props.id;

        if(this.props.loggedIn) {
            //TODO switch back to currently commented out code, when the mod-portal-api is back with all features!!
            /*return (
                <div className="box-body">
                    <form onSubmit={this.props.submitSearchMod}>
                        <div className="input-group col-lg-5">
                            <input type="text" className="form-control" placeholder="Search for Mod" name="search" />
                            <span className="input-group-btn">
                                <input className="btn btn-default" type="submit" value="Go!"/>
                            </span>
                        </div>
                    </form>

                    <ModFoundOverview
                        {...this.props}
                    />
                </div>
            )*/
            return (
                <div id={ids} className={classes}>
                    <form onSubmit={this.props.loadDownloadList}>
                        <div className="input-group col-lg-5">
                            <input type="text" className="form-control" placeholder="Download mod by ID" name="modId" />
                            <span className="input-group-btn">
                                <input className="btn btn-default" type="submit" value="Go!"/>
                            </span>
                        </div>
                    </form>
                </div>
            )
        } else {
            return (
                <div id={ids} className={classes}>
                    <form onSubmit={this.props.submitFactorioLogin}>
                        <h4>Login into Factorio</h4>
                        <div className="form-group">
                            <label htmlFor="factorio-account-name">Factorio Account Name:</label>
                            <input type="text" className="form-control" id="factorio-account-name" name="username" required />
                        </div>
                        <div className="form-group">
                            <label htmlFor="pwd">Password:</label>
                            <input type="password" className="form-control" id="pwd" name="password" required />
                        </div>
                        <input type="submit" className="btn btn-default" value="Login" />
                    </form>
                </div>
            )
        }
    }
}

ModSearch.propTypes = {
    submitSearchMod: PropTypes.func.isRequired,
    loggedIn: PropTypes.bool.isRequired,
    submitFactorioLogin: PropTypes.func.isRequired,
    loadDownloadList: PropTypes.func.isRequired,
    className: PropTypes.string,
    id: PropTypes.string
}

export default ModSearch;