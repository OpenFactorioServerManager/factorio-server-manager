import React from 'react';

class Save extends React.Component {
    render() {
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
            </tr>
        )
    }
}

Save.propTypes = {
    save: React.PropTypes.object.isRequired
}

export default Save
