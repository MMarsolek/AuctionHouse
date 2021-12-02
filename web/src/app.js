
import ReactDOM from 'react-dom';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './auth-provider';
import Navbar from './components/Navbar/nav-bar';
import AddItems from './components/AuctionItems/add-auction-items';
import ListItems from './components/AuctionItems/list-auction-items';
import UserLogIn  from './components/UserSignIn/user-login';
import { CookiesProvider } from 'react-cookie';
import UserCreation from './components/UserManagement/user-creation';



function App(){
    return (
        <CookiesProvider>
            <BrowserRouter>
                <AuthProvider>
                    <Navbar />
                    <Routes>
                        <Route index path="/" element={<UserLogIn />} />
                        <Route path="/auctions" element={<ListItems />} />
                        <Route path="/createUser" element={<UserCreation />} />
                        <Route path="/createItem" element={<AddItems />} />
                    </Routes>
                </AuthProvider>
            </BrowserRouter>
        </CookiesProvider>
    )
}

ReactDOM.render(App(), document.getElementById('app'));
