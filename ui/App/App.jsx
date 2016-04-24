import React from 'react';
import Header from './components/Header.jsx';
import Sidebar from './components/Sidebar.jsx';
import Footer from './components/Footer.jsx';
import HiddenSidebar from './components/HiddenSidebar.jsx';


class App extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return(
            <div className="wrapper">

                <Header />

                <Sidebar />
                
                {React.cloneElement(
                    this.props.children,
                    {message: ""}
                )}

                <Footer />

                <HiddenSidebar />

            </div>
        )
    }
}

export default App
