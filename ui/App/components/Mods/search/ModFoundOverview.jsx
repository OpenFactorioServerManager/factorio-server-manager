import React from 'react';
import PropTypes from 'prop-types';

class ModFoundOverview extends React.Component {
    render() {
        let imgStyle= {
            width: 144,
            height: 144,
            border: "1px outset #333",
            borderRadius: 2,
        }
        let noImgStyle = {
            container: {
                width: 144,
                height: 144,
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
                backgroundColor: "#333",
                border: "1px inset #333",
                fontSize: 20,
                color: "#949391",
                border: "1px outset #333",
                borderRadius: 2,
            },
        }
        let informationStyle = {
            container: {
                marginLeft: 20,
            }
        }

        let mods = [];
        this.props.shownModList.some((mod, index) => {
            if(index == 10) return true;
            let img =
                (mod.first_media_file != null) ?
                    <img src={mod.first_media_file.urls.thumb} style={imgStyle} />
                    :
                    <div style={noImgStyle.container}>
                        <div>No picture</div>
                    </div>;

            mods.push(
                <div className="list-group-item" key={mod.title}>
                    <div style={{display: "flex"}}>
                        {img}
                        <div style={informationStyle.container}>
                            <h4 className="list-group-item-heading">{mod.title} <small>by {mod.owner}</small></h4>
                            <div className="list-group-item-text">{mod.summary}</div>
                            <button style={{marginTop: 10, display: "flex"}} onClick={this.props.loadDownloadList} data-mod-id={mod.name}>INSTALL</button>
                        </div>
                    </div>
                </div>
            );
        });

        return (
            <div className="list-group">
                {mods}
            </div>
        );
    }
}

ModFoundOverview.propTypes = {
    shownModList: PropTypes.array.isRequired,
    loadDownloadList: PropTypes.func.isRequired
}

export default ModFoundOverview;
