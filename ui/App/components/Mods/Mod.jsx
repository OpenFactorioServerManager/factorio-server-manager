import React from 'react';

class Mod extends React.Component {
    togglePress(e) {
        e.preventDefault();
        console.log(this.refs.modName);
        const node = this.refs.modName;
        const modName = node.name;
        this.props.toggleMod(modName);
    }

    render() {
        if (this.props.mod.enabled === "false") {
            this.modStatus = <span className="label label-danger">Disabled</span>
        } else {
            this.modStatus = <span className="label label-success">Enabled</span>
        }
        return(
            <tr>
                <td>{this.props.mod.name}</td>
                <td>{this.modStatus}</td>
                <td>
                    <form onSubmit={this.togglePress.bind(this)}>
                        <input className='btn btn-default btn-sm'
                            ref='modName'
                            type='submit'
                            value='Toggle'
                            name={this.props.mod.name}
                        />
                    </form>
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
