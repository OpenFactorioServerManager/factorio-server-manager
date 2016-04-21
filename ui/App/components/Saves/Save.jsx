import React from 'react';

class Save extends React.Component {
    render() {
        let saveLocation = "/api/saves/dl/" + this.props.save.name 
        let saveSize = parseFloat(this.props.save.size / 1024 / 1024).toFixed(3)
        let saveLastMod = Date.parse(this.props.save.last_mod);
        let date = new Date(saveLastMod)
        let dateFmt = date.getFullYear() + '-' + date.getMonth() + '-' + date.getDay() + ' '
            + date.getHours() + ':' + date.getMinutes() + ':' + date.getSeconds();

        return(
            <tr>
                <td>{this.props.save.name}</td>
                <td>{dateFmt}</td>
                <td>{saveSize} MB</td>
                <td>
                    <a className="btn btn-default" href={saveLocation}>Download</a>
                </td>
                <td>
                    <button
                        className="btn btn-danger btn-small" 
                        ref="saveInput"
                        type="button" 
                        onClick={this.props.removeSave.bind(this, this.props.saves[this.props.index].name)}>
                    <i className="fa fa-trash"></i>
                    &nbsp;
                    Delete
                    </button>
                </td>
</tr>
        )
    }
}

Save.propTypes = {
    save: React.PropTypes.object.isRequired,
    saves: React.PropTypes.array.isRequired,
    index: React.PropTypes.number.isRequired,
    removeSave: React.PropTypes.func.isRequired
}

export default Save
