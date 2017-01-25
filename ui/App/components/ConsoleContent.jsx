import React from 'react';
import {IndexLink} from 'react-router';

class ConsoleContent extends React.Component {
    constructor(props) {
        super(props);
        this.componentDidMount = this.componentDidMount.bind(this);
        this.handleInput = this.handleInput.bind(this);
        this.clearInput = this.clearInput.bind(this);
        this.clearHistory = this.clearHistory.bind(this);
        this.addHistory = this.addHistory.bind(this);
        this.handleClick = this.handleClick.bind(this);
        this.newLogLine = this.newLogLine.bind(this);
        this.state = {
            commands: {},
            history: [],
            prompt: '$ ',
        }
    }

    componentDidMount() {
        this.props.socket.emit("log subscribe");
        this.setState({connected: true});
        
        this.props.socket.on('log update', this.newLogLine.bind(this));
    }

    componentDidUpdate() {
      var el = this.refs.output;
      var container = document.getElementById("console-output");
      container.scrollTop = this.refs.output.scrollHeight;
    }

    handleInput(e) {
        if (e.key === "Enter") {
            var input_text = this.refs.term.value;
            this.props.socket.emit("command send", input_text);
            this.addHistory(this.state.prompt + " " + input_text);
            this.clearInput();
        }
    }

    clearInput() {
        this.refs.term.value = "";
    }

    clearHistory() {
        ths.setState({ history: [] });
    }

    addHistory(output) {
        var history = this.state.history;
        history.push(output);
        this.setState({
            'history': history
        });
    }

    handleClick() {
        var term = this.refs.term;
        term.focus();
    }

    newLogLine(logline) {
        var history = this.state.history;
        history.push(logline);
        this.setState({
            'history': history
        });
    }

    render() {
        var output = this.state.history.map((op, i) => {
            return <p key={i}>{op}</p>
        });

        return(
            <div className="content-wrapper">
                <section className="content-header">
                    <h1>
                        Server Console
                        <small>Send commands and messages to the Factorio server</small>
                    </h1>
                    <ol className="breadcrumb">
                        <li><IndexLink to="/"><i className="fa fa-dashboard"></i>Server Control</IndexLink></li>
                        <li className="active">Here</li>
                    </ol>
                </section>

                <section className="content">
                <div className="console-box" >
                    <div id='console-output' className='console-container' onClick={this.handleClick} ref="output">
                        {output}
                    </div>
                    <p>
                        <span className="console-prompt-box">{this.state.prompt}
                        <input type="text" onKeyPress={this.handleInput} ref="term" /></span>
                    </p>

                </div>
                </section>
            </div>
        );
    }
}

ConsoleContent.propTypes = {
    socket: React.PropTypes.object.isRequired,
}

export default ConsoleContent;
