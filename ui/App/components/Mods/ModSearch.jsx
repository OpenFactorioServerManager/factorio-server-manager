import React from 'react';

class ModSearch extends React.Component {
    render() {
        if(this.props.userKey) {
            return (
                <form onSubmit={this.props.submitSearchMod}>
                    <div className="input-group col-lg-5">
                        <input type="text" className="form-control" placeholder="Search for Mod" name="search" />
                        <span className="input-group-btn">
                            <input className="btn btn-default" type="submit" value="Go!"/>
                        </span>
                    </div>
                </form>
            )
        } else {
            return (
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
            )
        }
    }
}

ModSearch.propTypes = {
    submitSearchMod: React.PropTypes.func.isRequired,
    userKey: React.PropTypes.string.isRequired,
    submitFactorioLogin: React.PropTypes.func.isRequired
}

export default ModSearch;