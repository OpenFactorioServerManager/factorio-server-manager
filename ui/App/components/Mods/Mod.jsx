import React from 'react';

class Mod extends React.Component {
    render() {
        if (this.props.mod.enabled === false) {
            this.modStatus = <span className="label label-danger">Disabled</span>
        } else {
            this.modStatus = <span className="label label-success">Enabled</span>
        }
        return(
            <tr data-mod-name={this.props.mod.name}
                data-file-name={this.props.mod.file_name}
            >
                <td>{this.props.mod.title}</td>
                <td>{this.modStatus}</td>
                <td>
                    <input className='btn btn-default btn-sm'
                        ref='modName'
                        type='submit'
                        value='Toggle'
                        onClick={this.props.toggleMod}
                    />
                </td>
            </tr>
        )
    }
}

Mod.propTypes = {
    mod: React.PropTypes.object.isRequired,
    toggleMod: React.PropTypes.func.isRequired
}

export default Mod
