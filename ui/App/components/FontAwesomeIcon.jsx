import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';

class FontAwesomeIcon extends React.Component {
    constructor(props) {
        super(props);

        console.log("huhu");
    }

    render() {
        let classes = classNames(this.props.prefix, {
            "fas": !this.props.prefix,
        }, 'fa-' + this.props.icon, this.props.className);
        console.log(classes);

        return (
            <i className={classes} {...this.props}></i>
        );
    }
}

FontAwesomeIcon.propTypes = {
    icon: PropTypes.string.isRequired,
    prefix: PropTypes.string,
    className: PropTypes.string
};

export default FontAwesomeIcon;
