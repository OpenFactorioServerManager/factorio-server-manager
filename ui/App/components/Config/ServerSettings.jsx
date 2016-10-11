import React from 'react';

class ServerSettings extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        return(
            <tr key={this.props.name}>
                <td>{this.props.name}</td>
                <td>{this.props.setting}</td>
            </tr>
        )
    }
}

ServerSettings.propTypes = {
    name: React.PropTypes.string.isRequired,
    setting: React.PropTypes.oneOfType([
        React.PropTypes.string,
        React.PropTypes.number,
        React.PropTypes.boolean,
        React.PropTypes.array,
    ]),
}

export default ServerSettings
