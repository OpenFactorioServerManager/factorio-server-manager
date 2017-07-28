import React from 'react';

class Mod extends React.Component {
    render() {
        if (this.props.mod.enabled === false) {
            this.modStatus = <span className="label label-danger">Disabled</span>
        } else {
            this.modStatus = <span className="label label-success">Enabled</span>
        }
        return(
            <tr>
                <td>{this.props.mod.name}</td>
                <td>{this.modStatus}</td>
                <td>
                    <input className='btn btn-default btn-sm'
                        ref='modName'
                        type='submit'
                        value='Toggle'
                        name={this.props.mod.name}
                           //TODO toggle and remove mod
                    />
                </td>
            </tr>
        )
    }
}

Mod.propTypes = {
    mod: React.PropTypes.object.isRequired
}

export default Mod
