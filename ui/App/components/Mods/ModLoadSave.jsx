import React from 'react';

class ModOverview extends React.Component {
    constructor(props) {
        super(props);

        this.loadMods = this.loadMods.bind(this);
    }

    componentDidMount() {
        //Load Saves
        this.props.getSaves();
    }

    loadMods() {
        console.log("baum");
    }

    render() {
        console.log(this.props.saves);

        let saves = [];
        this.props.saves.forEach((value, index) => {
            if(index != this.props.saves.length - 1) {
                saves.push(
                    <option key={index} value={value.name}>
                        {value.name}
                    </option>
                )
            }
        });

        return (
            <div className="box-body">
                <div className="input-group">
                    <select className="custom-select form-control" id="inputGroupSelect04">
                        {saves}
                    </select>
                    <div className="input-group-append">
                        <button className="btn btn-outline-secondary" type="button" onClick={this.loadMods}>Load Mods</button>
                    </div>
                </div>
            </div>
        )
    }
}

export default ModOverview;
