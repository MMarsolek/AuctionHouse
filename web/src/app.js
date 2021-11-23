import React from 'react';
import ReactDOM from 'react-dom';
import Navbar from './components/Navbar/nav-bar';
// import AuctionItems from './components/AuctionItems/auction-items';
import UserLogIn  from './components/UserSignIn/user-login';


function App(){
    return (
        <div className="app">
            <Navbar />
            {/* <AuctionItems /> */}
            <UserLogIn />
        </div>
    )
}

ReactDOM.render(App(), document.getElementById('app'));
