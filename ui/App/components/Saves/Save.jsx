import React from 'react';
import PropTypes from 'prop-types';
import FontAwesomeIcon from "../FontAwesomeIcon";

class Save extends React.Component {
    render() {
        let saveLocation = "/api/saves/dl/" + this.props.save.name 
        let saveSize = parseFloat(this.props.save.size / 1024 / 1024).toFixed(3)
        let saveLastMod = Date.parse(this.props.save.last_mod);
        let date = new Date(saveLastMod)
        let dateFmt = date.toISOString().replace('T', ' ').split('.')[0]; //Outputs date as "YYYY-MM-DD HH:MM:SS" with zero-padding and 24h

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
                        onClick={this.props.removeSave.bind(this, this.props.saves[this.props.index].name)}
                    >
                        <FontAwesomeIcon icon="trash"/>
                        &nbsp;
                        Delete
                    </button>
                </td>
</tr>
        )
    }
}

Save.propTypes = {
    save: PropTypes.object.isRequired,
    saves: PropTypes.array.isRequired,
    index: PropTypes.number.isRequired,
    removeSave: PropTypes.func.isRequired
}

export default Save
