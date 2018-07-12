import React from 'react';

class ModLoadSave extends React.Component {
    constructor(props) {
        super(props);

        this.loadMods = this.loadMods.bind(this);
    }

    componentDidMount() {
        //Load Saves
        this.props.getSaves();
    }

    loadMods(e) {
        e.preventDefault();

        // let save = $(e.target).find("select").val();

        $.ajax({
            url: "/api/mods/save/load",
            method: "POST",
            data: $(e.target).serialize(),
            dataType: "JSON",
            success: (data) => {
            },
            error: (jqXHR) => {
                let json_data = JSON.parse(jqXHR.responseJSON.data);

                swal({
                    title: json_data.detail,
                    type: "error"
                });
            }
        });
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
                <form action="" onSubmit={this.loadMods}>
                    <div className="input-group">
                        <select className="custom-select form-control" name="saveFile">
                            {saves}
                        </select>
                        <div className="input-group-append">
                            <button className="btn btn-outline-secondary" type="submit">Load Mods</button>
                        </div>
                    </div>
                </form>
            </div>
        )
    }
}

export default ModLoadSave;
