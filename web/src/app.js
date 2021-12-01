import React from 'react';
import ReactDOM from 'react-dom';
import Navbar from './components/Navbar/nav-bar';
import AddItems from './components/AuctionItems/add-auction-items';
import ListItems from './components/AuctionItems/list-auction-items';
import UserLogIn  from './components/UserSignIn/user-login';
import MakeBid  from './components/AuctionItems/Bid/make-bid';
import { CookiesProvider } from 'react-cookie'


function App(){
    return (
        <CookiesProvider>
            <Navbar />
            {/* <AuctionItems /> */}
            <UserLogIn />
            <AddItems />
            <ListItems />
        </CookiesProvider>
            
    )
}

ReactDOM.render(App(), document.getElementById('app'));
