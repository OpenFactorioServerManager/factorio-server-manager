import React from "react";

const Layout = (props) => {

    return (<div className="container mx-auto">
        {props.children}
    </div>);
}

export default Layout;