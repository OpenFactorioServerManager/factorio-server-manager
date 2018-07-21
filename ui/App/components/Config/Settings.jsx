import React from 'react';
import PropTypes from 'prop-types';

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
    section: PropTypes.string.isRequired,
    config: PropTypes.object.isRequired,
}

export default Settings
