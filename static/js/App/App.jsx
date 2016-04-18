import React from 'react';
import Hello from './components/hello.jsx';
import World from './components/world.jsx';

class App extends React.Component {
    render() {
        return(
            <div>
                <Hello />
                <World />
            </div>
        )
    }
}

export default App
