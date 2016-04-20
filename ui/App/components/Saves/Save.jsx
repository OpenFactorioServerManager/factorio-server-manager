import React from 'react';

class Save extends React.Component {
    downloadSave(e) {
        e.preventDefault();
        const node = this.refs.saveName;
        const saveName = node.name;
        $.ajax({
            url: "/api/saves/dl/" + saveName,
            dataType: "json",
            success: (data) => {
                console.log(data)
            },
            error: (xhr, status, err) => {
                console.log('api/mods/list', status, err.toString());
            }
        })
    }

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
                <td>
                    <form onSubmit={this.downloadSave.bind(this)}>
                        <input className='btn btn-default btn-sm'
                            ref='saveName'
                            type='submit'
                            value='Download Save'
                            name={this.props.save.name}
                        />
                    </form>
                </td>
            </tr>
        )
    }
}

Save.propTypes = {
    save: React.PropTypes.object.isRequired
}

export default Save
