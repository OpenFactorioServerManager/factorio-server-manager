import React from 'react';

class Settings extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        return(
            <tbody>
            {Object.keys(this.props.config).map(function(key) {
                return(
                    <tr key={key}>
                        <td>{key}</td>
                        <td>{this.props.config[key]}</td>
                    </tr>
                )                                                  
            }, this)}        
            </tbody>
        )
    }

}

Settings.propTypes = {
    section: React.PropTypes.string.isRequired,
    config: React.PropTypes.object.isRequired,
}

export default Settings
