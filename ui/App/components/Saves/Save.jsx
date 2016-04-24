import React from 'react';

class Save extends React.Component {
    hours12(date) { return (date.getHours() + 24) % 12 || 12; }

    render() {
        let saveLocation = "/api/saves/dl/" + this.props.save.name 
        let saveSize = parseFloat(this.props.save.size / 1024 / 1024).toFixed(3)
        let saveLastMod = Date.parse(this.props.save.last_mod);
        let date = new Date(saveLastMod)
        let dateFmt = date.getFullYear() + '-' + date.getMonth() + '-' + date.getDay() + '   '
            + this.hours12(date) + ':' + date.getMinutes() + ':' + date.getSeconds();

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
